package discov

type (
	subOptions struct {
		exclusive bool
	}

	SubOption func(opts *subOptions)
)

type Subscriber interface {
	AddListener(listener func())
	Values() []string
}
