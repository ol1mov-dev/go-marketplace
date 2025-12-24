package kafka

import (
	"net"
	"os"
)

type KafkaConfig struct {
	Host    string
	Port    string
	Address string
}

func LoadKafkaConfigs() (*KafkaConfig, error) {
	host := os.Getenv("KAFKA_HOST")
	port := os.Getenv("KAFKA_PORT")
	address := net.JoinHostPort(host, port)

	return &KafkaConfig{
		Host:    host,
		Port:    port,
		Address: address,
	}, nil
}
