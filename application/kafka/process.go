package kafka

import (
	"fmt"
	"os"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jfilipedias/codepix-go/application/factory"
	appmodel "github.com/jfilipedias/codepix-go/application/model"
	"github.com/jfilipedias/codepix-go/application/usecase"
	"github.com/jfilipedias/codepix-go/domain/model"
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
		"bootstrap.servers": os.Getenv("kafkaBootstrapServers"),
		"group.id":          os.Getenv("kafkaConsumerGroupId"),
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}

	topics := []string{os.Getenv("kafkaTransactionTopic"), os.Getenv("kafkaTransactionConfirmationTopic")}
	consumer.SubscribeTopics(topics, nil)

	fmt.Println("Kafka consumer has been started.")

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			kafkaProcessor.processMessage(msg)
		}
	}
}

func (kafkaProcessor *KafkaProcessor) processMessage(msg *ckafka.Message) {
	transactionsTopic := "transactions"
	transactionConfirmationTopic := "transaction_confirmation"

	switch topic := *msg.TopicPartition.Topic; topic {
	case transactionsTopic:
		kafkaProcessor.processTransaction(msg)
	case transactionConfirmationTopic:
		kafkaProcessor.processTransactionConfirmation(msg)
	default:
		fmt.Println("Not a valid topic.", string(msg.Value))
	}
}

func (kafkaProcessor *KafkaProcessor) processTransaction(msg *ckafka.Message) error {
	transaction := appmodel.NewTransaction()
	err := transaction.ParseJson(msg.Value)
	if err != nil {
		return err
	}

	transactionUseCase := factory.TransactionUseCaseFactory(kafkaProcessor.Database)

	createdTransaction, err := transactionUseCase.Register(
		transaction.AccountId,
		transaction.Amount,
		transaction.PixKeyTo,
		transaction.PixKeyToKind,
		transaction.Description,
	)
	if err != nil {
		fmt.Errorf("Error registering transaction.", err)
		return err
	}

	transaction.ID = createdTransaction.ID
	transaction.Status = model.TransactionPending

	transactionJson, err := transaction.ToJson()
	if err != nil {
		return err
	}

	topic := "bank" + createdTransaction.PixKeyTo.Account.Bank.Code
	err = Publish(string(transactionJson), topic, kafkaProcessor.Producer, kafkaProcessor.DeliveryChan)
	if err != nil {
		return err
	}

	return nil
}

func (kafkaProcessor *KafkaProcessor) confirmTransaction(transaction *appmodel.Transaction, transactionUseCase usecase.TransactionUseCase) error {
	confirmedTransaction, err := transactionUseCase.Confirm(transaction.ID)
	if err != nil {
		return err
	}

	topic := "banck" + confirmedTransaction.AccountFrom.Bank.Code
	transactionJson, err := transaction.ToJson()
	if err != nil {
		return err
	}

	err = Publish(string(transactionJson), topic, kafkaProcessor.Producer, kafkaProcessor.DeliveryChan)
	if err != nil {
		return err
	}

	return nil
}

func (kafkaProcessor *KafkaProcessor) processTransactionConfirmation(msg *ckafka.Message) error {
	transaction := appmodel.NewTransaction()
	err := transaction.ParseJson(msg.Value)
	if err != nil {
		return err
	}

	transactionUseCase := factory.TransactionUseCaseFactory(kafkaProcessor.Database)

	if transaction.Status == model.TransactionConfirmed {
		err = kafkaProcessor.confirmTransaction(transaction, transactionUseCase)
		if err != nil {
			return err
		}
	} else if transaction.Status == model.TransactionCompleted {
		_, err = transactionUseCase.Complete(transaction.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
