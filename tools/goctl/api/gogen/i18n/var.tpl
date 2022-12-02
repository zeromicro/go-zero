package i18n

import (
	"embed"
)

//go:embed locale/*.json
var LocaleFS embed.FS
