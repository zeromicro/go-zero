package sqlmodel

import "errors"

func sqlError(str string) error {
	return errors.New("sql error: " + str)
}
