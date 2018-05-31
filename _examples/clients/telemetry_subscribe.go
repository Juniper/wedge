/*
 * Copyright (c) 2018, Juniper Networks, Inc.
 * All rights reserved.
 */

package main

import (
	wc "WedgeClient"
	"bufio"

	"fmt"
	"log"
	"os"
	//	"strconv"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func pubread(p *kafka.Producer, channel interface{}) {
	doneChan := channel.(chan bool)
outer:
	for e := range p.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			m := ev
			if m.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
			} else {
				fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
					*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
			}
			break outer

		default:
			fmt.Printf("Ignored event: %s\n", ev)
		}
	}

	doneChan <- true
}

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <broker> \n",
			os.Args[0])
		os.Exit(1)
	}

	broker := os.Args[1]

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Producer %v\n", p)

	doneChan := make(chan bool)

	go pubread(p, doneChan)

	/*
	 * Register the topic with the wedge
	 */
	var input, topics string
	topics = topics + "_authentication.Login_LoginCheck" + ","
	topics = topics + "_telemetry.OpenConfigTelemetry_telemetrySubscribe" + ","
	topics = topics + "_wedge.Grpc_cancelCall" + ","
	topics = topics + "_wedge.Grpc_terminateClient"

	topic := "_wedge.Topic_registerTopic"
	message := kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny},
		Value: []byte(topics),
	}

	p.ProduceChannel() <- &message

	// wait for delivery report goroutine to finish
	_ = <-doneChan

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	var rpcKey wc.RpcKey
	rpcKey.BrokerId = "wedge_telemetry_broker"
	rpcKey.ClientId = "plugin_tester"
	rpcKey.TransactionId = "wedge_telemetry_test"

	rpcKey.IpAddress = append(rpcKey.IpAddress, "2.2.2.2")
	rpcKey.IpAddress = append(rpcKey.IpAddress, "3.3.3.3")

	rpcKey.Port = "50051"
	rpcKey.RpcId = "1"

	fmt.Println("invoking LoginCheck")
	var login wc.Authentication__Login_LoginCheck
	login.Request.UserName = "username"
	login.Request.Password = "password"
	login.Request.ClientId = rpcKey.TransactionId

	topic, bKey, bInput, err := login.Marshal(rpcKey)
	if err != nil {
		log.Fatalf("%v", err)
	}

	topic = strings.Replace(topic, "/", "_", -1)
	fmt.Println("topic is ", topic)
	fmt.Println("key is ", string(bKey))
	fmt.Println("input is", string(bInput))

	message = kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   bKey,
		Value: bInput,
	}

	go pubread(p, doneChan)

	p.ProduceChannel() <- &message

	// wait for delivery report goroutine to finish
	_ = <-doneChan

	time.Sleep(2 * time.Second)

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	fmt.Println("Invoking telemetry subscribe")
	rpcKey.RpcId = "2"

        var telSub wc.Telemetry__OpenConfigTelemetry_telemetrySubscribe

        var col wc.Telemetry__Collector
        col.Address = "13.1.1.1"

        Input := new(wc.Telemetry__SubscriptionInput);
        Input.CollectorList = append(Input.CollectorList, col)

        telSub.Request.Input = Input

        var path wc.Telemetry__Path
        path.Path = "/interfaces/interface[name='fxp0']"
        path.SampleFrequency = 2000

        //telSub.Request.AdditionalConfig.NeedEos = true

        telSub.Request.PathList = append(telSub.Request.PathList, path)


	topic, bKey, bInput, err = telSub.Marshal(rpcKey)

	if err != nil {
		log.Fatalf("%v", err)
	}

	topic = strings.Replace(topic, "/", "_", -1)

	message = kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   bKey,
		Value: bInput, 
	}

	go pubread(p, doneChan)

	p.ProduceChannel() <- &message

	// wait for delivery report goroutine to finish
	_ = <-doneChan

	time.Sleep(2 * time.Second)

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	/*
	 * Cancel a call
	 */
	fmt.Println("invoking _wedge.Grpc_cancelCall")
	input = topic
	topic = "_wedge.Grpc_cancelCall"
	message = kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   bKey,
		Value: []byte(input),
	}

	go pubread(p, doneChan)

	p.ProduceChannel() <- &message

	// wait for delivery report goroutine to finish
	_ = <-doneChan

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	/*
	 * terminate a client
	 */
	fmt.Println("invoking _wedge.Grpc_terminateClient")
	input = ""
	topic = "_wedge.Grpc_terminateClient"
	message = kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   bKey,
		Value: []byte(input),
	}

	go pubread(p, doneChan)

	p.ProduceChannel() <- &message

	// wait for delivery report goroutine to finish
	_ = <-doneChan

	p.Close()

}
