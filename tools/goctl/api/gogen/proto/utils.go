package proto

func convertProtoTypeToGoType(typeName string) string {
	switch typeName {
	case "float":
		typeName = "float32"
	case "double":
		typeName = "float64"
	}
	return typeName
}
