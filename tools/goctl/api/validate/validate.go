package validate

import (
	"errors"
	"fmt"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
)

// VarStringAPI describes an API.
var VarStringAPI string

// GoValidateApi verifies whether the api has a syntax error
func GoValidateApi(_ *cobra.Command, _ []string) error {
	apiFile := VarStringAPI

	if len(apiFile) == 0 {
		return errors.New("missing -api")
	}

	spec, err := parser.Parse(apiFile)
	if err != nil {
		return err
	}

	err = spec.Validate()
	if err == nil {
		fmt.Println(color.Green.Render("api format ok"))
	}
	return err
}
