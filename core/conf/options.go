package conf

type (
	Option func(opt *options)

	options struct {
		env bool
	}
)

func UseEnv() Option {
	return func(opt *options) {
		opt.env = true
	}
}
