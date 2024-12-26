package template

import (
	_ "embed"
	"fmt"

	"github.com/zeromicro/go-zero/tools/goctl/internal/version"
	"github.com/zeromicro/go-zero/tools/goctl/util"
)

// Customized defines a template for customized in model
//
//go:embed tpl/customized.tpl
var Customized string

// Vars defines a template for var block in model
//
//go:embed tpl/var.tpl
var Vars string

// Types defines a template for types in model.
//
//go:embed tpl/types.tpl
var Types string

// Tag defines a tag template text
//
//go:embed tpl/tag.tpl
var Tag string

// TableName defines a template that generate the tableName method.
//
//go:embed tpl/table-name.tpl
var TableName string

// New defines the template for creating model instance.
//
//go:embed tpl/model-new.tpl
var New string

// ModelCustom defines a template for extension
//
//go:embed tpl/model.tpl
var ModelCustom string

// ModelGen defines a template for model
var ModelGen = fmt.Sprintf(`%s
// versions:
//  goctl version: %s

package {{.pkg}}
{{.imports}}
{{.vars}}
{{.types}}
{{.new}}
{{.delete}}
{{.find}}
{{.insert}}
{{.update}}
{{.extraMethod}}
{{.tableName}}
{{.customized}}
`, util.DoNotEditHead, version.BuildVersion)

// Insert defines a template for insert code in model
//
//go:embed tpl/insert.tpl
var Insert string

// InsertMethod defines an interface method template for insert code in model
//
//go:embed tpl/interface-insert.tpl
var InsertMethod string

// Update defines a template for generating update codes
//
//go:embed tpl/update.tpl
var Update string

// UpdateMethod defines an interface method template for generating update codes
//
//go:embed tpl/interface-update.tpl
var UpdateMethod string

// Imports defines a import template for model in cache case
//
//go:embed tpl/import.tpl
var Imports string

// ImportsNoCache defines a import template for model in normal case
//
//go:embed tpl/import-no-cache.tpl
var ImportsNoCache string

// FindOne defines find row by id.
//
//go:embed tpl/find-one.tpl
var FindOne string

// FindOneByField defines find row by field.
//
//go:embed tpl/find-one-by-field.tpl
var FindOneByField string

// FindOneByFieldExtraMethod defines find row by field with extras.
//
//go:embed tpl/find-one-by-field-extra-method.tpl
var FindOneByFieldExtraMethod string

// FindOneMethod defines find row method.
//
//go:embed tpl/interface-find-one.tpl
var FindOneMethod string

// FindOneByFieldMethod defines find row by field method.
//
//go:embed tpl/interface-find-one-by-field.tpl
var FindOneByFieldMethod string

// Field defines a filed template for types
//
//go:embed tpl/field.tpl
var Field string

// Error defines an error template
//
//go:embed tpl/err.tpl
var Error string

// Delete defines a delete template
//
//go:embed tpl/delete.tpl
var Delete string

// DeleteMethod defines a delete template for interface method
//
//go:embed tpl/interface-delete.tpl
var DeleteMethod string
