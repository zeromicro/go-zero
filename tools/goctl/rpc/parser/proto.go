package parser

// Proto describes a proto file,
type Proto struct {
	Src       string
	Name      string
	Package   Package
	PbPackage string
	GoPackage string
	Import    []Import
	Message   []Message
	Service   Services
}
