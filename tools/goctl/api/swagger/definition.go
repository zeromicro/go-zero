package swagger

import (
	"github.com/go-openapi/spec"
	apiSpec "github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func definitionsFromTypes(ctx Context, types []apiSpec.Type) spec.Definitions {
	if !ctx.UseDefinitions {
		return nil
	}
	definitions := make(spec.Definitions)
	for _, tp := range types {
		typeName := tp.Name()
		definitions[typeName] = schemaFromType(ctx, tp)
	}
	return definitions
}

func schemaFromType(ctx Context, tp apiSpec.Type) spec.Schema {
	p, r := propertiesFromType(ctx, tp)
	props := spec.SchemaProps{
		Type:                 typeFromGoType(ctx, tp),
		Properties:           p,
		AdditionalProperties: mapFromGoType(ctx, tp),
		Items:                itemsFromGoType(ctx, tp),
		Required:             r,
	}
	return spec.Schema{
		SchemaProps: props,
	}
}
