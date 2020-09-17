# Kafka如何保证消息不丢失\不重复

## 消息传递保障
Kafka提供三种消息传递保障 (sarama相关配置参数：config.Producer.RequiredAcks)
* 最多一次    —>  消息可能丢失，但绝不会重发，producer发出消息即完成发送，**无需等待**来自 broker 的确认 (NoResponse) 0
* 至少一次    —>  消息绝不会丢失，但有可能重新发送 ，当且仅当leader收到消息**返回commit确认信号**后认为发送成功 (WaitForLocal) 1
* 正好一次    —>  发送端需要等待 **ISR 列表中所有列表都确认接收数据**后才算一次发送完成，每个消息传递一次且仅一次 (WaitForAll) -1

## 丢失数据
**Broker，Producer，Consumer 都有可能丢失数据。**

### Broker
数据已经保存在broker端，但是数据却丢失了。可能原因：
	* Broker机器宕机了

### Producer
生产者丢失数据，即发送的消息未保存到broker端。可能原因：
	* 网络抖动
	* 消息本身不合格导致broker拒收 (如消息过大，超出broker的承受能力)

### Consumer
Consumer 程序有个“位移”的概念，表示的是这个 Consumer 当前消费到的 Topic 分区的位置。Kafka默认是自动提交位移的，这样可能会有个问题，假如你在pull(拉取)30条数据，处理到第20条时自动提交了offset，但是在处理21条的时候出现了异常，当你再次pull数据时，由于之前是自动提交的offset，所以是从30条之后开始拉取数据，这也就意味着21-30条的数据发生了丢失。

## 如果保证数据不“丢”失
### Broker
* 设置Replica数量，一般设置  replication.factor >= 3 ；
* 设置 min.insync.replicas ，代表消息至少被写入几个副本才算已提交。确保 replication.factor > min.insync.replicas ，一般设置为 replication.factor = min.insync.replicas + 1 。如果两者相等，有一个副本挂机，整个分区就无法正常工作了。我们不仅要考虑消息的可靠性，防止消息丢失，更应该考虑可用性问题；
* leader的选举：
    * kafka中有 Leader Replica 和 Follower Replica ，而follower replica存在的唯一目的就是防止消息丢失，并不参与具体的业务逻辑的交互。只有leader 才参与服务，follower的作用就是充当leader的候补，平时的操作也只有信息同步。
    * 要保证消息不丢失，需设置：unclean.leader.election.enable=false ，但是Kafka的可用性就会降低，具体怎么选择需要读者根据实际的业务逻辑进行权衡，可靠性优先还是可用性优先。从Kafka 0.11.0.0版本开始将此参数从true设置为false，可以看出Kafka的设计者偏向于可靠性。

### Producer
* 设置重试次数，解决网络抖动；
* 设置config.Producer.RequiredAcks = sarama.WaitForAll，即当所有的Isr写入成功后才确认committed；

### Consumer
消费端保证不丢数据，最重要就是保证offset的准确性。确保消息消费完成后再提交 offset。Consumer端设置enable.auto.commit= false，并采用手动提交offset的方式。如果处理数据时发生异常，就提交异常发生时的offset，下次消费时继续从上次失败的offset处进行消费。

## 消息去重
### Producer
Producer端可通过幂等性（Idempotence）和事务（Transaction）的机制，提供了这种精确的消息保障。

#### 幂等性
props.put(“enable.idempotence”, ture)   启用幂等性；
	* 它只能保证单分区上的幂等性，即一个幂等性 Producer 能够保证某个主题的一个分区上不出现重复消息，无法实现多个分区的幂等性。
	* 它只能实现单会话上的幂等性，不能实现跨会话的幂等性。这里的会话，你可以理解为 Producer 进程的一次运行。当你重启了 Producer 进程之后，这种幂等性保证就丧失了。

#### 事务
* props.put(“enable.idempotence”, ture)
* Producer端设置：TransactionalID
* Consumer端配置：config.Consumer.IsolationLevel = sarama.ReadCommitted （表明 Consumer 只会读取事务型 Producer 成功提交事务写入的消息）
> 事务Producer虽然在多分区的数据处理上保证了幂等，但是处理性能上相应的是会有一些下降的。

### Consumer
最好借助业务消息本身的幂等性来做。其中有些大数据组件，如 hbase，elasticsearch 天然就支持幂等操作。

## 消息乱序
这里说的乱序不是全局顺序，目前Kafka并不保证全局的消息顺序，只是提供分区级别的顺序性。如果同一分区消息乱序，很可能是重试导致的：
假设a,b两条消息，a先发送后由于发送失败重试，这时顺序就会在b的消息后面，可以设置max.in.flight.requests.per.connection=1来避免。
producer向broker发送数据的时候，其实是有多个网络连接的（因为有多个broker）。每个网络连接可以忍受producer端发送给broker消息然后消息没有响应的个数。
max.in.flight.requests.per.connection：**限制客户端在单个连接上能够发送的未响应请求的个数**。设置此值是1表示kafka broker在响应请求之前client不能再向同一个broker发送请求，但吞吐量会下降	