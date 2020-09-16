package main

import (
	"log"

	"github.com/corrots/kafka/client/producer/client/async"
)

func main() {
	producer, err := async.NewProducer()
	if err != nil {
		log.Fatal(err)
	}
	topic := "sekiro"
	game := &async.GameInfo{
		Name:       "CyberPunk 2077",
		Price:      60,
		Currency:   "USD",
		ReleaseDay: "2020-11-17",
		Platform:   "XBOX ONE",
	}
	producer.SendMessage(topic, game)

	defer producer.Async.Close()
}
