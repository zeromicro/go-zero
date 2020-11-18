package generator

type NamingStyle = string

const (
	namingLower NamingStyle = "lower"
	namingCamel NamingStyle = "camel"
	namingSnake NamingStyle = "snake"
)

// IsNamingValid validates whether the namingStyle is valid or not,return
// namingStyle and true if it is valid, or else return empty string
// and false, and it is a valid value even namingStyle is empty string
func IsNamingValid(namingStyle string) (NamingStyle, bool) {
	if len(namingStyle) == 0 {
		namingStyle = namingLower
	}
	switch namingStyle {
	case namingLower, namingCamel, namingSnake:
		return namingStyle, true
	default:
		return "", false
	}
}
