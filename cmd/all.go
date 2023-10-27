/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jfilipedias/codepix-go/application/grpc"
	"github.com/jfilipedias/codepix-go/application/kafka"
	"github.com/jfilipedias/codepix-go/infra/db"
	"github.com/spf13/cobra"
)

var grpcPortNumber int

// allCmd represents the all command
var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Run gRPC and a Kafka Consumer",
	Run: func(cmd *cobra.Command, args []string) {
		database := db.ConnectDB()

		go grpc.StartGrpcServer(database, grpcPortNumber)

		deliveryChan := make(chan ckafka.Event)

		go kafka.DeliveryReport(deliveryChan)

		producer := kafka.NewKafkaProducer()
		kafkaProcessor := kafka.NewKafkaProcessor(database, producer, deliveryChan)
		kafkaProcessor.Consume()
	},
}

func init() {
	rootCmd.AddCommand(allCmd)
	allCmd.Flags().IntVarP(&grpcPortNumber, "grpc-port", "p", 50051, "gRPC port")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// allCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// allCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
