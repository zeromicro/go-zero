package internal

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/iox"
	"github.com/zeromicro/go-zero/core/lang"
	"golang.org/x/sys/unix"
)

const (
	cgroupDir   = "/sys/fs/cgroup"
	cpuMaxFile  = cgroupDir + "/cpu.max"
	cpuStatFile = cgroupDir + "/cpu.stat"
	cpusetFile  = cgroupDir + "/cpuset.cpus.effective"
)

var (
	isUnifiedOnce sync.Once
	isUnified     bool
	inUserNS      bool
	nsOnce        sync.Once
)

type cgroup interface {
	cpuQuota() (float64, error)
	cpuUsage() (uint64, error)
	effectiveCpus() (int, error)
}

func currentCgroup() (cgroup, error) {
	if isCgroup2UnifiedMode() {
		return currentCgroupV2()
	}

	return currentCgroupV1()
}

type cgroupV1 struct {
	cgroups map[string]string
}

func (c *cgroupV1) cpuQuota() (float64, error) {
	quotaUs, err := c.cpuQuotaUs()
	if err != nil {
		return 0, err
	}

	if quotaUs == -1 {
		return -1, nil
	}

	periodUs, err := c.cpuPeriodUs()
	if err != nil {
		return 0, err
	}

	return float64(quotaUs) / float64(periodUs), nil
}

func (c *cgroupV1) cpuPeriodUs() (uint64, error) {
	data, err := iox.ReadText(path.Join(c.cgroups["cpu"], "cpu.cfs_period_us"))
	if err != nil {
		return 0, err
	}

	return parseUint(data)
}

func (c *cgroupV1) cpuQuotaUs() (int64, error) {
	data, err := iox.ReadText(path.Join(c.cgroups["cpu"], "cpu.cfs_quota_us"))
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(data, 10, 64)
}

func (c *cgroupV1) cpuUsage() (uint64, error) {
	data, err := iox.ReadText(path.Join(c.cgroups["cpuacct"], "cpuacct.usage"))
	if err != nil {
		return 0, err
	}

	return parseUint(data)
}

func (c *cgroupV1) effectiveCpus() (int, error) {
	data, err := iox.ReadText(path.Join(c.cgroups["cpuset"], "cpuset.cpus"))
	if err != nil {
		return 0, err
	}

	cpus, err := parseUints(data)
	if err != nil {
		return 0, err
	}

	return len(cpus), nil
}

type cgroupV2 struct {
	cgroups map[string]string
}

func (c *cgroupV2) cpuQuota() (float64, error) {
	data, err := iox.ReadText(cpuMaxFile)
	if err != nil {
		return 0, err
	}

	fields := strings.Fields(data)
	if len(fields) != 2 {
		return 0, fmt.Errorf("cgroup: bad /sys/fs/cgroup/cpu.max file: %s", data)
	}

	if fields[0] == "max" {
		return -1, nil
	}

	quotaUs, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, err
	}

	periodUs, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return 0, err
	}

	return float64(quotaUs) / float64(periodUs), nil
}

func (c *cgroupV2) cpuUsage() (uint64, error) {
	usec, err := parseUint(c.cgroups["usage_usec"])
	if err != nil {
		return 0, err
	}

	return usec * uint64(time.Microsecond), nil
}

func (c *cgroupV2) effectiveCpus() (int, error) {
	data, err := iox.ReadText(cpusetFile)
	if err != nil {
		return 0, err
	}

	cpus, err := parseUints(data)
	if err != nil {
		return 0, err
	}

	return len(cpus), nil
}

