package internal

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/zeromicro/go-zero/core/iox"
	"github.com/zeromicro/go-zero/core/lang"
	"golang.org/x/sys/unix"
)

const cgroupDir = "/sys/fs/cgroup"

var (
	isUnifiedOnce sync.Once
	isUnified     bool
	inUserNS      bool
	nsOnce        sync.Once
)

// IsCgroup2UnifiedMode returns whether we are running in cgroup v2 unified mode.
func IsCgroup2UnifiedMode() bool {
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

type cgroup struct {
	cgroups map[string]string
}

func (c *cgroup) acctUsageAllCpus() (uint64, error) {
	data, err := iox.ReadText(path.Join(c.cgroups["cpuacct"], "cpuacct.usage"))
	if err != nil {
		return 0, err
	}

	return parseUint(data)
}

func (c *cgroup) acctUsagePerCpu() ([]uint64, error) {
	data, err := iox.ReadText(path.Join(c.cgroups["cpuacct"], "cpuacct.usage_percpu"))
	if err != nil {
		return nil, err
	}

	var usage []uint64
	for _, v := range strings.Fields(data) {
		u, err := parseUint(v)
		if err != nil {
			return nil, err
		}

		usage = append(usage, u)
	}

	return usage, nil
}

func (c *cgroup) cpuQuotaUs() (int64, error) {
	data, err := iox.ReadText(path.Join(c.cgroups["cpu"], "cpu.cfs_quota_us"))
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(data, 10, 64)
}

func (c *cgroup) cpuPeriodUs() (uint64, error) {
	data, err := iox.ReadText(path.Join(c.cgroups["cpu"], "cpu.cfs_period_us"))
	if err != nil {
		return 0, err
	}

	return parseUint(data)
}

func (c *cgroup) cpus() ([]uint64, error) {
	data, err := iox.ReadText(path.Join(c.cgroups["cpuset"], "cpuset.cpus"))
	if err != nil {
		return nil, err
	}

	return parseUints(data)
}

func currentCgroup() (*cgroup, error) {
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

	return &cgroup{
		cgroups: cgroups,
	}, nil
}

func parseUint(s string) (uint64, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		if err.(*strconv.NumError).Err == strconv.ErrRange {
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

	ints := make(map[uint64]lang.PlaceholderType)
	cols := strings.Split(val, ",")
	for _, r := range cols {
		if strings.Contains(r, "-") {
			fields := strings.SplitN(r, "-", 2)
			min, err := parseUint(fields[0])
			if err != nil {
				return nil, fmt.Errorf("cgroup: bad int list format: %s", val)
			}

			max, err := parseUint(fields[1])
			if err != nil {
				return nil, fmt.Errorf("cgroup: bad int list format: %s", val)
			}

			if max < min {
				return nil, fmt.Errorf("cgroup: bad int list format: %s", val)
			}

			for i := min; i <= max; i++ {
				ints[i] = lang.Placeholder
			}
		} else {
			v, err := parseUint(r)
			if err != nil {
				return nil, err
			}

			ints[v] = lang.Placeholder
		}
	}

	var sets []uint64
	for k := range ints {
		sets = append(sets, k)
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

		/*
		 * We assume we are in the initial user namespace if we have a full
		 * range - 4294967295 uids starting at uid 0.
		 */
		if a == 0 && b == 0 && c == 4294967295 {
			return
		}
		inUserNS = true
	})
	return inUserNS
}
