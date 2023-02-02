package vben

import (
	_ "embed"
)

var (
	//go:embed api.tpl
	apiTpl string

	//go:embed model.tpl
	modelTpl string

	//go:embed data.tpl
	dataTpl string

	//go:embed drawer.tpl
	drawerTpl string

	//go:embed index.tpl
	indexTpl string

	//go:embed locale.tpl
	localeTpl string

	//go:embed statusrender.tpl
	statusRenderTpl string
)
