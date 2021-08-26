package env

import (
	"encoding/json"
	"os"
	"strings"
)

// GetGoctlVersion obtains from the environment variable GOCTL_VERSION, prior to 1.1.11,
// the goctl version was 1.1.10 by default.
// the goctl version is set at runtime in the environment variable GOCTL_VERSION,
// see the detail at https://github.com/tal-tech/go-zero/blob/master/tools/goctl/goctl.go
func GetGoctlVersion() string {
	currentVersion := os.Getenv("GOCTL_VERSION")
	if currentVersion == "" {
		currentVersion = "1.1.10"
	}
	return currentVersion
}

var tag = map[string]int{"pre-alpha": 0, "alpha": 1, "pre-bata": 2, "beta": 3, "released": 4, "": 5}

// IsVersionGatherThan compares whether the current goctl version
// is gather than the target version
func IsVersionGatherThan(version, target string) bool {
	versionNumber, versionTag := convertVersion(version)
	targetVersionNumber, targetTag := convertVersion(target)
	if versionNumber > targetVersionNumber {
		return true
	} else if versionNumber < targetVersionNumber {
		return false
	} else { // unchecked case, in normal, the goctl version does not contains suffix in release.
		return tag[versionTag] > tag[targetTag]
	}
}

// version format: number[.number]*(-tag)
func convertVersion(version string) (versionNumber float64, tag string) {
	splits := strings.Split(version, "-")
	tag = strings.Join(splits[1:], "")
	var flag bool
	numberStr := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}

		if r == '.' {
			if flag {
				return '_'
			}
			flag = true
			return r
		}
		return '_'
	}, splits[0])
	numberStr = strings.ReplaceAll(numberStr, "_", "")
	versionNumber, _ = json.Number(numberStr).Float64()
	return
}
