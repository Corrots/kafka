package async

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

type Producer struct {
	Async sarama.AsyncProducer
}

func NewProducer() (*Producer, error) {
	addr := []string{"192.168.56.25:9092"}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	asyncProducer, err := sarama.NewAsyncProducer(addr, config)
	if err != nil {
		return nil, err
	}
	return &Producer{Async: asyncProducer}, nil
}

func (p *Producer) SendMessage(topic string, msg *GameInfo) {
	key := time.Now().Format("2006-01-02 15:04:05")
	p.Async.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: msg,
	}
	//
	go func() {
		for err := range p.Async.Errors() {
			log.Println("Failed to send msg:", err)
		}
	}()
	//
	go func() {
		for msg := range p.Async.Successes() {
			fmt.Printf("Partition: %d, Offset: %d\n", msg.Partition, msg.Offset)
		}
	}()
}

type GameInfo struct {
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Currency   string `json:"currency"`
	ReleaseDay string `json:"release_day"`
	Platform   string `json:"platform"`

	encoded []byte
	err     error
}

func (u *GameInfo) ensureEncoded() {
	if u.encoded == nil && u.err == nil {
		u.encoded, u.err = json.Marshal(u)
	}
}

func (u *GameInfo) Length() int {
	u.ensureEncoded()
	return len(u.encoded)
}

func (u *GameInfo) Encode() ([]byte, error) {
	u.ensureEncoded()
	return u.encoded, u.err
}
