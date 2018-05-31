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
	"strconv"
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
	topics = topics + "_routing.BgpRoute_BgpRouteInitialize" + ","
	topics = topics + "_routing.BgpRoute_BgpRouteAdd" + ","
	topics = topics + "_routing.Base_RoutePurgeTimeConfig" + ","
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
	rpcKey.BrokerId = "wedge_json_broker"
	rpcKey.ClientId = "plugin_tester"
	rpcKey.TransactionId = "bgp_wedge_test"

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

	fmt.Println("invoking route initialize")

	var rit wc.Routing__BgpRoute_BgpRouteInitialize

	topic, bKey, bInput, err = rit.Marshal(rpcKey)
	if err != nil {
		log.Fatalf("%v", err)
	}

	topic = strings.Replace(topic, "/", "_", -1)

	empty := make([]byte, 0)

	message = kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   bKey,
		Value: empty,
	}

	go pubread(p, doneChan)

	p.ProduceChannel() <- &message

	// wait for delivery report goroutine to finish
	_ = <-doneChan

	time.Sleep(2 * time.Second)

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	ipnet := 117
	snet := 1
	bnet := 1
	lnet := 0
	DEST_ROUTE_TABLE := "inet.0"
	DEST_NEXT_HOP := "120.1.1.1"
	var DEST_PREFIX_LEN uint32 = 32

	/*
	 * Create bgp route add request
	 */
	for rpcs := 0; rpcs < 30; rpcs++ {
		var routeList []wc.Routing__BgpRouteEntry
		for routes := 0; routes < 1000; routes++ {
			dest_prefix_add := strconv.Itoa(ipnet) + "." + strconv.Itoa(snet) +
				"." + strconv.Itoa(bnet) + "." + strconv.Itoa(lnet)

			addr := wc.JnxBase__IpAddress{AddrString: dest_prefix_add}
			destprefix := wc.Routing__RoutePrefix{Inet: &addr}

			rtbln := wc.Routing__RouteTableName{Name: DEST_ROUTE_TABLE}
			rtbl := wc.Routing__RouteTable{RttName: &rtbln}
			nextHop := wc.JnxBase__IpAddress{AddrString: DEST_NEXT_HOP}

			routeParams :=
				wc.Routing__BgpRouteEntry{
					DestPrefix:    &destprefix,
					DestPrefixLen: DEST_PREFIX_LEN,
					Table:         &rtbl,
					PathCookie:    10,
				}

			routeParams.ProtocolNexthops =
				append(routeParams.ProtocolNexthops, nextHop)

			routeList = append(routeList, routeParams)
			if lnet >= 255 {
				snet = snet + 1
				lnet = 0
			} else if snet >= 255 {
				snet = 1
				bnet = bnet + 1
			} else {
				lnet = lnet + 1
			}
		}

		var routeAdd wc.Routing__BgpRoute_BgpRouteAdd
		routeAdd.Request.BgpRoutes = routeList
		rpcKey.RpcId = strconv.Itoa(rpcs)

		topic, bKey, bInput, err = routeAdd.Marshal(rpcKey)
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

		//fmt.Print("Press 'Enter' to send the RPC ...")
		//bufio.NewReader(os.Stdin).ReadBytes('\n')

		go pubread(p, doneChan)

		p.ProduceChannel() <- &message

		// wait for delivery report goroutine to finish
		_ = <-doneChan

		//fmt.Print("Press 'Enter' to continue...")
		//bufio.NewReader(os.Stdin).ReadBytes('\n')
		time.Sleep(1 * time.Second)
	}

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
