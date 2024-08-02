package spec

// RoutePrefixKey is the prefix keyword for the routes.
const RoutePrefixKey = "prefix"

type (
	// Doc describes document
	Doc []string

	// Annotation defines key-value
	Annotation struct {
		Properties map[string]string
	}

	// ApiSyntax describes the syntax grammar
	ApiSyntax struct {
		Version string
		Doc     Doc
		Comment Doc
	}

	// ApiSpec describes an api file
	ApiSpec struct {
		Info    Info      // Deprecated: useless expression
		Syntax  ApiSyntax // Deprecated: useless expression
		Imports []Import  // Deprecated: useless expression
		Types   []Type
		Service Service
	}

	// Import describes api import
	Import struct {
		Value   string
		Doc     Doc
		Comment Doc
	}

	// Group defines a set of routing information
	Group struct {
		Annotation Annotation
		Routes     []Route
	}

	// Info describes info grammar block
	Info struct {
		// Deprecated: use Properties instead
		Title string
		// Deprecated: use Properties instead
		Desc string
		// Deprecated: use Properties instead
		Version string
		// Deprecated: use Properties instead
		Author string
		// Deprecated: use Properties instead
		Email      string
		Properties map[string]string
	}

	// Member describes the field of a structure
	Member struct {
		Name string
		// 数据类型字面值，如：string、map[int]string、[]int64、[]*User
		Type    Type
		Tag     string
		Comment string
		// 成员头顶注释说明
		Docs     Doc
		IsInline bool
	}

	// Route describes api route
	Route struct {
		// Deprecated: Use Service AtServer instead.
		AtServerAnnotation Annotation
		Method             string
		Path               string
		RequestType        Type
		ResponseType       Type
		Docs               Doc
		Handler            string
		AtDoc              AtDoc
		HandlerDoc         Doc
		HandlerComment     Doc
		Doc                Doc
		Comment            Doc
	}

	// Service describes api service
	Service struct {
		Name   string
		Groups []Group
	}

	// Type defines api type
	Type interface {
		Name() string
		Comments() []string
		Documents() []string
	}

	// DefineStruct describes api structure
	DefineStruct struct {
		RawName string
		Members []Member
		Docs    Doc
	}

	// NestedStruct describes a structure nested in structure.
	NestedStruct struct {
		RawName string
		Members []Member
		Docs    Doc
	}

	// PrimitiveType describes the basic golang type, such as bool,int32,int64, ...
	PrimitiveType struct {
		RawName string
	}

	// MapType describes a map for api
	MapType struct {
		RawName string
		// only support the PrimitiveType
		Key string
		// it can be asserted as PrimitiveType: int、bool、
		// PointerType: *string、*User、
		// MapType: map[${PrimitiveType}]interface、
		// ArrayType:[]int、[]User、[]*User
		// InterfaceType: interface{}
		// Type
		Value Type
	}

	// ArrayType describes a slice for api
	ArrayType struct {
		RawName string
		Value   Type
	}

	// InterfaceType describes an interface for api
	InterfaceType struct {
		RawName string
	}

	// PointerType describes a pointer for api
	PointerType struct {
		RawName string
		Type    Type
	}

	// AtDoc describes a metadata for api grammar: @doc(...)
	AtDoc struct {
		Properties map[string]string
		Text       string
	}
)
