//go:generate mockgen -package internal -destination updatelistener_mock.go -source updatelistener.go UpdateListener
package internal

type (
	KV struct {
		Key string
		Val string
	}

	UpdateListener interface {
		OnAdd(kv KV)
		OnDelete(kv KV)
	}
)
