package kafka

import (
	kafkago "github.com/segmentio/kafka-go"
)

func ConsumeMessages(topic string, groupId string) (*kafkago.Reader, error) {
	kafkaConf, err := LoadKafkaConfigs()

	if err != nil {
		return nil, err
	}

	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: []string{kafkaConf.Address},
		Topic:   topic,
		GroupID: groupId,
	})

	return reader, nil
}
