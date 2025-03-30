package cmd

type Config struct {
	HttpPort                  string
	DbHost                    string
	DbPort                    string
	DbUser                    string
	DbPassword                string
	DbDbName                  string
	DbSslMode                 string
	GeoServiceGrpcHost        string
	KafkaHost                 string
	KafkaConsumerGroup        string
	KafkaBasketConfirmedTopic string
}
