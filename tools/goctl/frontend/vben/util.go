package vben

func ConvertGoTypeToTsType(goType string) string {
	switch goType {
	case "int", "uint", "int64", "uint64", "float", "float32", "float64":
		goType = "number"
	case "[]int", "[]uint", "[]int64", "[]uint64", "[]float", "[]float32", "[]float64":
		goType = "number[]"
	case "string":
		goType = "string"
	case "[]string":
		goType = "string[]"

	}
	return goType
}
