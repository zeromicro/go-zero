package discov

type Publisher interface {
	SetID(int64)
	KeepAlive() error
	Pause()
	Resume()
	Stop()
}
