package converter

import (
	"fmt"
	"strings"
)

var (
	commonMysqlDataTypeMap = map[string]string{
		// For consistency, all integer types are converted to int64
		// number
		"bool":      "int64",
		"boolean":   "int64",
		"tinyint":   "int64",
		"smallint":  "int64",
		"mediumint": "int64",
		"int":       "int64",
		"integer":   "int64",
		"bigint":    "int64",
		"float":     "float64",
		"double":    "float64",
		"decimal":   "float64",
		// date&time
		"date":      "time.Time",
		"datetime":  "time.Time",
		"timestamp": "time.Time",
		"time":      "string",
		"year":      "int64",
		// string
		"char":       "string",
		"varchar":    "string",
		"binary":     "string",
		"varbinary":  "string",
		"tinytext":   "string",
		"text":       "string",
		"mediumtext": "string",
		"longtext":   "string",
		"enum":       "string",
		"set":        "string",
		"json":       "string",
	}
)

func ConvertDataType(dataBaseType string) (goDataType string, err error) {
	tp, ok := commonMysqlDataTypeMap[strings.ToLower(dataBaseType)]
	if !ok {
		err = fmt.Errorf("unexpected database type: %s", dataBaseType)
		return
	}
	goDataType = tp
	return
}
