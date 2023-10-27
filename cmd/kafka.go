/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jfilipedias/codepix-go/application/kafka"
	"github.com/spf13/cobra"
)

// kafkaCmd represents the kafka command
var kafkaCmd = &cobra.Command{
	Use:   "kafka",
	Short: "Start combining transaction using Apache Kafka.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Produzindo mensagem.")
		producer := kafka.NewKafkaProducer()
		deliveryChan := make(chan ckafka.Event)
		kafka.Publish("Ola Kafka", "Teste", producer, deliveryChan)
		kafka.DeliveryReport(deliveryChan)
	},
}

func init() {
	rootCmd.AddCommand(kafkaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kafkaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kafkaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
