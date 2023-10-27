package kafka

import (
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"gorm.io/gorm"
)

type KafkaProcessor struct {
	Database     *gorm.DB
	Producer     *ckafka.Producer
	DeliveryChan chan ckafka.Event
}

func NewKafkaProcessor(database *gorm.DB, producer *ckafka.Producer, deliveryChan chan ckafka.Event) *KafkaProcessor {
	return &KafkaProcessor{
		Database:     database,
		Producer:     producer,
		DeliveryChan: deliveryChan,
	}
}

func (kafkaProcessor *KafkaProcessor) Consume() {
	consumer, err := ckafka.NewConsumer(&ckafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
		"group.id":          "consumergroup",
		"auto.offset.reset": "earlist",
	})
	if err != nil {
		panic(err)
	}

	topics := []string{"Teste"}
	consumer.SubscribeTopics(topics, nil)

	fmt.Println("Kafka consumer has been started.")

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Println(string(msg.Value))
		}
	}
}
