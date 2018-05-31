/*
 * Copyright (c) 2018, Juniper Networks, Inc.
 * All rights reserved.
 */

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	avro "github.com/elodina/go-avro"
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

func getKeyRecord(rpcid string) []byte {
	RpcKeySchema, err := ioutil.ReadFile("../reference_config/avsc/rpc_key.avsc")
	if err != nil {
		fmt.Println("Error", err, "parsing rpc_key.avsc")
		os.Exit(0)
	}

	schema := avro.MustParseSchema(string(RpcKeySchema))
	record := avro.NewGenericRecord(schema)

	var client_id []string
	client_id = append(client_id, "confluent_wedge_test")

	record.Set("BrokerId", "wedge_avro_broker")
	record.Set("ClientId", "plugin_tester")
	record.Set("TransactionId", "confluent_wedge_test")
	record.Set("RpcId", rpcid)
	record.Set("Port", "50051")

	var ip_address []string
	ip_address = append(ip_address, "2.2.2.2")
	record.Set("IpAddress", ip_address)

	keyWriter := avro.NewGenericDatumWriter()

	keyWriter.SetSchema(schema)

	// Create a new Buffer and Encoder to write to this Buffer
	keybuf := new(bytes.Buffer)
	encoder := avro.NewBinaryEncoder(keybuf)

	// Write the record
	err = keyWriter.Write(record, encoder)
	if err != nil {
		panic(err)
	}

	return []byte(string(keybuf.Bytes()))
}

