package model

var (
	CommonMysqlDataTypeMap = map[string]string{
		"tinyint":    "int",
		"smallint":   "int",
		"mediumint":  "int64",
		"int":        "int64",
		"integer":    "int64",
		"bigint":     "int64",
		"float":      "float32",
		"double":     "float64",
		"decimal":    "float64",
		"date":       "time.Time",
		"time":       "string",
		"year":       "int",
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
	ModeDirPerm = 0755
)
