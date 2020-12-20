package visitor

type VisitResult struct {
	err error
	v   interface{}
}

func (vr *VisitResult) Result() (interface{}, error) {
	return vr.v, vr.err
}
