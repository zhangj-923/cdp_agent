package common

import (
	"github.com/Shopify/sarama"
)

type KafkaClient struct {
	Brokers []string `toml:"brokers"`
	Topic   string   `toml:"topic"`
	client  sarama.SyncProducer
}

func (info *KafkaClient) Connect() error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	client, err := sarama.NewSyncProducer(info.Brokers, config)
	if err != nil {
		return err
	}
	info.client = client
	return nil
}

func (info *KafkaClient) Close() {
	info.client.Close()
}

func (info *KafkaClient) Send(msgValue string) (pid int32, offset int64, err error) {
	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = info.Topic
	msg.Value = sarama.StringEncoder(msgValue)

	return info.client.SendMessage(msg)
}
