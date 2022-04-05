package ast

import "errors"

var ErrImportCycleNotAllowed = errors.New("import cycle not allowed")

// ImportStack a stack of import paths
type ImportStack []string

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

func (s *ImportStack) Copy() []string {
	return append([]string{}, *s...)
}

func (s *ImportStack) Top() string {
	if len(*s) == 0 {
		return ""
	}
	return (*s)[len(*s)-1]
}
