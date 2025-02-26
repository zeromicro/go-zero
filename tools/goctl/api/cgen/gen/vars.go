package gen

import "github.com/zeromicro/go-zero/tools/goctl/api/spec"

var AllTypes []spec.Type

func findTypeByName(name string) spec.Type {
	for _, t := range AllTypes {
		if t.Name() == name {
			return t
		}
	}
	return nil
}
