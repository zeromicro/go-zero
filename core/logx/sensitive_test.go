package logx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const maskedContent = "******"

type User struct {
	Name string
	Pass string
}

func (u User) MaskSensitive() any {
	return User{
		Name: u.Name,
		Pass: maskedContent,
	}
}

type NonSensitiveUser struct {
	Name string
	Pass string
}

func TestMaskSensitive(t *testing.T) {
	t.Run("sensitive", func(t *testing.T) {
		user := User{
			Name: "kevin",
			Pass: "123",
		}

		mu := maskSensitive(user)
		assert.Equal(t, user.Name, mu.(User).Name)
		assert.Equal(t, maskedContent, mu.(User).Pass)
	})

	t.Run("non-sensitive", func(t *testing.T) {
		user := NonSensitiveUser{
			Name: "kevin",
			Pass: "123",
		}

		mu := maskSensitive(user)
		assert.Equal(t, user.Name, mu.(NonSensitiveUser).Name)
		assert.Equal(t, user.Pass, mu.(NonSensitiveUser).Pass)
	})
}
