package filex

import "gopkg.in/cheggaaa/pb.v1"

type (
	// A Scanner is used to read lines.
	Scanner interface {
		// Scan checks if it has remaining to read.
		Scan() bool
		// Text returns next line.
		Text() string
	}

	progressScanner struct {
		Scanner
		bar *pb.ProgressBar
	}
)

// NewProgressScanner returns a Scanner with progress indicator.
func NewProgressScanner(scanner Scanner, bar *pb.ProgressBar) Scanner {
	return &progressScanner{
		Scanner: scanner,
		bar:     bar,
	}
}

func (ps *progressScanner) Text() string {
	s := ps.Scanner.Text()
	ps.bar.Add64(int64(len(s)) + 1) // take newlines into account
	return s
}
