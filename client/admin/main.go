package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/corrots/kafka/client/admin/client"
)

func main() {
	adminClient := client.NewAdminCluster()
	topic := "corrots-topic"
	//create
	if err := adminClient.CreateTopic(topic); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("succeed, create topic: [%s]\n", topic)

	// delete
	if err := adminClient.DeleteTopic(topic); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("succeed, del topic: [%s]\n", topic)

	// list
	topics := adminClient.ListTopic()
	for k, v := range topics {
		if k != "__consumer_offsets" {
			fmt.Printf("%s: %v\n", k, v)
		}
	}

	//desc
	description, err := adminClient.GetTopicDescription(topic)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(description)

	// Topic config
	config, err := adminClient.GetTopicConfig(topic)
	if err != nil {
		log.Fatal(err)
	}
	b, err := json.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
