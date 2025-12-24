package kafka

import (
	kafkago "github.com/segmentio/kafka-go"
)

func ProduceMessage(kafkaTopic string) (*kafkago.Writer, error) {
	kafkaConf, err := LoadKafkaConfigs()
	if err != nil {
		return nil, err
	}

	writer := &kafkago.Writer{
		Addr:     kafkago.TCP(kafkaConf.Address), // адрес вашего Kafka брокера
		Topic:    kafkaTopic,
		Balancer: &kafkago.LeastBytes{}, // балансировщик для распределения по партициям
	}

	return writer, nil
}
