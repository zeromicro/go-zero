package filex

import "gopkg.in/cheggaaa/pb.v1"

type (
	Scanner interface {
		Scan() bool
		Text() string
	}

	progressScanner struct {
		Scanner
		bar *pb.ProgressBar
	}
)

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
