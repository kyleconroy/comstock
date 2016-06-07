package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

type KafkaConfig struct {
	Url           string // KAFKA_URL
	TrustedCert   string // KAFKA_TRUSTED_CERT
	ClientCertKey string // KAFKA_CLIENT_CERT_KEY
	ClientCert    string // KAFKA_CLIENT_CERT
	Topic         string // COMSTOCK_KAFKA_TOPIC
}

type Message struct {
	Partition int32           `json:"partition"`
	Offset    int64           `json:"offset"`
	Value     string          `json:"value"`
	Metadata  MessageMetadata `json:"metadata"`
}

type MessageMetadata struct {
	ReceivedAt time.Time `json:"received_at"`
}

const (
	ConsumerId    = "comstock"
	ConsumerGroup = "comstock-group"
)

func newProducer() (sarama.AsyncProducer, error) {
	config, err := parseConfig()
	if err != nil {
		return nil, err
	}
	tlsConfig, err := config.createTlsConfig()
	if err != nil {
		return nil, err
	}

	brokerAddrs, err := config.brokerAddresses()
	if err != nil {
		return nil, err
	}

	return config.createKafkaProducer(brokerAddrs, tlsConfig)
}

func parseConfig() (KafkaConfig, error) {
	keys := []string{
		"KAFKA_URL",
		"KAFKA_TRUSTED_CERT",
		"KAFKA_CLIENT_CERT_KEY",
		"KAFKA_CLIENT_CERT",
		"COMSTOCK_KAFKA_TOPIC",
	}

	for _, k := range keys {
		if os.Getenv(k) == "" {
			return KafkaConfig{}, fmt.Errorf("Missing %s environment variable", k)
		}
	}

	c := KafkaConfig{
		Url:           os.Getenv("KAFKA_URL"),
		TrustedCert:   os.Getenv("KAFKA_TRUSTED_CERT"),
		ClientCertKey: os.Getenv("KAFKA_CLIENT_CERT_KEY"),
		ClientCert:    os.Getenv("KAFKA_CLIENT_CERT"),
		Topic:         os.Getenv("COMSTOCK_KAFKA_TOPIC"),
	}
	return c, nil
}

// It receives messages as http bodies on /messages,
// and posts them directly to a Kafka topic.
// func (kc *KafkaClient) messagesPOST(c *gin.Context) {
// 	message, err := ioutil.ReadAll(c.Request.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	msg := &sarama.ProducerMessage{
// 		Topic: KafkaTopic,
// 		Key:   sarama.ByteEncoder(c.Request.RemoteAddr),
// 		Value: sarama.ByteEncoder(message),
// 	}
//
// 	kc.producer.Input() <- msg
// }

// Create the TLS context, using the key and certificates provided.
func (ac *KafkaConfig) createTlsConfig() (*tls.Config, error) {
	cert, err := tls.X509KeyPair([]byte(ac.ClientCert), []byte(ac.ClientCertKey))
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM([]byte(ac.TrustedCert))

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		RootCAs:            certPool,
	}
	tlsConfig.BuildNameToCertificate()
	return tlsConfig, nil
}

// Connect a consumer. Consumers in Kafka have a "group" id, which
// denotes how consumers balance work. Each group coordinates
// which partitions to process between its nodes.
// For the demo app, there's only one group, but a production app
// could use separate groups for e.g. processing events and archiving
// raw events to S3 for longer term storage
func (ac *KafkaConfig) createKafkaConsumer(brokers []string, tc *tls.Config) (*cluster.Consumer, error) {
	config := cluster.NewConfig()

	config.Net.TLS.Config = tc
	config.Net.TLS.Enable = true
	config.Group.PartitionStrategy = cluster.StrategyRoundRobin
	config.ClientID = ConsumerId
	config.Consumer.Return.Errors = true

	err := config.Validate()
	if err != nil {
		return nil, err
	}

	consumer, err := cluster.NewConsumer(brokers, ConsumerGroup, []string{ac.Topic}, config)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

// Create the Kafka asynchronous producer
func (ac *KafkaConfig) createKafkaProducer(brokers []string, tc *tls.Config) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()

	config.Net.TLS.Config = tc
	config.Net.TLS.Enable = true
	config.Producer.Return.Errors = true
	config.ClientID = ConsumerId

	err := config.Validate()
	if err != nil {
		return nil, err
	}
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return producer, nil
}

// Extract the host:port pairs from the Kafka URL(s)
func (ac *KafkaConfig) brokerAddresses() ([]string, error) {
	urls := strings.Split(ac.Url, ",")
	addrs := make([]string, len(urls))
	for i, v := range urls {
		u, err := url.Parse(v)
		if err != nil {
			return []string{}, err
		}
		addrs[i] = u.Host
	}
	return addrs, nil
}
