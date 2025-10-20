package migrationnotes

import (
	"github.com/lerity-yao/go-zero/tools/cztctl/config"
	"github.com/lerity-yao/go-zero/tools/cztctl/util/format"
)

// BeforeCommands run before command run to show some migration notes
func BeforeCommands(dir, style string) error {
	if err := migrateBefore1_3_4(dir, style); err != nil {
		return err
	}
	return nil
}

func getModelSuffix(style string) (string, error) {
	cfg, err := config.NewConfig(style)
	if err != nil {
		return "", err
	}
	baseSuffix, err := format.FileNamingFormat(cfg.NamingFormat, "_model")
	if err != nil {
		return "", err
	}
	return baseSuffix + ".go", nil
}
