package importstack

import "errors"

// ErrImportCycleNotAllowed defines an error for circular importing
var ErrImportCycleNotAllowed = errors.New("import cycle not allowed")

// ImportStack a stack of import paths
type ImportStack []string

func New() *ImportStack {
	return &ImportStack{}
}

func (s *ImportStack) Push(p string) error {
	for _, x := range *s {
		if x == p {
			return ErrImportCycleNotAllowed
		}
	}
	*s = append(*s, p)
	return nil
}

func (s *ImportStack) Pop() {
	*s = (*s)[0 : len(*s)-1]
}

func (s *ImportStack) List() []string {
	return *s
}
