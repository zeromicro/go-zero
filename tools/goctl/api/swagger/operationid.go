package swagger

import "fmt"

func generateDefaultOperationId(groupName string, handlerName string) string {
	if groupName == "" {
		return handlerName
	}

	return fmt.Sprintf("%s_%s", groupName, handlerName)
}
