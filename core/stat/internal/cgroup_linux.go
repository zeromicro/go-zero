package internal

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/tal-tech/go-zero/core/iox"
	"github.com/tal-tech/go-zero/core/lang"
)

const cgroupDir = "/sys/fs/cgroup"

type cgroup struct {
	cgroups map[string]string
}

func (c *cgroup) acctUsageAllCpus() (uint64, error) {
	data, err := iox.ReadText(path.Join(c.cgroups["cpuacct"], "cpuacct.usage"))
	if err != nil {
		return 0, err
	}

	return parseUint(string(data))
}

func (c *cgroup) acctUsagePerCpu() ([]uint64, error) {
	data, err := iox.ReadText(path.Join(c.cgroups["cpuacct"], "cpuacct.usage_percpu"))
	if err != nil {
		return nil, err
	}

	var usage []uint64
	for _, v := range strings.Fields(string(data)) {
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

	return strconv.ParseInt(string(data), 10, 64)
}

func (c *cgroup) cpuPeriodUs() (uint64, error) {
	data, err := iox.ReadText(path.Join(c.cgroups["cpu"], "cpu.cfs_period_us"))
	if err != nil {
		return 0, err
	}

	return parseUint(string(data))
}

func (c *cgroup) cpus() ([]uint64, error) {
	data, err := iox.ReadText(path.Join(c.cgroups["cpuset"], "cpuset.cpus"))
	if err != nil {
		return nil, err
	}

	return parseUints(string(data))
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

		cgroups[subsys] = path.Join(cgroupDir, subsys)
		if strings.Contains(subsys, ",") {
			for _, k := range strings.Split(subsys, ",") {
				cgroups[k] = path.Join(cgroupDir, k)
			}
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
		} else {
			return 0, fmt.Errorf("cgroup: bad int format: %s", s)
		}
	} else {
		if v < 0 {
			return 0, nil
		} else {
			return uint64(v), nil
		}
	}
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
