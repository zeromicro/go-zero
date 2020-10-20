package parser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
	"github.com/tal-tech/go-zero/core/collection"
	"github.com/tal-tech/go-zero/core/lang"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

const (
	AnyImport = "google/protobuf/any.proto"
)

var (
	enumTypeTemplate = `{{.name}} int32`
	enumTemplate     = `const (
	{{.element}}
)`
	enumFiledTemplate = `{{.key}} {{.name}} = {{.value}}`
)

type (
	MessageField struct {
		Type string
		Name stringx.String
	}
	Message struct {
		Name    stringx.String
		Element []*MessageField
		*proto.Message
	}
	Enum struct {
		Name    stringx.String
		Element []*EnumField
		*proto.Enum
	}
	EnumField struct {
		Key   string
		Value int
	}

	Proto struct {
		Package string
		Import  []*Import
		PbSrc   string
		// deprecated: containsAny will be removed in the feature
		ContainsAny bool
		Message     map[string]lang.PlaceholderType
		Enum        map[string]*Enum
	}
	Import struct {
		ProtoImportName   string
		PbImportName      string
		OriginalDir       string
		OriginalProtoPath string
		OriginalPbPath    string
		BridgeImport      string
		exists            bool
		//xx.proto
		protoName string
		// xx.pb.go
		pbName string
	}
)

func checkImport(src string) error {
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	parser := proto.NewParser(r)
	parseRet, err := parser.Parse()
	if err != nil {
		return err
	}
	var base = filepath.Base(src)
	proto.Walk(parseRet, proto.WithImport(func(i *proto.Import) {
		if err != nil {
			return
		}
		err = fmt.Errorf("%v:%v the external proto cannot import other proto files", base, i.Position.Line)
	}))
	if err != nil {
		return err
	}
	return nil
}
func ParseImport(src string) ([]*Import, bool, error) {
	bridgeImportM := make(map[string]string)
	r, err := os.Open(src)
	if err != nil {
		return nil, false, err
	}
	defer r.Close()

	workDir := filepath.Dir(src)
	parser := proto.NewParser(r)
	parseRet, err := parser.Parse()
	if err != nil {
		return nil, false, err
	}
	protoImportSet := collection.NewSet()
	var containsAny bool
	proto.Walk(parseRet, proto.WithImport(func(i *proto.Import) {
		if i.Filename == AnyImport {
			containsAny = true
			return
		}
		protoImportSet.AddStr(i.Filename)
		if i.Comment != nil {
			lines := i.Comment.Lines
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if !strings.HasPrefix(line, "@") {
					continue
				}
				line = strings.TrimPrefix(line, "@")
				bridgeImportM[i.Filename] = line
			}
		}
	}))
	var importList []*Import

	for _, item := range protoImportSet.KeysStr() {
		pb := strings.TrimSuffix(filepath.Base(item), filepath.Ext(item)) + ".pb.go"
		var pbImportName, brideImport string
		if v, ok := bridgeImportM[item]; ok {
			pbImportName = v
			brideImport = "M" + item + "=" + v
		} else {
			pbImportName = item
		}
		var impo = Import{
			ProtoImportName: item,
			PbImportName:    pbImportName,
			BridgeImport:    brideImport,
		}
		protoSource := filepath.Join(workDir, item)
		pbSource := filepath.Join(filepath.Dir(protoSource), pb)
		if util.FileExists(protoSource) && util.FileExists(pbSource) {
			impo.OriginalProtoPath = protoSource
			impo.OriginalPbPath = pbSource
			impo.OriginalDir = filepath.Dir(protoSource)
			impo.exists = true
			impo.protoName = filepath.Base(item)
			impo.pbName = pb
		} else {
			return nil, false, fmt.Errorf("「%v」: import must be found in the relative directory of 「%v」", item, filepath.Base(src))
		}
		importList = append(importList, &impo)
	}

	return importList, containsAny, nil
}

func parseProto(src string, messageM map[string]lang.PlaceholderType, enumM map[string]*Enum) (*Proto, error) {
	if !filepath.IsAbs(src) {
		return nil, fmt.Errorf("expected absolute path,but found: %v", src)
	}

	r, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	parser := proto.NewParser(r)
	parseRet, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	// xx.proto
	fileBase := filepath.Base(src)
	var resp Proto

	proto.Walk(parseRet, proto.WithPackage(func(p *proto.Package) {
		if err != nil {
			return
		}

		if len(resp.Package) != 0 {
			err = fmt.Errorf("%v:%v duplicate package「%v」", fileBase, p.Position.Line, p.Name)
		}

		if len(p.Name) == 0 {
			err = errors.New("package not found")
		}

		resp.Package = p.Name
	}), proto.WithMessage(func(message *proto.Message) {
		if err != nil {
			return
		}

		for _, item := range message.Elements {
			switch item.(type) {
			case *proto.NormalField, *proto.MapField, *proto.Comment:
				continue
			default:
				err = fmt.Errorf("%v: unsupport inline declaration", fileBase)
				return
			}
		}
		name := stringx.From(message.Name)
		if _, ok := messageM[name.Lower()]; ok {
			err = fmt.Errorf("%v:%v duplicate message 「%v」", fileBase, message.Position.Line, message.Name)
			return
		}

		messageM[name.Lower()] = lang.Placeholder
	}), proto.WithEnum(func(enum *proto.Enum) {
		if err != nil {
			return
		}

		var node Enum
		node.Enum = enum
		node.Name = stringx.From(enum.Name)
		for _, item := range enum.Elements {
			v, ok := item.(*proto.EnumField)
			if !ok {
				continue
			}
			node.Element = append(node.Element, &EnumField{
				Key:   v.Name,
				Value: v.Integer,
			})
		}
		if _, ok := enumM[node.Name.Lower()]; ok {
			err = fmt.Errorf("%v:%v duplicate enum 「%v」", fileBase, node.Position.Line, node.Name.Source())
			return
		}

		lower := stringx.From(enum.Name).Lower()
		enumM[lower] = &node
	}))

	if err != nil {
		return nil, err
	}
	resp.Message = messageM
	resp.Enum = enumM

	return &resp, nil
}

func (e *Enum) GenEnumCode() (string, error) {
	var element []string
	for _, item := range e.Element {
		code, err := item.GenEnumFieldCode(e.Name.Source())
		if err != nil {
			return "", err
		}
		element = append(element, code)
	}
	buffer, err := util.With("enum").Parse(enumTemplate).Execute(map[string]interface{}{
		"element": strings.Join(element, util.NL),
	})
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (e *Enum) GenEnumTypeCode() (string, error) {
	buffer, err := util.With("enumAlias").Parse(enumTypeTemplate).Execute(map[string]interface{}{
		"name": e.Name.Source(),
	})
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (e *EnumField) GenEnumFieldCode(parentName string) (string, error) {
	buffer, err := util.With("enumField").Parse(enumFiledTemplate).Execute(map[string]interface{}{
		"key":   e.Key,
		"name":  parentName,
		"value": e.Value,
	})
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
