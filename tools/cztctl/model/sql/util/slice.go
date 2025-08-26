package util

// TrimStringSlice returns a copy slice without empty string item
func TrimStringSlice(list []string) []string {
	var out []string
	for _, item := range list {
		if len(item) == 0 {
			continue
		}
		out = append(out, item)
	}
	return out
}