func currentCgroupV1() (cgroup, error) {
	cgroupFile := fmt.Sprintf("/proc/%d/cgroup", os.Getpid())
	lines, err := iox.ReadTextLines(cgroupFile, iox.WithoutBlank())
	if err != nil {
		return nil, err
	}

	cgroups := make(map[string]string)
	for _, line := range lines {
		cols := strings.Split(line, ":")
		if len(cols) != 3 {
			return nil, fmt.Errorf("invalid cgroup line: %s", line)
		}

		subsys := cols[1]
		// only read cpu staff
		if !strings.HasPrefix(subsys, "cpu") {
			continue
		}

		// https://man7.org/linux/man-pages/man7/cgroups.7.html
		// comma-separated list of controllers for cgroup version 1
		fields := strings.Split(subsys, ",")
		for _, val := range fields {
			cgroups[val] = path.Join(cgroupDir, val)
		}
	}

	return &cgroupV1{
		cgroups: cgroups,
	}, nil
}

func currentCgroupV2() (cgroup, error) {
	lines, err := iox.ReadTextLines(cpuStatFile, iox.WithoutBlank())
	if err != nil {
		return nil, err
	}

	cgroups := make(map[string]string)
	for _, line := range lines {
		cols := strings.Fields(line)
		if len(cols) != 2 {
			return nil, fmt.Errorf("invalid cgroupV2 line: %s", line)
		}

		cgroups[cols[0]] = cols[1]
	}

	return &cgroupV2{
		cgroups: cgroups,
	}, nil
}

// isCgroup2UnifiedMode returns whether we are running in cgroup v2 unified mode.
func isCgroup2UnifiedMode() bool {
	isUnifiedOnce.Do(func() {
		var st unix.Statfs_t
		err := unix.Statfs(cgroupDir, &st)
		if err != nil {
			if os.IsNotExist(err) && runningInUserNS() {
				// ignore the "not found" error if running in userns
				isUnified = false
				return
			}
			panic(fmt.Sprintf("cannot statfs cgroup root: %s", err))
		}
		isUnified = st.Type == unix.CGROUP2_SUPER_MAGIC
	})

	return isUnified
}

func parseUint(s string) (uint64, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		if errors.Is(err, strconv.ErrRange) {
			return 0, nil
		}

		return 0, fmt.Errorf("cgroup: bad int format: %s", s)
	}

	if v < 0 {
		return 0, nil
	}

	return uint64(v), nil
}

func parseUints(val string) ([]uint64, error) {
	if val == "" {
		return nil, nil
	}

	var sets []uint64
	ints := make(map[uint64]lang.PlaceholderType)
	cols := strings.Split(val, ",")
	for _, r := range cols {
		if strings.Contains(r, "-") {
			fields := strings.SplitN(r, "-", 2)
			minimum, err := parseUint(fields[0])
			if err != nil {
				return nil, fmt.Errorf("cgroup: bad int list format: %s", val)
			}

			maximum, err := parseUint(fields[1])
			if err != nil {
				return nil, fmt.Errorf("cgroup: bad int list format: %s", val)
			}

			if maximum < minimum {
				return nil, fmt.Errorf("cgroup: bad int list format: %s", val)
			}

			for i := minimum; i <= maximum; i++ {
				if _, ok := ints[i]; !ok {
					ints[i] = lang.Placeholder
					sets = append(sets, i)
				}
			}
		} else {
			v, err := parseUint(r)
			if err != nil {
				return nil, err
			}

			if _, ok := ints[v]; !ok {
				ints[v] = lang.Placeholder
				sets = append(sets, v)
			}
		}
	}

	return sets, nil
}

// runningInUserNS detects whether we are currently running in a user namespace.
func runningInUserNS() bool {
	nsOnce.Do(func() {
		file, err := os.Open("/proc/self/uid_map")
		if err != nil {
			// This kernel-provided file only exists if user namespaces are supported
			return
		}
		defer file.Close()

		buf := bufio.NewReader(file)
		l, _, err := buf.ReadLine()
		if err != nil {
			return
		}

		line := string(l)
		var a, b, c int64
		fmt.Sscanf(line, "%d %d %d", &a, &b, &c)

		// We assume we are in the initial user namespace if we have a full
		// range - 4294967295 uids starting at uid 0.
		if a == 0 && b == 0 && c == math.MaxUint32 {
			return
		}

		inUserNS = true
	})

	return inUserNS
}
