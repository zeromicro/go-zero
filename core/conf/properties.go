package conf

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/tal-tech/go-zero/core/iox"
)

// PropertyError represents a configuration error message.
type PropertyError struct {
	error
	message string
}

// Properties interface provides the means to access configuration.
type Properties interface {
	GetString(key string) string
	SetString(key, value string)
	GetInt(key string) int
	SetInt(key string, value int)
	ToString() string
}

// Properties config is a key/value pair based configuration structure.
type mapBasedProperties struct {
	properties map[string]string
	lock       sync.RWMutex
}

// Loads the properties into a properties configuration instance. May return the
// configuration itself along with an error that indicates if there was a problem loading the configuration.
func LoadProperties(filename string) (Properties, error) {
	lines, err := iox.ReadTextLines(filename, iox.WithoutBlank(), iox.OmitWithPrefix("#"))
	if err != nil {
		return nil, nil
	}

	raw := make(map[string]string)
	for i := range lines {
		pair := strings.Split(lines[i], "=")
		if len(pair) != 2 {
			// invalid property format
			return nil, &PropertyError{
				message: fmt.Sprintf("invalid property format: %s", pair),
			}
		}

		key := strings.TrimSpace(pair[0])
		value := strings.TrimSpace(pair[1])
		raw[key] = value
	}

	return &mapBasedProperties{
		properties: raw,
	}, nil
}

func (config *mapBasedProperties) GetString(key string) string {
	config.lock.RLock()
	ret := config.properties[key]
	config.lock.RUnlock()

	return ret
}

func (config *mapBasedProperties) SetString(key, value string) {
	config.lock.Lock()
	config.properties[key] = value
	config.lock.Unlock()
}

func (config *mapBasedProperties) GetInt(key string) int {
	config.lock.RLock()
	// default 0
	value, _ := strconv.Atoi(config.properties[key])
	config.lock.RUnlock()

	return value
}

func (config *mapBasedProperties) SetInt(key string, value int) {
	config.lock.Lock()
	config.properties[key] = strconv.Itoa(value)
	config.lock.Unlock()
}

// Dumps the configuration internal map into a string.
func (config *mapBasedProperties) ToString() string {
	config.lock.RLock()
	ret := fmt.Sprintf("%s", config.properties)
	config.lock.RUnlock()

	return ret
}

// Returns the error message.
func (configError *PropertyError) Error() string {
	return configError.message
}

// Builds a new properties configuration structure
func NewProperties() Properties {
	return &mapBasedProperties{
		properties: make(map[string]string),
	}
}
