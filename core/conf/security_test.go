package conf

import (
	"os"
	"testing"

	"github.com/WqyJh/confcrypt"
	"github.com/stretchr/testify/assert"
)

func TestSecurityLoad(t *testing.T) {
	key := "testkey"
	type testConfig struct {
		SecurityConf

		User   string `json:"user"`
		Pass   string `json:"pass"`
		Secret string `json:"secret"`
	}
	expected := testConfig{
		SecurityConf: SecurityConf{
			Enable: true,
			Env:    "CONFIG_KEY",
		},
		User:   "testuser",
		Pass:   "testpass",
		Secret: "testsecret",
	}
	encryptedPass, err := confcrypt.EncryptString(expected.Pass, key)
	assert.Nil(t, err)
	encryptedSecret, err := confcrypt.EncryptString(expected.Secret, key)
	assert.Nil(t, err)
	text := `{
		"user": "testuser",
		"pass": "` + encryptedPass + `",
		"secret": "` + encryptedSecret + `"
}`
	tmpfile, err := createTempFile(".json", text)
	assert.Nil(t, err)
	defer os.Remove(tmpfile)

	os.Setenv("CONFIG_KEY", key)
	var config testConfig
	err = SecurityLoad(tmpfile, &config)
	assert.Nil(t, err)
	assert.NotEqual(t, encryptedPass, config.Pass)
	assert.NotEqual(t, encryptedPass, config.Secret)
	assert.Equal(t, expected.Pass, config.Pass)
	assert.Equal(t, expected.Secret, config.Secret)
}
