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
		Info    Info
		Syntax  ApiSyntax
		Imports []Import
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
		Docs               Doc
		IsInline           bool
		IsRequired         bool
		NotAllowEmptyValue bool
	}

	// Route describes api route
	Route struct {
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
		Nullable() bool
	}

	// DefineStruct describes api structure
	DefineStruct struct {
		BaseType
		Comment  string
		Members  []Member
		Docs     Doc
		RawName  string
		Required []string
	}

	// PrimitiveType describes the basic golang type, such as bool,int32,int64, ...
	PrimitiveType struct {
		BaseType

		RawName string
		Comment string

		//number
		Min          *float64
		Max          *float64
		MultipleOf   *float64
		ExclusiveMin bool
		ExclusiveMax bool

		//string
		MinLength uint64
		MaxLength *uint64
		Pattern   string
		Format    string //

		Enum []interface{}
	}

	// MapType describes a map for api
	MapType struct {
		BaseType

		Comment string
		RawName string
		// only support the PrimitiveType
		Key string
		// it can be asserted as PrimitiveType: int、bool、
		// PointerType: *string、*User、
		// MapType: map[${PrimitiveType}]interface、
		// ArrayType:[]int、[]User、[]*User
		// InterfaceType: interface{}
		// Type
		Value       Type
		MinItems    int64
		MaxItems    int64
		UniqueItems bool
	}

	// ArrayType describes a slice for api
	ArrayType struct {
		BaseType

		Comment     string
		Value       Type
		RawName     string
		MinItems    uint64
		MaxItems    *uint64
		UniqueItems bool
	}

	// InterfaceType describes an interface for api
	InterfaceType struct {
		BaseType

		RawName string
		RawID   string
	}

	// PointerType describes a pointer for api
	PointerType struct {
		BaseType

		RawName    string
		Type       Type
		RawID      string
		RefRawName string
	}

	// AtDoc describes a metadata for api grammar: @doc(...)
	AtDoc struct {
		Properties map[string]string
		Text       string
	}

	BaseType struct {
		RawNullable bool
	}
)
