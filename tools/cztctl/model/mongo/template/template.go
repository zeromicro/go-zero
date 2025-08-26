package template

import _ "embed"

var (
	// ModelText provides the default template for model to generate.
	//go:embed model.tpl
	ModelText string

	// ModelCustomText provides the default template for model to generate.
	//go:embed model_custom.tpl
	ModelCustomText string

	// ModelTypesText provides the default template for model to generate.
	//go:embed types.tpl
	ModelTypesText string

	// Error provides the default template for error definition in mongo code generation.
	//go:embed error.tpl
	Error string
)
