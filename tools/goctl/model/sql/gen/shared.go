package gen

var (
	commonMysqlDataTypeMap = map[string]string{
		"tinyint":    "int64",
		"smallint":   "int64",
		"mediumint":  "int64",
		"int":        "int64",
		"integer":    "int64",
		"bigint":     "int64",
		"float":      "float64",
		"double":     "float64",
		"decimal":    "float64",
		"date":       "time.Time",
		"time":       "string",
		"year":       "int64",
		"datetime":   "time.Time",
		"timestamp":  "time.Time",
		"char":       "string",
		"varchar":    "string",
		"tinyblob":   "string",
		"tinytext":   "string",
		"blob":       "string",
		"text":       "string",
		"mediumblob": "string",
		"mediumtext": "string",
		"longblob":   "string",
		"longtext":   "string",
	}
)

const (
	QueryNone  QueryType = 0
	QueryOne   QueryType = 1 // 仅支持单个字段为查询条件
	QueryAll   QueryType = 2 // 可支持多个字段为查询条件，且关系均为and
	QueryLimit QueryType = 3 // 可支持多个字段为查询条件，且关系均为and
)

type (
	QueryType int

	Case struct {
		SnakeCase      string
		LowerCamelCase string
		UpperCamelCase string
	}
	InnerWithField struct {
		Case
		DataType string
	}
	InnerTable struct {
		Case
		ContainsCache  bool
		CreateNotFound bool
		PrimaryField   *InnerField
		Fields         []*InnerField
		CacheKey       map[string]Key // key-数据库字段
	}
	InnerField struct {
		IsPrimaryKey bool
		InnerWithField
		DataBaseType string // 数据库中字段类型
		Tag          string // 标签，格式：`db:"xxx"`
		Comment      string // 注释,以"// 开头"
		Cache        bool   // 是否缓存模式
		QueryType    QueryType
		WithFields   []InnerWithField
		Sort         []InnerSort
	}
	InnerSort struct {
		Field Case
		Asc   bool
	}

	OuterTable struct {
		Table          string        `json:"table"`
		CreateNotFound bool          `json:"createNotFound,optional"`
		Fields         []*OuterFiled `json:"fields"`
	}
	OuterWithField struct {
		Name         string `json:"name"`
		DataBaseType string `json:"dataBaseType"`
	}
	OuterSort struct {
		Field string `json:"fields"`
		Asc   bool   `json:"asc,optional"`
	}
	OuterFiled struct {
		IsPrimaryKey bool   `json:"isPrimaryKey,optional"`
		Name         string `json:"name"`
		DataBaseType string `json:"dataBaseType"`
		Comment      string `json:"comment"`
		Cache        bool   `json:"cache,optional"`
		// if IsPrimaryKey==false下面字段有效
		QueryType  QueryType        `json:"queryType"`           // 查找类型
		WithFields []OuterWithField `json:"withFields,optional"` // 其他字段联合组成条件的字段列表
		OuterSort  []OuterSort      `json:"sort,optional"`
	}
)
