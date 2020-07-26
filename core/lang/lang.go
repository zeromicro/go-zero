package lang

import "log"

var Placeholder PlaceholderType

type (
	GenericType     = interface{}
	PlaceholderType = struct{}
)

func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
