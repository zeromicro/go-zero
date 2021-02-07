package discovk8s

type Registry interface {
	NewSubscriber(serviceName string) Subscriber
}
