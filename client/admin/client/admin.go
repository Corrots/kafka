package client

import (
	"log"
	"time"

	"github.com/Shopify/sarama"
)

var (
	addrs = []string{"192.168.56.25:9092"}
)

type AdminCluster struct {
	Client sarama.ClusterAdmin
}

func NewAdminCluster() *AdminCluster {
	config := sarama.NewConfig()
	//等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	config.Producer.Return.Successes = true
	config.Producer.Timeout = time.Second * 5

	config.Version = sarama.V2_4_0_0

	clusterAdmin, err := sarama.NewClusterAdmin(addrs, config)
	if err != nil {
		log.Fatal(err)
	}
	return &AdminCluster{Client: clusterAdmin}
}

func (c *AdminCluster) CreateTopic(topicName string) error {
	detail := &sarama.TopicDetail{
		// topic分区数
		NumPartitions: 1,
		// topic的副本数
		ReplicationFactor: 1,
		ReplicaAssignment: nil,
		ConfigEntries:     nil,
	}
	//defer c.Client.Close()
	return c.Client.CreateTopic(topicName, detail, false)
}

func (c *AdminCluster) ListTopic() map[string]sarama.TopicDetail {
	topics, err := c.Client.ListTopics()
	if err != nil {
		log.Fatal(err)
	}
	//defer c.Client.Close()
	return topics
}

func (c *AdminCluster) DeleteTopic(topicName string) error {
	//defer c.Client.Close()
	return c.Client.DeleteTopic(topicName)
}

func (c *AdminCluster) GetTopicDescription(topicName string) (*sarama.TopicMetadata, error) {
	//defer c.Client.Close()
	topics, err := c.Client.DescribeTopics([]string{topicName})
	if err != nil {
		return nil, err
	}
	return topics[0], nil
}

func (c *AdminCluster) GetTopicConfig(topicName string) ([]sarama.ConfigEntry, error) {
	//defer c.Client.Close()
	configResource := sarama.ConfigResource{
		Type: sarama.TopicResource,
		Name: topicName,
	}
	return c.Client.DescribeConfig(configResource)
}
