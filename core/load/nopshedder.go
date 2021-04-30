package load

type nopShedder struct{}

func newNopShedder() Shedder {
	return nopShedder{}
}

func (s nopShedder) Allow() (Promise, error) {
	return nopPromise{}, nil
}

type nopPromise struct{}

func (p nopPromise) Pass() {
}

func (p nopPromise) Fail() {
}
