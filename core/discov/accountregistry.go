package discov

import "github.com/tal-tech/go-zero/core/discov/internal"

func RegisterAccount(endpoints []string, user, pass string) {
	internal.AddAccount(endpoints, user, pass)
}
