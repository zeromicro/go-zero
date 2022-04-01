package template

import (
	"fmt"

	"github.com/zeromicro/go-zero/tools/goctl/util"
)

// ModelGen defines a template for model
var ModelGen = fmt.Sprintf(`%s

package {{.pkg}}
{{.imports}}
{{.vars}}
{{.types}}
{{.new}}
{{.insert}}
{{.find}}
{{.update}}
{{.delete}}
{{.extraMethod}}
`, util.DoNotEditHead)

// ModelCustom defines a template for extension
var ModelCustom = fmt.Sprintf(`package {{.pkg}}
type {{.upperStartCamelObject}}Model interface {
	{{.lowerStartCamelObject}}Model
}
`)