func getWedgeMsgRecord(message string) []byte {
	WedgeMsgSchema, err := ioutil.ReadFile("../reference_config/avsc/wedge_message.avsc")
	if err != nil {
		fmt.Println("Error parsing wedge_message.avsc")
		os.Exit(0)
	}

	schema := avro.MustParseSchema(string(WedgeMsgSchema))
	record := avro.NewGenericRecord(schema)

	record.Set("Message", message)

	msgWriter := avro.NewGenericDatumWriter()

	msgWriter.SetSchema(schema)

	// Create a new Buffer and Encoder to write to this Buffer
	msgbuf := new(bytes.Buffer)
	encoder := avro.NewBinaryEncoder(msgbuf)

	// Write the record
	err = msgWriter.Write(record, encoder)
	if err != nil {
		panic(err)
	}

	return []byte(string(msgbuf.Bytes()))
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

	/*
	 * Register the topic with the wedge
	 */
	var topics string
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

	go pubread(p, doneChan)

	p.ProduceChannel() <- &message

	// wait for delivery report goroutine to finish
	_ = <-doneChan

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	fmt.Println("Invoking Login check RPC")

	/*
	 * Login check schema parsing
	 */

	loginSchema, err := ioutil.ReadFile("../reference_config/avsc/authentication_service.avsc")
	fmt.Println("schema is", string(loginSchema))
	ls := avro.MustParseSchema(string(loginSchema))
	loginRecord := avro.NewGenericRecord(ls)

	loginRecord.Set("user_name", "username")
	loginRecord.Set("password", "password")
	loginRecord.Set("client_id", "avsc_wedge_test")

	// Create a new Buffer and Encoder to write to this Buffer
	writer := avro.NewGenericDatumWriter()
	writer.SetSchema(ls)

	valuebuf := new(bytes.Buffer)
	valencoder := avro.NewBinaryEncoder(valuebuf)

	// Write the record
	err = writer.Write(loginRecord, valencoder)
	if err != nil {
		panic(err)
	}

	topic = "_authentication.Login_LoginCheck"
	message = kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   getKeyRecord("1"),
		Value: []byte(string(valuebuf.Bytes())),
	}

	go pubread(p, doneChan)

	p.ProduceChannel() <- &message

	// wait for delivery report goroutine to finish
	_ = <-doneChan

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	fmt.Println("Invoking BgpRouteInitialize RPC")

	bgpSchema, err := ioutil.ReadFile("../reference_config/avsc/bgp_route_service.avsc")
	bschema := avro.MustParseSchema(string(bgpSchema))
	/*
	 * Perform BGP route initialize
	 */
	routeInitRecord := avro.NewGenericRecord(bschema)
	routeInitRecord.Set("WedgePlaceholder", nil)

	initWriter := avro.NewGenericDatumWriter()
	initWriter.SetSchema(bschema)

	valuebuf = new(bytes.Buffer)
	valencoder = avro.NewBinaryEncoder(valuebuf)

	err = initWriter.Write(routeInitRecord, valencoder)
	if err != nil {
		panic(err)
	}

	topic = "_routing.BgpRoute_BgpRouteInitialize"
	message = kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   getKeyRecord("2"),
		Value: []byte(string(valuebuf.Bytes())),
	}

	go pubread(p, doneChan)

	p.ProduceChannel() <- &message

	// wait for delivery report goroutine to finish
	_ = <-doneChan

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	fmt.Println("Invoking BgpRouteAdd RPC")

	/*
	 * Perform BGP route add
	 */
	bgpRecord := avro.NewGenericRecord(bschema)
	ipnet := 117
	snet := 1
	bnet := 1
	lnet := 0
	DEST_ROUTE_TABLE := "inet.0"
	DEST_NEXT_HOP := "120.1.1.1"
	var DEST_PREFIX_LEN int32 = 32

	for i := 0; i < 30; i++ {
		routeEntries := make([]*avro.GenericRecord, 1000)
		for routes := 0; routes < 1000; routes++ {
			routeEntries[routes] = avro.NewGenericRecord(bschema)
			dest_prefix_add := strconv.Itoa(ipnet) + "." + strconv.Itoa(snet) +
				"." + strconv.Itoa(bnet) + "." + strconv.Itoa(lnet)

			routePrefix := avro.NewGenericRecord(bschema)
			ipAddress := avro.NewGenericRecord(bschema)

			ipAddress.Set("addr_string", dest_prefix_add)
			routePrefix.Set("inet", ipAddress)

			record_Routing__RouteTableName := avro.NewGenericRecord(bschema)
			record_Routing__RouteTableName.Set("name", DEST_ROUTE_TABLE)

			record_Routing__RouteTable := avro.NewGenericRecord(bschema)
			record_Routing__RouteTable.Set("rtt_name", record_Routing__RouteTableName)

			subRecords := make([]*avro.GenericRecord, 1)
			pAddress := avro.NewGenericRecord(bschema)
			pAddress.Set("addr_string", DEST_NEXT_HOP)
			subRecords[0] = pAddress

			routeEntries[routes].Set("dest_prefix", routePrefix)
			routeEntries[routes].Set("dest_prefix_len", DEST_PREFIX_LEN)
			routeEntries[routes].Set("table", record_Routing__RouteTable)
			routeEntries[routes].Set("path_cookie", int64(10))
			routeEntries[routes].Set("protocol_nexthops", subRecords)

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

		bgpRecord.Set("bgp_routes", routeEntries)

		// Create a new Buffer and Encoder to write to this Buffer
		bgpWriter := avro.NewGenericDatumWriter()
		bgpWriter.SetSchema(bschema)

		valuebuf = new(bytes.Buffer)
		valencoder = avro.NewBinaryEncoder(valuebuf)

		// Write the record
		err = bgpWriter.Write(bgpRecord, valencoder)
		if err != nil {
			panic(err)
		}

		topic = "_routing.BgpRoute_BgpRouteAdd"
		message = kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Key:   getKeyRecord("3"),
			Value: []byte(string(valuebuf.Bytes())),
		}

		go pubread(p, doneChan)

		p.ProduceChannel() <- &message

		// wait for delivery report goroutine to finish
		_ = <-doneChan

		time.Sleep(1 * time.Second)
	}

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	/*
	 * Cancel a call
	 */
	fmt.Println("invoking _wedge.Grpc_cancelCall")
	input := topic
	topic = "_wedge.Grpc_cancelCall"
	message = kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   getKeyRecord("4"),
		Value: getWedgeMsgRecord(input),
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
		Key:   getKeyRecord("5"),
		Value: getWedgeMsgRecord(input),
	}

	go pubread(p, doneChan)

	p.ProduceChannel() <- &message

	// wait for delivery report goroutine to finish
	_ = <-doneChan

	p.Close()
}
