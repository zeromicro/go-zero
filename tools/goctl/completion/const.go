package completion

const (
	BashCompletionFlag        = `generate-goctl-completion`
	defaultCompletionFilename = "goctl_autocomplete"
)

const (
	magic = 1 << iota
	flagZsh
	flagBash
)
