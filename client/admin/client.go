package admin

import (
	"log"
	"time"

	"github.com/Shopify/sarama"
)

var (
	addrs = []string{"192.168.56.25:9092"}
)

type Admin struct {
	Client sarama.ClusterAdmin
}

func NewAdminCluster() *Admin {
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
	return &Admin{Client: clusterAdmin}
}

func (a *Admin) CreateTopic(topicName string) error {
	detail := &sarama.TopicDetail{
		// topic分区数
		NumPartitions: 1,
		// topic的副本数
		ReplicationFactor: 1,
		ReplicaAssignment: nil,
		ConfigEntries:     nil,
	}
	//defer a.Client.Close()
	return a.Client.CreateTopic(topicName, detail, false)
}

func (a *Admin) CreatePartition(topic string) error {
	return a.Client.CreatePartitions(topic, 2, nil, false)
}

func (a *Admin) ListTopic() map[string]sarama.TopicDetail {
	topics, err := a.Client.ListTopics()
	if err != nil {
		log.Fatal(err)
	}
	//defer a.Client.Close()
	return topics
}

func (a *Admin) DeleteTopic(topicName string) error {
	//defer a.Client.Close()
	return a.Client.DeleteTopic(topicName)
}

func (a *Admin) GetTopicDescription(topicName string) (*sarama.TopicMetadata, error) {
	//defer a.Client.Close()
	topics, err := a.Client.DescribeTopics([]string{topicName})
	if err != nil {
		return nil, err
	}
	return topics[0], nil
}

func (a *Admin) GetTopicConfig(topicName string) ([]sarama.ConfigEntry, error) {
	//defer a.Client.Close()
	configResource := sarama.ConfigResource{
		Type: sarama.TopicResource,
		Name: topicName,
	}
	return a.Client.DescribeConfig(configResource)
}
