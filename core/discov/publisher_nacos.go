package discov

type PublisherNacos struct {
}

func NewPublisherNacos(endpoints []string, key, value string, opts ...PublisherOption) Publisher {
	publisher := &PublisherNacos{}

	for _, opt := range opts {
		opt(publisher)
	}

	return publisher
}

func (p *PublisherNacos) SetID(int64) {
}
func (p *PublisherNacos) KeepAlive() error {
	return nil
}
func (p *PublisherNacos) Pause() {

}
func (p *PublisherNacos) Resume() {

}
func (p *PublisherNacos) Stop() {

}
