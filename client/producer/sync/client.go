package sync

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

type Producer struct {
	syncProducer sarama.SyncProducer
}

const (
	BrokerList = "192.168.56.25:9092"
	TopicName  = "sekiro"
)

func newConfig() *sarama.Config {
	config := sarama.NewConfig()
	// Wait for all in-sync replicas to ack the message
	config.Producer.RequiredAcks = sarama.WaitForAll
	// Retry up to 10 times to produce the message
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true
	// Set kafka version
	config.Version = sarama.V2_4_0_0
	return config
}

//
func NewProducer() (*Producer, error) {
	cfg := newConfig()
	brokerList := strings.Split(BrokerList, ",")
	producer, err := sarama.NewSyncProducer(brokerList, cfg)
	if err != nil {
		return nil, fmt.Errorf("new async syncProducer err: %v\n", err)
	}
	return &Producer{
		syncProducer: producer,
	}, nil
}

type UserInfo struct {
	Method   string `json:"method"`
	Username string `json:"username"`

	encoded []byte
	err     error
}

func (u *UserInfo) ensureEncoded() {
	if u.encoded == nil && u.err == nil {
		u.encoded, u.err = json.Marshal(u)
	}
}

func (u *UserInfo) Length() int {
	u.ensureEncoded()
	return len(u.encoded)
}

func (u *UserInfo) Encode() ([]byte, error) {
	u.ensureEncoded()
	return u.encoded, u.err
}

// 通过sync producer 发送message
func (p *Producer) SendMessage() error {
	key := time.Now().Format("2006-01-02 15:04:05")
	user := &UserInfo{
		Method:   "PUT",
		Username: "HaHa",
	}

	partition, offset, err := p.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicName,
		Key:   sarama.StringEncoder(key),
		Value: user,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("partition: %d, offset: %d\n", partition, offset)
	return nil
}
