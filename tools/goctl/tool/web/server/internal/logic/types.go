package logic

type (
	RouteGroup struct {
		Jwt        bool
		Prefix     string
		Group      string
		Timeout    int
		Middleware string
		MaxBytes   int64
	}

	APIRoute struct {
		Handler     string
		Method      string
		Path        string
		ContentType string
	}

	json2APITypeReq struct {
		root           bool
		indentCount    int
		parentTypeName string
		typeName       string
		v              any
	}

	goctlAPIMemberResult struct {
		TypeExpr         string
		TypeName         string
		IsStruct         bool
		IsArray          bool
		ExternalTypeExpr []string
	}
	KV map[string]any
)
