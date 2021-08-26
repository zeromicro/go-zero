package errorx

// Chain runs funs one by one until an error occurred.
func Chain(fns ...func() error) error {
	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}
