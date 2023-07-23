package queue

import "strings"

func generateName(pushers []Pusher) string {
	names := make([]string, len(pushers))
	for i, pusher := range pushers {
		names[i] = pusher.Name()
	}

	return strings.Join(names, ",")
}
