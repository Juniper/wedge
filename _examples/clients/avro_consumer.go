/*
 * Copyright (c) 2018, Juniper Networks, Inc.
 * All rights reserved.
 */

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	avro "github.com/elodina/go-avro"
)

func main() {

	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <broker> <group> <topics..>\n",
			os.Args[0])
		os.Exit(1)
	}

	broker := os.Args[1]
	group := os.Args[2]
	topics := os.Args[3:]

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               broker,
		"group.id":                        group,
		"session.timeout.ms":              6000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"}})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Consumer %v\n", c)

	err = c.SubscribeTopics(topics, nil)

	run := true

	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false

		case ev := <-c.Events():
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.Unassign()
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\nKey: %s\n\nValue: %s\n",
					e.TopicPartition, string(e.Key),
					string(e.Value))

				/*
				 * Parse the key
				 */
				keySchema, err := ioutil.ReadFile("../reference_config/wedge_key.avsc")
				kschema := avro.MustParseSchema(string(keySchema))

				reader := avro.NewGenericDatumReader()
				reader.SetSchema(kschema)

				reader.SetSchema(kschema)
				decoder := avro.NewBinaryDecoder(e.Key)

				decodedRecord := avro.NewGenericRecord(kschema)

				err = reader.Read(decodedRecord, decoder)
				if err != nil {
					fmt.Println(err)
					break
				}

				fmt.Println("Key:\n", decodedRecord)

				if strings.Contains(string(e.Key),
					"Bgp") == true {
					bgpSchema, err := ioutil.ReadFile("../reference_config/bgp_route_service.avsc")
					bschema := avro.MustParseSchema(string(bgpSchema))

					reader := avro.NewGenericDatumReader()
					reader.SetSchema(bschema)

					reader.SetSchema(bschema)
					decoder := avro.NewBinaryDecoder(e.Value)

					decodedRecord := avro.NewGenericRecord(bschema)

					err = reader.Read(decodedRecord, decoder)
					if err == nil {
						fmt.Println("Value:\n", decodedRecord)
					} else {
						fmt.Println(err)
					}

				} else if strings.Contains(string(e.Key),
					"LoginCheck") {
					aSchema, err := ioutil.ReadFile("../reference_config/authentication_service.avsc")
					aschema := avro.MustParseSchema(string(aSchema))

					reader := avro.NewGenericDatumReader()
					reader.SetSchema(aschema)

					reader.SetSchema(aschema)
					decoder := avro.NewBinaryDecoder(e.Value)

					decodedRecord := avro.NewGenericRecord(aschema)

					err = reader.Read(decodedRecord, decoder)
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Println("Value:\n", decodedRecord)
					}
				} else if strings.Contains(string(e.Key),
					"Grpc_cancelCall") == true ||
					strings.Contains(string(e.Key),
						"Grpc_terminateClient") {
					wSchema, err := ioutil.ReadFile("../reference_config/wedge_message.avsc")
					wschema := avro.MustParseSchema(string(wSchema))

					reader := avro.NewGenericDatumReader()
					reader.SetSchema(wschema)

					reader.SetSchema(wschema)
					decoder := avro.NewBinaryDecoder(e.Value)

					decodedRecord := avro.NewGenericRecord(wschema)

					err = reader.Read(decodedRecord, decoder)
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Println("Value:\n", decodedRecord)
					}
				} else {
					eSchema, err := ioutil.ReadFile("../reference_config/wedge_error.avsc")
					eschema := avro.MustParseSchema(string(eSchema))

					reader := avro.NewGenericDatumReader()
					reader.SetSchema(eschema)

					reader.SetSchema(eschema)
					decoder := avro.NewBinaryDecoder(e.Value)

					decodedRecord := avro.NewGenericRecord(eschema)

					err = reader.Read(decodedRecord, decoder)
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Println("Value:\n", decodedRecord)
					}
				}
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				run = false
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	c.Close()
}
