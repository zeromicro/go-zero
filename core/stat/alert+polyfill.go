//go:build !linux

package stat

// Report reports given message.
func Report(string) {
}

// SetReporter sets the given reporter.
func SetReporter(func(string)) {
}
