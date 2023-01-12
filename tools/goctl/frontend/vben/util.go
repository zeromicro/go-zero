package vben

func ConvertGoTypeToTsType(goType string) string {
	switch goType {
	case "int", "uint", "int32", "int64", "uint32", "uint64", "float", "float32", "float64":
		goType = "number"
	case "[]int", "[]uint", "[]int32", "[]int64", "[]uint32", "[]uint64", "[]float", "[]float32", "[]float64":
		goType = "number[]"
	case "string":
		goType = "string"
	case "[]string":
		goType = "string[]"
	case "bool":
		goType = "boolean"

	}
	return goType
}
