package ast

import "errors"

// ErrImportCycleNotAllowed defines an error for circular importing
var ErrImportCycleNotAllowed = errors.New("import cycle not allowed")

// importStack a stack of import paths
type importStack []string

func (s *importStack) push(p string) error {
	for _, x := range *s {
		if x == p {
			return ErrImportCycleNotAllowed
		}
	}
	*s = append(*s, p)
	return nil
}

func (s *importStack) pop() {
	*s = (*s)[0 : len(*s)-1]
}
