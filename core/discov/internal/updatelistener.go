//go:generate mockgen -package internal -destination updatelistener_mock.go -source updatelistener.go UpdateListener

package internal

type (
	// A KV is used to store an etcd entry with key and value.
	KV struct {
		Key string
		Val string
	}

	// UpdateListener wraps the OnAdd and OnDelete methods.
	UpdateListener interface {
		OnAdd(kv KV)
		OnDelete(kv KV)
	}
)
