package main

import (
	"fmt"
	"log"

	client "github.com/corrots/kafka/client/admin"
)

func main() {
	adminClient := client.NewAdminCluster()
	defer adminClient.Client.Close()
	topic := "sekiro"
	//
	////create
	//if err := adminClient.CreateTopic(topic); err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("succeed, create topic: [%s]\n", topic)
	//
	// delete
	//if err := adminClient.DeleteTopic(topic); err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("succeed, del topic: [%s]\n", topic)

	//// list
	//topics := adminClient.ListTopic()
	//for k, v := range topics {
	//	if k != "__consumer_offsets" {
	//		fmt.Printf("%s: %v\n", k, v)
	//	}
	//}

	// create partition
	if err := adminClient.CreatePartition(topic); err != nil {
		log.Fatal(err)
	}
	fmt.Println("create partition succeed")

	//desc
	desc, err := adminClient.GetTopicDescription(topic)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Name: %s, Partitions: %d, IsInternal: %v\n", desc.Name, desc.Partitions, desc.IsInternal)
	//
	//// Topic config
	//config, err := adminClient.GetTopicConfig(topic)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//b, err := json.Marshal(config)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(string(b))

}
