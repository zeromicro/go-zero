package internal

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/iox"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	cpuTicks  = 100
	cpuFields = 8
)

var (
	preSystem uint64
	preTotal  uint64
)

// if /proc not present, ignore the cpu calculation, like wsl linux
func init() {
	preSystem, err := systemCpuUsage()
	if err != nil {
		logx.Error(err)
		return
	}

	preTotal, err = totalCpuUsage()
	if err != nil {
		logx.Error(err)
		return
	}
}

// RefreshCpu refreshes cpu usage and returns.
func RefreshCpu() uint64 {
	total, err := totalCpuUsage()
	if err != nil {
		return 0
	}
	system, err := systemCpuUsage()
	if err != nil {
		return 0
	}

	var usage uint64
	cpuDelta := total - preTotal
	systemDelta := system - preSystem
	if cpuDelta > 0 && systemDelta > 0 {
		usage = uint64(float64(cpuDelta*1e3) / float64(systemDelta))
	}
	preSystem = system
	preTotal = total

	return usage
}

func cpuQuota() (int64, error) {
	cg, err := currentCgroup()
	if err != nil {
		return 0, err
	}

	return cg.cpuQuotaUs()
}

func cpuPeriod() (uint64, error) {
	cg, err := currentCgroup()
	if err != nil {
		return 0, err
	}

	return cg.cpuPeriodUs()
}

func cpuSets() ([]uint64, error) {
	cg, err := currentCgroup()
	if err != nil {
		return nil, err
	}

	return cg.cpus()
}

func systemCpuUsage() (uint64, error) {
	lines, err := iox.ReadTextLines("/proc/stat", iox.WithoutBlank())
	if err != nil {
		return 0, err
	}

	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			if len(fields) < cpuFields {
				return 0, fmt.Errorf("bad format of cpu stats")
			}

			var totalClockTicks uint64
			for _, i := range fields[1:cpuFields] {
				v, err := parseUint(i)
				if err != nil {
					return 0, err
				}

				totalClockTicks += v
			}

			return (totalClockTicks * uint64(time.Second)) / cpuTicks, nil
		}
	}

	return 0, errors.New("bad stats format")
}

func totalCpuUsage() (usage uint64, err error) {
	var cg cgroup
	if cg, err = currentCgroup(); err != nil {
		return
	}

	return cg.usageAllCpus()
}
