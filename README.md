# kafka

## 介绍
* A distributed streaming platform
* 基于zookeeper的分布式信息系统(**早期版本使用 ZooKeeper 来存储它的元数据，比如brokers信息、分区的位置和topic的配置等数据**)
* 具有高吞吐率，高性能，实时及高可靠性等特点

## 优势
1. 多Producer；
无论 kafka 多个生产者的客户端正在使用很多 topic 还是同一个 topic ， Kafka 都能够无缝处理好这些生产者。这使得 kafka 成为一个从多个前端系统聚合数据，然后提供一致的数据格式的理想系统。

2. 多Consumer；
多个消费者去读取任意的单个消息流而不相互影响。而其他的很多消息队列系统，一旦一个消息被一个客户端消费，那么这个消息就不能被其他客户端消费，这是 kafka 与其他队列不同的地方。

3. 基于硬盘的持久化消息存储；
这意味着消费者不需要实时的处理数据，也不会丢失数据。

4. 可扩展性和高性能；

## 配置文件
**config  —> server.properties**

* broker.id：broker表示当前服务器上的Kafka进程。每个节点都应该有不一样的id
* host.name：表示主机名
* log.dirs：数据存放目录，很重要。之后要看partition在哪就要去这里。
* zookeeper.connect：KAFKA是依赖ZK的。

## 总结
一个典型的kafka集群包含若干个producer（向主题发布新消息），若干consumer（从主题订阅新消息，用Consumer Offset表征消费者消费进度），cousumergroup（多个消费者实例共同组成的一个组，共同消费多个分区），若干broker（服务器端进程）。还有zookeeper。
kafka发布订阅的对象叫主题，每个Topic下可以有多个Partition，Partition中每条消息的位置信息又叫做消息位移(Offset)，Partition有副本机制，使得同一条消息能够被拷贝到多个地方以提供数据冗余，副本分为Leader Replica和 Follower Replica。

## 客户端操作
* Producers
* Consumers
* Connectors
* Stream Processors

## 客户端API类型
* AdminClient API：管理和检测Topic，broker及其他kafka对象
* Producer API：发布消息到1个或多个Topic
* Consumer API：订阅一个或多个Topic，并处理产生的消息
* Streams API：高效地将输入流转换成输出流
* Connector API：从一些源系统或应用程序中拉取数据到kafka

## Producer
* 了解Producer各项重点配置
* Producer的负载均衡等高级特性

### Producer发送模式
* 同步发送
* 异步发送
* 异步回调发送

> invalid configuration (Producer.Return.Successes must be true to be used in a SyncProducer)

