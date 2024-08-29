package entity

const (
	BrokerAddress      = "localhost:9092"
	ProducerPort       = ":8090"
	ConsumerReportPort = ":8091"
	ConsumerFraudPort  = ":8092"
	Topic              = "transfer"
	GroupReportID      = "report-consumer"
	GroupFraudID       = "fraud-consumer"
)
