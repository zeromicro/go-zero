package template

import _ "embed"

// Text provides the default template for model to generate.
//go:embed model.tpl
var Text string

// Error provides the default template for error definition in mongo code generation.
//go:embed error.tpl
var Error string
