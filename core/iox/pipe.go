package iox

import "os"

// RedirectInOut redirects stdin to r, stdout to w, and callers need to call restore afterwards.
func RedirectInOut() (restore func(), err error) {
	var r, w *os.File
	r, w, err = os.Pipe()
	if err != nil {
		return
	}

	ow := os.Stdout
	os.Stdout = w
	or := os.Stdin
	os.Stdin = r
	restore = func() {
		os.Stdin = or
		os.Stdout = ow
	}

	return
}
