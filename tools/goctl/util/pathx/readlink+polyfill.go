//go:build windows

package pathx

func ReadLink(name string) (string, error) {
	return name, nil
}
