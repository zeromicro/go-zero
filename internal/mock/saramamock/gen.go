package saramamock

//go:generate mockgen -destination=sarama.gen.go -package=saramamock github.com/IBM/sarama Client,PartitionConsumer,Consumer,SyncProducer

//go:generate mockgen -destination=sync_producer.mock.go -package=saramamock github.com/IBM/sarama AsyncProducer
