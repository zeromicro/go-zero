package makefile

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/suyuan32/knife/core/io/filex"
)

// extractInfo extracts the information to context
func extractInfo(g *GenContext) error {
	makefileData, err := filex.ReadFileString(g.TargetPath)
	if err != nil {
		return err
	}

	if strings.Contains(makefileData, "gen-api") {
		if strings.Contains(makefileData, "gen-ent") {
			g.UseEnt = true
			g.IsSingle = true
		} else {
			g.IsApi = true
		}
	}

	if strings.Contains(makefileData, "gen-rpc") {
		g.IsRpc = true
		if strings.Contains(makefileData, "gen-ent") {
			g.UseEnt = true
		}
	}

	dataSplit := strings.Split(makefileData, "\n")

	if g.Style == "" && !strings.Contains(makefileData, "PROJECT_STYLE=") {
		return errors.New("style not set, use -s to set style")
	} else if g.Style == "" && strings.Contains(makefileData, "PROJECT_STYLE=") {
		style := findDefined("PROJECT_STYLE=", dataSplit)
		if style == "" {
			return errors.New("failed to find style definition, please set it manually by -s")
		}
		g.Style = style
	}

	if g.ServiceName == "" {
		g.ServiceName = findDefined("SERVICE", dataSplit)
	}

	return err
}

func findDefined(target string, data []string) string {
	for _, v := range data {
		if strings.Contains(v, target) {
			dataSplit := strings.Split(v, "=")
			if len(dataSplit) == 2 {
				return dataSplit[1]
			} else {
				return ""
			}
		}
	}

	return ""
}
