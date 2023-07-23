package conf

type (
	// Option defines the method to customize the config options.
	Option func(opt *options)

	options struct {
		env bool
	}
)

// UseEnv customizes the config to use environment variables.
func UseEnv() Option {
	return func(opt *options) {
		opt.env = true
	}
}
