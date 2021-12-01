package internal

type etcdBuilder struct {
	discovBuilder
}

func (b *etcdBuilder) Scheme() string {
	return EtcdScheme
}
