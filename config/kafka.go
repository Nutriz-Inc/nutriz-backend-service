package config

import (
	"strings"

	"github.com/IBM/sarama"
	fluxgo "github.com/MMortari/FluxGo"
)

type Kafka struct {
	ServiceName string `env:"SERVICE_NAME"`
	Brokers     string `env:"KAFKA_BROKERS"`
}

func (k Kafka) GetConfig() fluxgo.KafkaOptions {
	return fluxgo.KafkaOptions{
		Brokers: strings.Split(k.Brokers, ","),
		Auth: fluxgo.KafkaAuth{
			TlsEnabled: false,
		},
		Consumer: &fluxgo.KafkaConsumerOptions{
			GroupId:    k.ServiceName,
			AutoCommit: true,
		},
		Producer: &fluxgo.KafkaProducerOptions{
			Acks: sarama.WaitForAll,
		},
	}
}
