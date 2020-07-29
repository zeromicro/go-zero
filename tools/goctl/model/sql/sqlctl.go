package sqlgen

import (
	"errors"

	"zero/tools/goctl/model/sql/gen"
)

var (
	ErrCircleQuery = errors.New("circle query with other fields")
)

func Gen(in gen.OuterTable) (string, error) {
	t, err := gen.TableConvert(in)
	if err != nil {
		if err == gen.ErrCircleQuery {
			return "", ErrCircleQuery
		}
		return "", err
	}
	return gen.GenModel(t)
}
