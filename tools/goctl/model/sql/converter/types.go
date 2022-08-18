package converter

import (
	"fmt"
	"strings"

	"github.com/zeromicro/ddl-parser/parser"
)

var unsignedTypeMap = map[string]string{
	"int":   "uint",
	"int8":  "uint8",
	"int16": "uint16",
	"int32": "uint32",
	"int64": "uint64",
}

// Inspired by https://github.com/go-sql-driver/mysql/blob/master/fields.go#L140
var commonMysqlDataTypeMapInt = map[int]string{
	// For consistency, all integer types are converted to int64
	// number
	parser.Bit:       "sql.RawBytes",
	parser.TinyInt:   "int8",
	parser.SmallInt:  "int16",
	parser.MediumInt: "int32",
	parser.Int:       "int32",
	parser.MiddleInt: "int32",
	parser.Int1:      "int8",
	parser.Int2:      "int16",
	parser.Int3:      "int32",
	parser.Int4:      "int32",
	parser.Int8:      "int64",
	parser.Integer:   "int32",
	parser.BigInt:    "int64",
	parser.Float:     "float32",
	parser.Float4:    "float32",
	parser.Float8:    "float64",
	parser.Double:    "float64",
	parser.Decimal:   "sql.RawBytes",
	parser.Dec:       "sql.RawBytes",
	parser.Fixed:     "sql.RawBytes",
	parser.Numeric:   "sql.RawBytes",
	parser.Real:      "sql.RawBytes",
	// date&time
	parser.Date:      "time.Time",
	parser.DateTime:  "time.Time",
	parser.Timestamp: "time.Time",
	parser.Time:      "sql.RawBytes",
	parser.Year:      "int16",
	// string
	parser.Char:            "string",
	parser.VarChar:         "string",
	parser.NVarChar:        "string",
	parser.NChar:           "string",
	parser.Character:       "string",
	parser.LongVarChar:     "string",
	parser.LineString:      "string",
	parser.MultiLineString: "string",
	parser.Binary:          "string",
	parser.VarBinary:       "string",
	parser.TinyText:        "sql.RawBytes",
	parser.Text:            "sql.RawBytes",
	parser.MediumText:      "sql.RawBytes",
	parser.LongText:        "sql.RawBytes",
	parser.Enum:            "sql.RawBytes",
	parser.Set:             "sql.RawBytes",
	parser.Json:            "sql.RawBytes",
	parser.Blob:            "sql.RawBytes",
	parser.LongBlob:        "sql.RawBytes",
	parser.MediumBlob:      "sql.RawBytes",
	parser.TinyBlob:        "sql.RawBytes",
	// bool
	parser.Bool:    "bool",
	parser.Boolean: "bool",
}

// Inspired by https://github.com/go-sql-driver/mysql/blob/master/fields.go#L140
var commonMysqlDataTypeMapString = map[string]string{
	// For consistency, all integer types are converted to int64
	// bool
	"bool":    "bool",
	"boolean": "bool",
	// number
	"tinyint":   "int8",
	"smallint":  "int16",
	"mediumint": "int32",
	"int":       "int32",
	"int1":      "int8",
	"int2":      "int16",
	"int3":      "int32",
	"int4":      "int32",
	"int8":      "int64",
	"integer":   "int32",
	"bigint":    "int64",
	"float":     "float32",
	"float4":    "float32",
	"float8":    "float64",
	"double":    "float64",
	"decimal":   "sql.RawBytes",
	"dec":       "sql.RawBytes",
	"fixed":     "sql.RawBytes",
	"real":      "sql.RawBytes",
	"bit":       "sql.RawBytes",
	// date & time
	"date":      "time.Time",
	"datetime":  "time.Time",
	"timestamp": "time.Time",
	"time":      "sql.RawBytes",
	"year":      "int16",
	// string
	"linestring":      "string",
	"multilinestring": "string",
	"nvarchar":        "string",
	"nchar":           "string",
	"char":            "string",
	"character":       "string",
	"varchar":         "string",
	"binary":          "string",
	"bytea":           "string",
	"longvarbinary":   "string",
	"varbinary":       "string",
	"tinytext":        "sql.RawBytes",
	"text":            "sql.RawBytes",
	"mediumtext":      "sql.RawBytes",
	"longtext":        "sql.RawBytes",
	"enum":            "sql.RawBytes",
	"set":             "sql.RawBytes",
	"json":            "sql.RawBytes",
	"jsonb":           "sql.RawBytes",
	"blob":            "sql.RawBytes",
	"longblob":        "sql.RawBytes",
	"mediumblob":      "sql.RawBytes",
	"tinyblob":        "sql.RawBytes",
}

// ConvertDataType converts mysql column type into golang type
func ConvertDataType(dataBaseType int, isDefaultNull, unsigned bool) (string, error) {
	tp, ok := commonMysqlDataTypeMapInt[dataBaseType]
	if !ok {
		return "", fmt.Errorf("unsupported database type: %v", dataBaseType)
	}

	return mayConvertNullType(tp, isDefaultNull, unsigned), nil
}

// ConvertStringDataType converts mysql column type into golang type
func ConvertStringDataType(dataBaseType string, isDefaultNull, unsigned bool) (string, error) {
	tp, ok := commonMysqlDataTypeMapString[strings.ToLower(dataBaseType)]
	if !ok {
		return "", fmt.Errorf("unsupported database type: %s", dataBaseType)
	}

	return mayConvertNullType(tp, isDefaultNull, unsigned), nil
}

func mayConvertNullType(goDataType string, isDefaultNull, unsigned bool) string {
	if !isDefaultNull {
		if unsigned {
			ret, ok := unsignedTypeMap[goDataType]
			if ok {
				return ret
			}
		}
		return goDataType
	}

	switch goDataType {
	case "int8", "int16", "int32", "int64",
		"uint8", "uint16", "uint32", "uint64":
		return "sql.NullInt64"
	case "float32", "float64":
		return "sql.NullFloat64"
	case "bool":
		return "sql.NullBool"
	case "string":
		return "sql.NullString"
	case "time.Time":
		return "sql.NullTime"
	default:
		if unsigned {
			ret, ok := unsignedTypeMap[goDataType]
			if ok {
				return ret
			}
		}
		return goDataType
	}
}
