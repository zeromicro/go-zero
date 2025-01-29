package version

import (
	"encoding/json"
	"strings"
)

// BuildVersion is the version of goctl.
const BuildVersion = "1.7.6"

var tag = map[string]int{"pre-alpha": 0, "alpha": 1, "pre-bata": 2, "beta": 3, "released": 4, "": 5}

// GetGoctlVersion returns BuildVersion
func GetGoctlVersion() string {
	return BuildVersion
}

// IsVersionGreaterThan compares whether the current goctl version
// is greater than the target version
func IsVersionGreaterThan(version, target string) bool {
	versionNumber, versionTag := convertVersion(version)
	targetVersionNumber, targetTag := convertVersion(target)
	if versionNumber > targetVersionNumber {
		return true
	} else if versionNumber < targetVersionNumber {
		return false
	} else {
		// unchecked case, in normal, the goctl version does not contain suffix in release.
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
	numberStr = strings.Replace(numberStr, "_", "", -1)
	versionNumber, _ = json.Number(numberStr).Float64()
	return
}
