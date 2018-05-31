/*
 * Copyright (c) 2018, Juniper Networks, Inc.
 * All rights reserved.
 */

package implementations

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	wk "git.juniper.net/sksubra/wedge/plugins/kafka"
	wu "git.juniper.net/sksubra/wedge/util"
	kafka "github.com/confluentinc/confluent-kafka-go/kafka"
	avro "github.com/elodina/go-avro"
	"github.com/srikanth2212/jsonez"
)

//var KafkaClientID = wk.DEFAULT_KAFKA_CLIENT_ID
var InFluxTopic = "telegraf"
var FluentdTopic = "fluentd.wedge"

type jsonFunc func(val interface{}, paths ...string) error

/*
 * Callback map for pre write functions
 */
var PreWriteFuncMap map[string]wu.PreWriteFunc

func PreWriteMapInit() {
	PreWriteFuncMap = make(map[string]wu.PreWriteFunc)
	PreWriteFuncMap["kafka-json-pre-write"] = JSONPreWriteCB
	PreWriteFuncMap["kafka-avro-pre-write"] = AvroPreWriteCB
	PreWriteFuncMap["kafka-influx-pre-write"] = InfluxPreWriteCB
	PreWriteFuncMap["kafka-fluentd-pre-write"] = FluentdPreWriteCB
	PreWriteFuncMap["kafka-telemetry-pre-write"] = TelemetryPreWriteCB
}

/*
 * Callback map for post read functions
 */
var PostReadFuncMap map[string]wu.PostReadFunc

func PostReadMapInit() {
	PostReadFuncMap = make(map[string]wu.PostReadFunc)
	PostReadFuncMap["kafka-json-post-read"] = JSONPostReadCB
	PostReadFuncMap["kafka-avro-post-read"] = AvroPostReadCB
}

/*
 * Callback map for error functions
 */
var ErrorFuncMap map[string]wu.ErrorFunc

func ErrorFuncMapInit() {
	ErrorFuncMap = make(map[string]wu.ErrorFunc)
	ErrorFuncMap["avro"] = AvroErrorCB
}

/*
 * Avro payload specific callback functions
 */
var RpcKeySchema = `[
    {
    	"type": "record",
    	"name": "RpcMetadata",
    	"fields": [
    		{
    			"name": "Key",
 				"type": "string"
 			},
    		{
    			"name": "Value",
 				"type": "string"
 			}
    	]
    },
 	{
 		"type": "record",
 		"name": "RpcKey",
 		"fields": [
 		 	{
 				"name": "Rpc",
 				"type": ["null", "string"],
                "default": null
 			},
 			{
 				"name": "BrokerId",
 				"type": ["null", "string"],
                "default": null
 			},
 			{
 				"name": "ClientId",
 				"type": ["null", "string"],
                "default": null
 			},
 			{
 				"name": "TransactionId",
 				"type": ["null", "string"],
                "default": null
 			},
 			{
 				"name": "RpcId",
 				"type": ["null", "string"],
                "default": null
 			},
 			{
 				"name": "IpAddress",
 				"type": ["null", {"type": "array", "items": "string"}],
 				"default": null
 			},
 			{
 				"name": "Port",
 				"type": ["null", "string"],
                "default": null
 			},
 			{
 				"name": "Metadata",
 				"type": ["null", {"type": "array", "items": "RpcMetadata"}],
 				"default": null
 			
 			}
 		]
	}
 ]`

var ErrorMsgSchema = `[
    {
    	"type": "record",
    	"name": "Error",
    	"fields": [
    		{
 				"name": "Message",
 				"type": ["null", "string"],
                "default": null
 			}
    	]
    }
]`

var WedgeMsgSchema = `[
    {
    	"type": "record",
    	"name": "Payload",
    	"fields": [
    		{
 				"name": "Message",
 				"type": ["null", "string"],
                "default": null
 			}
    	]
    }
]`

/*
 * Function to build JSON payload for a specific
 * AVRO type. For fields that belong to a oneof, make
 * sure that only one of thr options is provided
 */
func addJSONVal(record *avro.GenericRecord, cur *jsonez.GoJSON,
	rName, key string, fDesc *wu.AvroFieldDesc, val interface{},
	repeated bool, curOneofMap map[string]string) error {

	var child *jsonez.GoJSON
	var jFunc jsonFunc
	var ok bool
	var prevVal string

	if repeated == true {
		jFunc = cur.AddToArray
	} else {
		jFunc = cur.AddVal
	}

	/*
	 * If the field belongs to a one ofm then check if
	 * another option with thin the same field is not
	 */
	if fDesc.ParentOneof != "" && fDesc.ParentOneof != "NULL" {
		/*
		 * Check if the currentOneof Map processed so far
		 * has the oneof entry and if so report an error
		 */
		if prevVal, ok = curOneofMap[fDesc.ParentOneof]; ok {
			errorStr := fmt.Sprintf("%s: Oneof field %s has two options "+
				"%s and %s specified", wu.FuncName(), key,
				fDesc.ParentOneof, prevVal)
			return errors.New(errorStr)
		} else {
			curOneofMap[fDesc.ParentOneof] = key
		}
	}

	switch fDesc.Ftype {
	case wu.AVRO_TYPE_BOOLEAN:
		jFunc(val.(bool), key)
	case wu.AVRO_TYPE_INT:
		jFunc(int(val.(int32)), key)
	case wu.AVRO_TYPE_LONG:
		jFunc(int(val.(int64)), key)
	case wu.AVRO_TYPE_FLOAT:
		jFunc(float64(val.(float32)), key)
	case wu.AVRO_TYPE_DOUBLE:
		jFunc(val.(float64), key)
	case wu.AVRO_TYPE_BYTES:
		v := val.(bytes.Buffer)
		jFunc(v.String(), key)
	case wu.AVRO_TYPE_STRING:
		jFunc(val.(string), key)
	case wu.AVRO_TYPE_ENUM:
		/*
		 * Get the integer value corresponding to enum value
		 * and add it
		 */
		eStr := val.(string)
		eval, ok := fDesc.EnumMap[eStr]
		if !ok {
			errorStr := fmt.Sprintf("%s: Enum value %s not found for "+
				"Field %s in record %s", wu.FuncName(), eStr, key, rName)
			return errors.New(errorStr)
		}
		jFunc(eval, key)
	case wu.AVRO_TYPE_RECORD:
		/*
		 * Create an object entry and tag it as the child of the
		 * current root. Then make a recursive call to process
		 * the contents
		 */
		child = jsonez.AllocObject()
		child.Key = key
		if cur.Jsontype == jsonez.JSON_OBJECT {
			cur.AddEntryToObject(key, child)
		} else if cur.Jsontype == jsonez.JSON_ARRAY {
			cur.AddEntryToArray(child)
		} else {
			errorStr := fmt.Sprintf("%s: Unexpected type JSON object with "+
				"key %s for record %s", wu.FuncName(), key, rName)
			return errors.New(errorStr)
		}

		return avroWalkRecordMap(val.(*avro.GenericRecord), child, key,
			fDesc.SubRecMap)

	default:
		errorStr := fmt.Sprintf("%s: Unknown type found for "+
			"Field %s in record %s", wu.FuncName(), key, rName)
		return errors.New(errorStr)
	}
	return nil
}

/*
 * Function to walk a record map and build JSON
 * object
 */
func avroWalkRecordMap(record *avro.GenericRecord, cur *jsonez.GoJSON,
	rName string, fDescMap wu.AFdMap) error {
	var child *jsonez.GoJSON
	var err error
	var val interface{}
	var key string
	var curOneofMap = make(map[string]string)

	for key, val = range record.Map() {
		fDesc, ok := fDescMap[key]
		if !ok {
			errorStr := fmt.Sprintf("%s: Field %s not found in record %s",
				wu.FuncName(), key, rName)
			return errors.New(errorStr)
		}

		/*
		 * If the value is null, then check if a placeholder
		 * message field for proto. If not, then this is the
		 * default value for the field which needs to be ignored
		 */
		if val == nil {
			if fDesc.Ftype == wu.AVRO_TYPE_NULL {
				if strings.Compare(key, wu.FIELD_PLACEHOLDER) == 0 {

					/*
						child = jsonez.AllocObject()
						child.Jsontype = jsonez.JSON_NULL
						cur.AddEntryToObject(key, child)
					*/
				} else {
					log.Println("Ignoring field", key)
				}
			}
		} else if fDesc.Repeated == true {
			/*
			 * The field of type array. For a record type,
			 * create an array object and add it to the current
			 * node and pass that as the parent to the add
			 * function
			 */
			if fDesc.Ftype == wu.AVRO_TYPE_RECORD {
				child = jsonez.AllocArray()
				cur.Key = key

				if cur.Jsontype == jsonez.JSON_OBJECT {
					cur.AddEntryToObject(key, child)
				} else if cur.Jsontype == jsonez.JSON_ARRAY {
					cur.AddEntryToArray(child)
				}
			} else {
				child = cur
			}

			/*
			 * To fetch the *GenericRecord corresponding to a record object,
			 * use the Get() method
			 */
			decodedArray := record.Get(key).([]interface{})

			for _, entry := range decodedArray {
				if err = addJSONVal(record, child, rName, key, fDesc,
					entry, true, curOneofMap); err != nil {
					return err
				}
			}
		} else {

			val = record.Get(key).(interface{})
			if err = addJSONVal(record, cur, rName, key, fDesc,
				val, false, curOneofMap); err != nil {
				return err
			}
		}
	}

	return nil
}

/*
 * Function to parse the AVRO input and build
 * a JSON payload for further processing
 */
func avroBuildJsonOutput(avroInput []byte, rpc string) (string, error) {
	var pDesc *wu.AvroProtocolDesc
	var root *jsonez.GoJSON
	var err error
	var jsonPayload string

	/*
	 * Fetch the protocol descriptor and get the schema
	 */
	if pDesc, err = wu.GetProtocolDesc(rpc); err != nil {
		return "", err
	}

	schema := avro.MustParseSchema(wu.GetAvroSchema(pDesc))

	reader := avro.NewGenericDatumReader()

	// SetSchema must be called before calling Read
	reader.SetSchema(schema)

	// Create a new Decoder with a given buffer
	decoder := avro.NewBinaryDecoder(avroInput)

	decodedRecord := avro.NewGenericRecord(schema)

	// Read data into given GenericRecord with a given Decoder
	err = reader.Read(decodedRecord, decoder)
	if err != nil {
		errorStr := fmt.Sprintf("%s: Reading Avro Input failed for RPC %s "+
			"with error %v", wu.FuncName(), rpc, err)
		return "", errors.New(errorStr)
	}

	/*
	 * Walk the record map and build the json input
	 */
	root = jsonez.AllocObject()

	if err = avroWalkRecordMap(decodedRecord, root, pDesc.Request,
		pDesc.RequestDesc.FieldDescMap); err != nil {
		return "", err
	}

	jsonPayload = string(jsonez.GoJSONPrint(root))

	return jsonPayload, nil
}

/*
 * Call back function to be invoked post reading the data
 */
func AvroPostReadCB(in interface{}) ([]wu.MsgFormat, error) {
	var err error
	var input *kafka.Message
	var errorStr string
	var msg wu.MsgFormat
	var ret []wu.MsgFormat
	var bid, cid, tid, rid, port string
	var metadata map[string]string
	var jsonPayload string

	input = in.(*kafka.Message)

	/*
	 * Build the JSON Key structure from Avro
	 */
	KeySchema := avro.MustParseSchema(RpcKeySchema)
	reader := avro.NewGenericDatumReader()
	reader.SetSchema(KeySchema)
	decoder := avro.NewBinaryDecoder(input.Key)
	decodedRecord := avro.NewGenericRecord(KeySchema)
	err = reader.Read(decodedRecord, decoder)

	if err != nil {
		log.Printf("%s: Error %v while parsing key", wu.FuncName(),
			err)
		errorStr = fmt.Sprintf("Error parsing key")
		return nil, errors.New(errorStr)
	}

	/*
	 * Walk the record map and populate the Key fields
	 */
	decodedMap := decodedRecord.Map()

	val, ok := decodedMap["BrokerId"]
	if !ok {
		log.Printf("%s: Error %v while fetching BrokerId", wu.FuncName(),
			err)
		errorStr = fmt.Sprintf("Mandatory parameter BrokerId not found")
		return nil, errors.New(errorStr)
	}
	bid = val.(string)
	

	val, ok = decodedMap["ClientId"]
	if !ok {
		log.Printf("%s: Error %v while fetching ClientId", wu.FuncName(),
			err)
		errorStr = fmt.Sprintf("Mandatory parameter ClientId not found")
		return nil, errors.New(errorStr)
	}
	cid = val.(string)

	val, ok = decodedMap["TransactionId"]
	if !ok {
		log.Printf("%s: Error %v while fetching TransactionId", wu.FuncName(),
			err)
		errorStr = fmt.Sprintf("Mandatory parameter TransactionId not found")
		return nil, errors.New(errorStr)
	}
	tid = val.(string)

	val, ok = decodedMap["RpcId"]
	if !ok {
		log.Printf("%s: Error %v while fetching RpcId", wu.FuncName(),
			err)
		errorStr = fmt.Sprintf("Mandatory parameter RpcId not found")
		return nil, errors.New(errorStr)
	}
	rid = val.(string)

	val, ok = decodedMap["Port"]
	if !ok {
		log.Printf("%s: Error %v while fetching Port", wu.FuncName(),
			err)
		errorStr = fmt.Sprintf("Mandatory parameter Port not found")
		return nil, errors.New(errorStr)
	}
	port = val.(string)

	val, ok = decodedMap["Metadata"]
	if ok {
		metadata = make(map[string]string)
		if val != nil {
			for _, m := range val.([]interface{}) {
				meta := m.(*avro.GenericRecord)

				for mk, mv := range meta.Map() {
					metadata[mk] = mv.(string)
				}
			}
		}
	}

	rpc := strings.Replace(*input.TopicPartition.Topic, "_", "/", -1)

	/*
	 * Get the JSON payload from avro input
	 */
	if strings.Contains(rpc, "wedge") == false {
		if jsonPayload, err = avroBuildJsonOutput(input.Value, rpc); err != nil {
			return nil, err
		}
	} else {
		/*
		 * For broker related messages the payload will be a string
		 * Use WedgeMsgSchema for these purposes
		 */
		MsgSchema := avro.MustParseSchema(WedgeMsgSchema)
		reader := avro.NewGenericDatumReader()
		reader.SetSchema(MsgSchema)
		decoder := avro.NewBinaryDecoder(input.Value)
		decodedRecord := avro.NewGenericRecord(MsgSchema)
		err = reader.Read(decodedRecord, decoder)

		if err != nil {
			log.Printf("%s: Error %v while parsing Value", wu.FuncName(),
				err)
			errorStr = fmt.Sprintf("Error parsing Value")
			return nil, errors.New(errorStr)
		}

		decodedMap := decodedRecord.Map()

		val, ok = decodedMap["Message"]
		if !ok {
			log.Printf("%s: Error while fetching Message", wu.FuncName())
			errorStr = fmt.Sprintf("Mandatory parameter Message not found")
			return nil, errors.New(errorStr)
		}

		if val != nil {
			jsonPayload = val.(string)
		} else {
			jsonPayload = ""
		}
	}

	/*
	 * Get the schema corresponding to the rpc and build the
	 * equivalent JSON value
	 */

	val, ok = decodedMap["IpAddress"]
	if !ok {
		log.Printf("%s: Error %v while fetching IpAddress", wu.FuncName(),
			err)
		errorStr = fmt.Sprintf("Mandatory parameter IpAddress not found")
		return nil, errors.New(errorStr)
	}

	/*
	 * Walk the ip address map and append the message
	 * structure to the message list
	 */
	for _, ip := range val.([]interface{}) {
		msg.IpAddress = ip.(string)
		msg.BrokerId = bid
		msg.ClientId = cid
		msg.TransactionId = tid
		msg.RpcId = rid
		msg.Port = port
		msg.Rpc = rpc
		msg.Value = jsonPayload
		msg.Metadata = metadata
		ret = append(ret, msg)
	}

	return ret, nil

}

/*
 * Function to build JSON payload for a specific
 * AVRO type
 */
func addAvroVal(record *avro.GenericRecord, cur *jsonez.GoJSON,
	rName string, fDesc *wu.AvroFieldDesc, schema *avro.Schema) error {

	switch fDesc.Ftype {
	case wu.AVRO_TYPE_BOOLEAN:
		record.Set(cur.Key, cur.Valbool)
	case wu.AVRO_TYPE_INT:
		var val int32
		if cur.Jsontype == jsonez.JSON_UINT {
			val = int32(cur.Valuint)
		} else {
			val = int32(cur.Valint)
		}
		record.Set(cur.Key, val)
	case wu.AVRO_TYPE_LONG:
		var val int64
		if cur.Jsontype == jsonez.JSON_UINT {
			val = int64(cur.Valuint)
		} else {
			val = int64(cur.Valint)
		}
		record.Set(cur.Key, val)
	case wu.AVRO_TYPE_FLOAT:
		record.Set(cur.Key, float32(cur.Valdouble))
	case wu.AVRO_TYPE_DOUBLE:
		record.Set(cur.Key, cur.Valdouble)
	case wu.AVRO_TYPE_BYTES:
		record.Set(cur.Key, []byte(cur.Valstr))
	case wu.AVRO_TYPE_STRING:
		record.Set(cur.Key, cur.Valstr)
	case wu.AVRO_TYPE_ENUM:
		/*
		 * Walk the enum map of this field to fetch
		 * the key
		 */
		var val int32
		var found bool
		if cur.Jsontype == jsonez.JSON_UINT {
			val = int32(cur.Valuint)
		} else {
			val = int32(cur.Valint)
		}

		for k, v := range fDesc.EnumMap {
			if int32(v) == val {
				record.Set(cur.Key, k)
				found = true
				break
			}

			if !found {
				errorStr := fmt.Sprintf("Enum corresponding to value %d was "+
					"not found in field %s of record %s", val, cur.Key, rName)
				return errors.New(errorStr)
			}
		}
	case wu.AVRO_TYPE_RECORD:
		subRecord := avro.NewGenericRecord(*schema)
		record.Set(cur.Key, subRecord)

		return avroWalkJSON(subRecord, cur, fDesc.SubRecord, fDesc.SubRecMap,
			schema)

	default:
		record.Set(cur.Key, nil)
	}

	return nil
}

/*
 * Walk the GoJSON structs and construct AVRO payload
 */

func avroWalkJSON(record *avro.GenericRecord, cur *jsonez.GoJSON, rName string,
	fDescMap wu.AFdMap, schema *avro.Schema) error {
	var err error
	var ok bool
	var fDesc *wu.AvroFieldDesc

	child := cur.Child
	for {
		if child != nil {
			/*
			 * Get the avro field descriptor correponding to
			 * the child
			 */
			fDesc, ok = fDescMap[child.Key]
			if !ok {
				errorStr := fmt.Sprintf("Field %s not found in record %s",
					child.Key, rName)
				return errors.New(errorStr)

			}

			/*
			 * If the field is an array, allocate memory for the
			 * the array entries and populate values
			 */
			if fDesc.Repeated == true {
				if child.Jsontype != jsonez.JSON_ARRAY {
					errorStr := fmt.Sprintf("Expected field %s to be an array"+
						", but got field type as %d", child.Key, child.Jsontype)
					return errors.New(errorStr)
				}

				repeatedRecords := make([]*avro.GenericRecord,
					child.GetArraySize())

				arrEntry := child.Child
				i := 0
				for {
					if arrEntry != nil {
						repeatedRecords[i] = avro.NewGenericRecord(*schema)
						if err = addAvroVal(repeatedRecords[i], arrEntry,
							rName, fDesc, schema); err != nil {
							return nil
						}
						arrEntry = arrEntry.Next
						i++
					} else {
						break
					}
				}

			} else {
				addAvroVal(record, child, rName, fDesc, schema)
			}
			child = child.Next
		} else {
			break
		}
	}

	/*
	 * Check if the record has a placeholder field
	 * and if so generate a null record for that
	 */
	if fDesc, ok = fDescMap[wu.FIELD_PLACEHOLDER]; ok {
		record.Set(wu.FIELD_PLACEHOLDER, nil)
	}

	return nil
}

/*
 * Function to parse the JSON input and build
 * an AVRO payload for further processing
 */
func avroBuildAvroOutput(jsonInput []byte, rpc string) ([]byte, error) {
	var pDesc *wu.AvroProtocolDesc
	var root *jsonez.GoJSON
	var err error

	/*
	 * Fetch the protocol descriptor and get the schema
	 */
	if pDesc, err = wu.GetProtocolDesc(rpc); err != nil {
		return nil, err
	}

	schema := avro.MustParseSchema(wu.GetAvroSchema(pDesc))
	record := avro.NewGenericRecord(schema)

	/*
	 * Walk the JSON payload and build AVRO
	 */
	if root, err = jsonez.GoJSONParse(jsonInput); err != nil {
		log.Printf("%s: Parsing JSON payload failed with error %v",
			wu.FuncName(), err)
		errorStr := fmt.Sprintf("Error %v while parsing message %s", err,
			rpc)
		return nil, errors.New(errorStr)
	}

	if err = avroWalkJSON(record, root, pDesc.Response,
		pDesc.ResponseDesc.FieldDescMap, &schema); err != nil {
		return nil, err
	}

	writer := avro.NewGenericDatumWriter()
	writer.SetSchema(schema)

	// Create a new Buffer and Encoder to write to this Buffer
	buffer := new(bytes.Buffer)
	encoder := avro.NewBinaryEncoder(buffer)

	err = writer.Write(record, encoder)
	if err != nil {
		panic(err)
	}

	return []byte(string(buffer.Bytes())), nil
}

/*
 * Function to encode kafka message key in Avro format
 */
func encodeAvroKey(msg wu.MsgFormat) ([]byte, error) {
	var key string
	var err error

	/*
	 * Build the Avro Key
	 */
	KeySchema := avro.MustParseSchema(RpcKeySchema)
	keyRecord := avro.NewGenericRecord(KeySchema)

	keyRecord.Set(wu.MSG_FORMAT_BROKER_ID, msg.BrokerId)
	keyRecord.Set(wu.MSG_FORMAT_CLIENT_ID, msg.BrokerId)
	keyRecord.Set(wu.MSG_FORMAT_RPC, msg.Rpc)

	var arr []string
	arr = append(arr, msg.IpAddress)
	keyRecord.Set(wu.MSG_FORMAT_IPADDRESS, arr)
	keyRecord.Set(wu.MSG_FORMAT_PORT, msg.Port)

	/*
	 * Fill in the metadata info
	 */
	if len(msg.Metadata) > 0 {
		subRecords := make([]*avro.GenericRecord, len(msg.Metadata))
		i := 0
		for k, v := range msg.Metadata {
			subRecord := avro.NewGenericRecord(KeySchema)
			subRecord.Set("Key", k)
			subRecord.Set("Value", v)
			subRecords[i] = subRecord
			i++
		}

		keyRecord.Set("Metadata", subRecords)
	}

	/*
	 * Generate a binary AVRO payload
	 */
	writer := avro.NewGenericDatumWriter()

	/*
	 * SetSchema must be called before calling Write
	 */
	writer.SetSchema(KeySchema)

	/*
	 *Create a new Buffer and Encoder to write to this Buffer
	 */
	buffer := new(bytes.Buffer)
	encoder := avro.NewBinaryEncoder(buffer)
	err = writer.Write(keyRecord, encoder)
	if err != nil {
		errorStr := fmt.Sprintf("Encoding Key for RPC %s failed with error %v",
			msg.Rpc, err)
		return nil, errors.New(errorStr)
	}

	key = string(buffer.Bytes())

	return []byte(key), nil
}

/*
 * Callback function to be invoked to perform processing
 * before writing the output
 */
func AvroPreWriteCB(msg wu.MsgFormat) ([]interface{}, error) {
	var err error
	var output *wk.PublishMsg
	var pubList []interface{}

	output = new(wk.PublishMsg)
	output.Topic = msg.TransactionId

	/*
	 * Build the Avro payload from JSON input
	 */
	output.Key, err = encodeAvroKey(msg)
	if err != nil {
		return nil, err
	}

	/*
	 * Build the Avro Payload from JSON
	 */
	if strings.Contains(msg.Rpc, "wedge") == false {
		if output.Payload, err = avroBuildAvroOutput([]byte(msg.Value.(string)),
			msg.Rpc); err != nil {
			errorStr := fmt.Sprintf("Building Avro output failed for RPC %s "+
				"with error %v", msg.Rpc, err)
			return nil, errors.New(errorStr)
		}
	} else {
		/*
		 * For broker related messages th payload will be a string
		 * Use WedgeMsgSchema for these purposes
		 */
		schema := avro.MustParseSchema(WedgeMsgSchema)
		msgRecord := avro.NewGenericRecord(schema)

		msgRecord.Set("Message", msg.Value.(string))

		writer := avro.NewGenericDatumWriter()
		writer.SetSchema(schema)

		/*
		 *Create a new Buffer and Encoder to write to this Buffer
		 */
		buffer := new(bytes.Buffer)
		encoder := avro.NewBinaryEncoder(buffer)
		err = writer.Write(msgRecord, encoder)
		if err != nil {
			errorStr := fmt.Sprintf("Encoding Message for RPC %s failed "+
				"with error %v", msg.Rpc, err)
			return nil, errors.New(errorStr)
		}

		output.Payload = []byte(string(buffer.Bytes()))
	}

	pubList = append(pubList, output)
	return pubList, nil
}

/*
 * Function to emit error messages
 */
func AvroErrorCB(msg *wu.MsgFormat, input interface{},
	errStr string) (interface{}, error) {

	var inputMsg *kafka.Message
	var errMsg = new(wk.PublishMsg)
	var err error

	if input != nil {
		inputMsg = input.(*kafka.Message)
		errMsg.Key = inputMsg.Key
	} else if msg != nil {
		/*
		 * Build the avro key from msg struct
		 */
		errMsg.Key, err = encodeAvroKey(*msg)
		if err != nil {
			return nil, err
		}
	}

	errMsg.Topic = wk.KAFKA_ERROR_TOPIC

	/*
	 * Build the avro error message payload
	 */
	errSchema := avro.MustParseSchema(ErrorMsgSchema)
	errRecord := avro.NewGenericRecord(errSchema)
	errRecord.Set("Message", errStr)

	/*
	 * Generate a binary AVRO payload
	 */
	writer := avro.NewGenericDatumWriter()

	/*
	 * SetSchema must be called before calling Write
	 */
	writer.SetSchema(errSchema)

	/*
	 * Create a new Buffer and Encoder to write to this Buffer
	 */
	buffer := new(bytes.Buffer)
	encoder := avro.NewBinaryEncoder(buffer)
	err = writer.Write(errRecord, encoder)
	if err != nil {
		return nil, err

	}
	/*
	 * Build the Avro payload from JSON input
	 */
	errMsg.Payload = []byte(string(buffer.Bytes()))

	return errMsg, nil
}

/*
 * JSON related CB functions
 */

/*
 * Function to validate if the oneof fields have only
 * one option set
 */
func validateRPCInput(cur *jsonez.GoJSON, fMap wu.FdMap) error {
	var child *jsonez.GoJSON
	var fDesc *wu.FieldDesc
	var curOneOfMap = make(map[string]string)
	var ok bool

	child = cur.Child

	for {
		if child != nil {
			/*
			 * Get the field descriptor for this field
			 */
			if fDesc, ok = fMap[child.Key]; !ok {
				errorStr := fmt.Sprintf("%s: Field %s not found for message %s",
					wu.FuncName(), child.Key, cur.Key)
				return errors.New(errorStr)
			}

			/*
			 * If the field belongs to a oneof, makes sure another
			 * option is not set already
			 */
			if fDesc.Ftype == wu.PROTOBUF_TYPE_MESSAGE {
				if fDesc.SubField == nil {
					errorStr := fmt.Sprintf("%s: Field %s has an empty sub "+
						"message field descriptor map", wu.FuncName(),
						child.Key)
					return errors.New(errorStr)
				}
				validateRPCInput(child, fDesc.SubField)
			} else {
				if fDesc.ParentOneof != "" && fDesc.ParentOneof != "NULL" {
					prevVal, ok := curOneOfMap[fDesc.ParentOneof]
					if ok {
						errorStr := fmt.Sprintf("%s: Oneof field %s has two options "+
							"%s and %s specified", wu.FuncName(),
							child.Key, fDesc.ParentOneof, prevVal)
						return errors.New(errorStr)
					}
				}
			}
			child = child.Next
		} else {
			break
		}
	}

	return nil
}

/*
 * JSON payload related callback functions
 */
func JSONPostReadCB(in interface{}) ([]wu.MsgFormat, error) {
	var err error
	var input *kafka.Message
	var errorStr string
	var msg wu.MsgFormat
	var root, ipAddrList, entry *jsonez.GoJSON
	var ret []wu.MsgFormat

	log.Printf("%s: post read callback invoked", wu.FuncName())

	input = in.(*kafka.Message)

	if root, err = jsonez.GoJSONParse(input.Key); err != nil {
		/*
		 * Error parsing the key, Publish the Error though
		 * the topic "_wedge.Topic_keyError"
		 */
		log.Printf("%s: Error %v while parsing key %s", wu.FuncName(),
			err, string(input.Key))
		errorStr = fmt.Sprintf("Error parsing key")
		return nil, errors.New(errorStr)
	}

	/*
	 * Fetch the IP address array
	 */
	ipAddrList, err = root.Get(wu.MSG_FORMAT_IPADDRESS)
	if err != nil {
		log.Printf("%s: Fetching Array %s returned error %v", wu.FuncName(),
			wu.MSG_FORMAT_IPADDRESS, err)
		errorStr = "Mandatory parameter IpAddress not found"
		return nil, errors.New(errorStr)
	}

	entry, err = ipAddrList.GetArrayElemByIndex(0)
	if err != nil {
		log.Printf("%s: Error %v while fetching first array element of ipAddrList",
			wu.FuncName(), err)
		errorStr = fmt.Sprintf("Error %v while fetching first array element "+
			"of ipAddrList for RPC %s", err, msg.Rpc)
		return nil, errors.New(errorStr)
	}

	ipAddLen := ipAddrList.GetArraySize()
	i := 0
	for {
		if i >= ipAddLen {
			break
		} else {
			msg.BrokerId, err = root.GetStringVal(wu.MSG_FORMAT_BROKER_ID)
			if err != nil {
				log.Printf("%s: Fetching %s returned error %v", wu.FuncName(),
					wu.MSG_FORMAT_BROKER_ID, err)
				errorStr = "Mandatory parameter BrokerId not found"
				return nil, errors.New(errorStr)
			}
			log.Printf("%s: msg.BrokerId is %s", wu.FuncName(), msg.BrokerId)

			msg.ClientId, err = root.GetStringVal(wu.MSG_FORMAT_CLIENT_ID)
			if err != nil {
				log.Printf("%s: Fetching %s returned error %v", wu.FuncName(),
					wu.MSG_FORMAT_CLIENT_ID, err)
				errorStr = "Mandatory parameter ClientId not found"
				return nil, errors.New(errorStr)
			}
			log.Printf("%s: msg.ClientId is %s", wu.FuncName(), msg.ClientId)

			msg.TransactionId, err = root.GetStringVal(wu.MSG_FORMAT_TRANSACTION_ID)
			if err != nil {
				log.Printf("%s: Fetching %s returned error %v", wu.FuncName(),
					wu.MSG_FORMAT_TRANSACTION_ID, err)

				errorStr = "Mandatory parameter TransactionId not found"
				return nil, errors.New(errorStr)
			}
			log.Printf("%s: msg.TransactionId is %s", wu.FuncName(),
				msg.TransactionId)

			msg.RpcId, err = root.GetStringVal(wu.MSG_FORMAT_RPC_ID)
			if err != nil {
				log.Printf("%s: Fetching %s returned error %v", wu.FuncName(),
					wu.MSG_FORMAT_RPC_ID, err)

				errorStr = "Mandatory parameter RpcId not found"
				return nil, errors.New(errorStr)

			}

			msg.IpAddress = entry.Valstr
			if err != nil {
				log.Printf("%s: Fetching %s returned error %v", wu.FuncName(),
					wu.MSG_FORMAT_IPADDRESS, err)

				errorStr = "Mandatory parameter IpAddress not found"
				return nil, errors.New(errorStr)
			}
			log.Printf("%s: msg.IpAddress is %s", wu.FuncName(), msg.IpAddress)

			msg.Port, err = root.GetStringVal(wu.MSG_FORMAT_PORT)
			if err != nil {
				log.Printf("%s: Fetching %s returned error %v", wu.FuncName(),
					wu.MSG_FORMAT_PORT, err)

				errorStr = "Mandatory parameter Port not found"
				return nil, errors.New(errorStr)
			}
			log.Printf("%s: msg.Port is %s", wu.FuncName(), msg.Port)

			metadata, err := root.Get(wu.MSG_FORMAT_METADATA)
			if metadata != nil && err != nil {
				for entry := metadata.Child; entry != nil; entry = entry.Next {
					msg.Metadata[entry.Child.Key] = entry.Child.Valstr
				}
			}

			/*
			 * If the clientId of the key matches wedge's clientId,
			 * then drop it as it a copy of the message that was
			 * produced by wedge
			 */
			if strings.Compare(msg.BrokerId, msg.ClientId) != 0 &&
				strings.Compare(msg.BrokerId, msg.ClientId) != 0 {
				msg.Rpc = strings.Replace(*input.TopicPartition.Topic,
					"_", "/", -1)

				/*
				 * If the RPC is for a call cancellation, then
				 * change the value i.e the RPC to be terminated
				 * to the required format
				 */
				if strings.Compare(msg.Rpc, wu.CALL_CANCELLATION_RPC) == 0 {
					msg.Value = strings.Replace(string(input.Value), "_", "/",
						-1)
				} else {
					msg.Value = string(input.Value)

					/*
					 * Verify if the oneof fields have only one option
					 * specified and report an error if it's not so
					 */
					if strings.Contains(msg.Rpc, "wedge") == false {
						rdesc, err := wu.GetRpcDesc(msg.Rpc)

						if err != nil {
							errorStr := fmt.Sprintf("%s: Could not fetch RPC "+
								"descriptor for %s", wu.FuncName(), msg.Rpc)
							return nil, errors.New(errorStr)
						}

						/*
						 * Parse the JSON data into goJSON objects
						 */
						if len(input.Value) > 0 {
							root, err := jsonez.GoJSONParse(input.Value)
							if err != nil {
								errorStr := fmt.Sprintf("%s: JSON input parsing "+
									"failed with error %v", wu.FuncName(), err)
								return nil, errors.New(errorStr)
							}

							validateRPCInput(root,
								rdesc.InMsgDescriptor.FieldDescMap)
						}
					}
				}

				ret = append(ret, msg)

			}
		}

		i++
		entry = entry.Next
	}

	return ret, nil
}

/*
 * Function to escapte special characters for influx DB
 */
func escInfluxSpecialChars(input string) string {
	var output []byte

	for _, c := range input {
		switch c {
		case ',':
			output = append(output, '\\', ',')
		case '=':
			output = append(output, '\\', '=')
		case ' ':
			output = append(output, '\\', ' ')
		default:
			output = append(output, byte(rune(c)))
		}
	}

	return string(output)
}

/*
 * Function to build the kafka key in JSON format from MsgFormat
 * structure
 */
func GetKafKaJSONKey(msg wu.MsgFormat) ([]byte, error) {
	var err error
	var errorStr string

	/*
	 * Make a Key for the message
	 */
	root := jsonez.AllocObject()
	/*
	 * Set the client address to that of wedge kafka
	 * producer
	 */
	root.AddVal(msg.BrokerId, wu.MSG_FORMAT_BROKER_ID)
	root.AddVal(msg.BrokerId, wu.MSG_FORMAT_CLIENT_ID)
	root.AddVal(msg.Rpc, "Rpc")
	root.AddVal(msg.IpAddress, wu.MSG_FORMAT_IPADDRESS)
	root.AddVal(msg.Port, wu.MSG_FORMAT_PORT)

	if len(msg.Metadata) > 0 {
		for k, v := range msg.Metadata {
			obj := jsonez.AllocObject()
			obj.AddEntryToObject(k, jsonez.AllocString(v))
			err = root.AddToArray(obj, wu.MSG_FORMAT_METADATA)
			if err != nil {
				log.Printf("%s: Adding Metadata with key %s and "+
					"value %s failed with error %v\n", wu.FuncName(),
					k, v, err)
				errorStr = fmt.Sprintf("Adding Metadata with key %s and "+
					"value %s failed with error %v\n", k, v, err)
				return nil, errors.New(errorStr)
			}
		}
	}

	log.Printf("%s: Resultant JSON is %s", wu.FuncName(),
		jsonez.GoJSONPrint(root))

	return jsonez.GoJSONPrint(root), nil
}

/*
 * Callback function to publish a JSON payload to a kafka bus
 */
func JSONPreWriteCB(msg wu.MsgFormat) ([]interface{}, error) {
	//var fluentRoot, mRoot, entry, val *jsonez.GoJSON
	var err error
	//var key, prefix, orig_prefix, errorStr string
	var pubList []interface{}
	var pub *wk.PublishMsg
	//	var i int

	pub = new(wk.PublishMsg)
	pub.Topic = msg.TransactionId
	pub.Payload = []byte(msg.Value.(string))

	/*
	 * Build the key for kafka message
	 */
	if pub.Key, err = GetKafKaJSONKey(msg); err != nil {
		return nil, err
	}

	/*
	 * Kafka doesn't accept '/'' in topic names. So replace
	 * all '/' with '_'
	 */
	pub.Topic = strings.Replace(pub.Topic, "/", "_", -1)

	pubList = append(pubList, pub)

	return pubList, nil
}

/*
 * Callback function to publish a influx DB payload to a
 * kafka bus for openconfig telemetry data
 */
func InfluxPreWriteCB(msg wu.MsgFormat) ([]interface{}, error) {
	var mRoot, entry, val *jsonez.GoJSON
	var err error
	var key, prefix, orig_prefix, errorStr string
	var pubList []interface{}
	var i int

	if strings.Compare(msg.Rpc,
		"/telemetry.OpenConfigTelemetry/telemetrySubscribe") != 0 {
		log.Printf("%s: RPC is not telemetrySubscribe but %s", wu.FuncName(),
			msg.Rpc)
		errorStr = fmt.Sprintf("Callback %s is called for a non-telemetry RPC %s",
			wu.FuncName(), msg.Rpc)
		return nil, errors.New(errorStr)
	}

	/*
	 * Publish to the influx topic with value in
	 * influx line protocol format
	 */
	if mRoot, err = jsonez.GoJSONParse([]byte(msg.Value.(string))); err != nil {
		/*
		 * Error parsing the key, Publish the Error though
		 * the topic "_wedge.Topic_keyError"
		 */
		log.Printf("%s: Error %v while parsing message %s", wu.FuncName(),
			err, msg.Value.(string))
		errorStr = fmt.Sprintf("Error %v while parsing message %s", err,
			msg.Rpc)
		return nil, errors.New(errorStr)
	}

	systemId, err := mRoot.GetStringVal("system_id")
	if err != nil {
		log.Printf("%s: Error %v while fetching system_id", wu.FuncName(),
			err)
		errorStr = fmt.Sprintf("Error %v while fetching system_id for RPC %s",
			err, msg.Rpc)
		return nil, errors.New(errorStr)
	}
	systemId = escInfluxSpecialChars(systemId)

	compId, err := mRoot.GetStringVal("component_id")
	if err != nil {
		compId = "0"
	}
	compId = escInfluxSpecialChars(compId)

	path, err := mRoot.GetStringVal("path")
	if err != nil {
		log.Printf("%s: Error %v while fetching system_id", wu.FuncName(),
			err)
		errorStr = fmt.Sprintf("Error %v while fetching system_id for RPC %s",
			err, msg.Rpc)
		return nil, errors.New(errorStr)
	}
	path = escInfluxSpecialChars(path)

	seqno, err := mRoot.GetUIntVal("sequence_number")
	if err != nil {
		seqno = 0
	}

	timestamp, err := mRoot.GetUIntVal("timestamp")
	if err != nil {
		log.Printf("%s: Error %v while fetching sequence_number", wu.FuncName(),
			err)
		errorStr = fmt.Sprintf("Error %v while fetching sequence_number for "+
			"RPC %s", err, msg.Rpc)
		return nil, errors.New(errorStr)
	}
	timestamp *= 1e6

	/*
	 * Fetch the key value pairs array
	 */
	kv, err := mRoot.Get("kv")
	if err != nil {
		log.Printf("%s: Error %v while fetching kv", wu.FuncName(),
			err)
		errorStr = fmt.Sprintf("Error %v while fetching kv for "+
			"RPC %s", err, msg.Rpc)
		return nil, errors.New(errorStr)
	}

	/*
	 * walk the key value pairs and build the influx messages
	 * accordingly
	 */
	entry, err = kv.GetArrayElemByIndex(0)
	if err != nil {
		log.Printf("%s: Error %v while fetching first array element of kv",
			wu.FuncName(), err)
		errorStr = fmt.Sprintf("Error %v while fetching first array element "+
			"of kv for RPC %s", err, msg.Rpc)
		return nil, errors.New(errorStr)
	}

	kvLength := kv.GetArraySize()
	var influxPub *wk.PublishMsg
	var measurement, influxMsg string
	for {
		if i >= kvLength {
			break
		} else {
			if entry != nil {
				key, err = entry.GetStringVal("key")
				orig_key := key
				key = escInfluxSpecialChars(key)

				if key != "__timestamp__" {
					if key == "__prefix__" {
						/*
						 * If influxPub is nil, then this is the first
						 * entry in the message received
						 */
						if influxPub != nil && influxMsg != "" &&
							measurement != "" {
							/*
							 * Append the previous entry to the list
							 */
							influxPub.Payload = []byte(measurement + influxMsg +
								" " + strconv.FormatUint(timestamp, 10))
							influxPub.Key = []byte("")
							pubList = append(pubList, influxPub)
							influxPub = nil
							influxMsg = ""
							measurement = ""

						}

						sval, err := entry.GetStringVal("str_value")
						if err != nil {
							log.Printf("%s: Error %v while fetching value of "+
								"key %s", wu.FuncName(), err, key)
							errorStr = fmt.Sprintf("Error %v while fetching "+
								"value of key %s for RPC %s", err, key, msg.Rpc)
							return nil, errors.New(errorStr)
						}

						orig_prefix = sval
						prefix = escInfluxSpecialChars(sval)
						influxPub = new(wk.PublishMsg)
						influxPub.Topic = InFluxTopic
						measurement = systemId + "," + "component_id=" + compId +
							",path=\"" + path + "\",sequence_number=" +
							func() string {
								if seqno > 0 {
									return strconv.FormatUint(seqno, 10)
								} else {
									return "0"
								}
							}() + " "
					} else {
						/*
						 * Search for the accompanying *_value field
						 * and populate the contents accordingly
						 */
						val, err = entry.Get("double_value")
						if err != nil {
							if !strings.Contains(err.Error(), "Path not found") {
								log.Printf("%s: Error %v while fetching value "+
									"of key double_value", wu.FuncName(), err)
								errorStr = fmt.Sprintf("Error %v while "+
									"fetching value of key double_value for "+
									"RPC %s", err, msg.Rpc)
								return nil, errors.New(errorStr)
							}
						} else {
							if influxMsg == "" {
								influxMsg = prefix + key + "="
							} else {
								influxMsg += "," + prefix + key + "="
							}
							influxMsg += strconv.FormatFloat(val.Valdouble,
								'E', -1, 64)

							fv := jsonez.AllocObject()
							fv.AddVal(val.Valdouble, orig_prefix+orig_key)
							goto ITERATE
						}

						val, err = entry.Get("int_value")
						if err != nil {
							if !strings.Contains(err.Error(), "Path not found") {
								log.Printf("%s: Error %v while fetching value "+
									"of key int_value", wu.FuncName(), err)
								errorStr = fmt.Sprintf("Error %v while "+
									"fetching value of key int_value for "+
									"RPC %s", err, msg.Rpc)
								return nil, errors.New(errorStr)
							}
						} else {
							if influxMsg == "" {
								influxMsg = prefix + key + "="
							} else {
								influxMsg += "," + prefix + key + "="
							}

							if val.Jsontype == jsonez.JSON_UINT {
								influxMsg += strconv.FormatUint(val.Valuint,
									10) + "i"

								fv := jsonez.AllocObject()
								fv.AddVal(val.Valuint, orig_prefix+orig_key)
							} else if val.Jsontype == jsonez.JSON_INT {
								influxMsg += strconv.FormatInt(val.Valint,
									10) + "i"

								fv := jsonez.AllocObject()
								fv.AddVal(val.Valint, orig_prefix+orig_key)
								fv.Key = orig_prefix + orig_key
							} else {
								log.Printf("%s: Error %v while fetching value "+
									"of key int_value", wu.FuncName(), err)
								errorStr = fmt.Sprintf("Error %v while "+
									"fetching value of key int_value for "+
									"RPC %s", err, msg.Rpc)
								return nil, errors.New(errorStr)
							}

							goto ITERATE
						}

						val, err = entry.Get("sint_value")
						if err != nil {
							if !strings.Contains(err.Error(), "Path not found") {
								log.Printf("%s: Error %v while fetching value "+
									"of key sint_value", wu.FuncName(), err)
								errorStr = fmt.Sprintf("Error %v while "+
									"fetching value of key int_value for "+
									"RPC %s", err, msg.Rpc)
								return nil, errors.New(errorStr)
							}
						} else {
							if influxMsg == "" {
								influxMsg = prefix + key + "="
							} else {
								influxMsg += "," + prefix + key + "="
							}

							if val.Jsontype == jsonez.JSON_UINT {
								influxMsg += strconv.FormatUint(val.Valuint,
									10) + "i"

								fv := jsonez.AllocObject()
								fv.AddVal(val.Valuint, orig_prefix+orig_key)
							} else if val.Jsontype == jsonez.JSON_INT {
								influxMsg += strconv.FormatInt(val.Valint,
									10) + "i"

								fv := jsonez.AllocObject()
								fv.AddVal(val.Valint, orig_prefix+orig_key)

								fv.Key = orig_prefix + orig_key
							} else {
								log.Printf("%s: Error %v while fetching value "+
									"of key sint_value", wu.FuncName(), err)
								errorStr = fmt.Sprintf("Error %v while "+
									"fetching value of key sint_value for "+
									"RPC %s", err, msg.Rpc)
								return nil, errors.New(errorStr)
							}

							goto ITERATE
						}

						val, err = entry.Get("uint_value")
						if err != nil {
							if !strings.Contains(err.Error(), "Path not found") {
								log.Printf("%s: Error %v while fetching value "+
									"of key uint_value", wu.FuncName(), err)
								errorStr = fmt.Sprintf("Error %v while "+
									"fetching value of key uint_value for "+
									"RPC %s", err, msg.Rpc)
								return nil, errors.New(errorStr)
							}
						} else {
							if influxMsg == "" {
								influxMsg = prefix + key + "="
							} else {
								influxMsg += "," + prefix + key + "="
							}
							influxMsg += strconv.FormatUint(val.Valuint,
								10) + "i"

							fv := jsonez.AllocObject()
							fv.AddVal(val.Valuint, orig_prefix+orig_key)

							fv.Key = orig_prefix + orig_key

							goto ITERATE
						}

						val, err = entry.Get("bool_value")
						if err != nil {
							if !strings.Contains(err.Error(), "Path not found") {
								log.Printf("%s: Error %v while fetching value "+
									"of key bool_value", wu.FuncName(), err)
								errorStr = fmt.Sprintf("Error %v while "+
									"fetching value of key bool_value for "+
									"RPC %s", err, msg.Rpc)
								return nil, errors.New(errorStr)
							}
						} else {
							if influxMsg == "" {
								influxMsg = prefix + key + "="
							} else {
								influxMsg += "," + prefix + key + "="
							}
							if val.Valbool == true {
								influxMsg += "true"
							} else {
								influxMsg += "false"
							}

							fv := jsonez.AllocObject()
							fv.AddVal(val.Valbool, orig_prefix+orig_key)

							goto ITERATE
						}

						val, err = entry.Get("str_value")
						if err != nil {
							if !strings.Contains(err.Error(), "Path not found") {
								log.Printf("%s: Error %v while fetching value "+
									"of key string_value", wu.FuncName(), err)
								errorStr = fmt.Sprintf("Error %v while "+
									"fetching value of key string_value for "+
									"RPC %s", err, msg.Rpc)
								return nil, errors.New(errorStr)
							}
						} else {
							if influxMsg == "" {
								influxMsg = prefix + key + "=\"" + val.Valstr +
									"\""
							} else {
								influxMsg += "," + prefix + key + "=\"" +
									val.Valstr + "\""
							}

							fv := jsonez.AllocObject()
							fv.AddVal(val.Valstr, orig_prefix+orig_key)

							goto ITERATE
						}

						val, err = entry.Get("bytes_value")
						if err != nil {
							if !strings.Contains(err.Error(), "Path not found") {
								log.Printf("%s: Error %v while fetching value "+
									"of key string_value", wu.FuncName(), err)
								errorStr = fmt.Sprintf("Error %v while "+
									"fetching value of key string_value for "+

									"RPC %s", err, msg.Rpc)
								return nil, errors.New(errorStr)
							}
						} else {
							if influxMsg == "" {
								influxMsg = prefix + key + "=\"" + val.Valstr +
									"\""
							} else {
								influxMsg += "," + prefix + key + "=\"" +
									val.Valstr + "\""
							}

							fv := jsonez.AllocObject()
							fv.AddVal(val.Valstr, orig_prefix+orig_key)

							goto ITERATE
						}
					}
				}
			}
		ITERATE:
			i++
			entry = entry.Next
		}
	}

	/*
	 * To publish the final influx message
	 */
	if influxPub != nil && influxMsg != "" &&
		measurement != "" {
		/*
		 * Append the previous entry to the list
		 */
		influxPub.Payload = []byte(measurement + influxMsg + " " +
			strconv.FormatUint(timestamp, 10))
		influxPub.Key = []byte("")
		pubList = append(pubList, influxPub)
	}

	return pubList, nil

}

/*
 * Callback function to publish a fluentd json payload to a
 * kafka bus for openconfig telemetry data
 */
func FluentdPreWriteCB(msg wu.MsgFormat) ([]interface{}, error) {
	var fluentRoot, mRoot, entry, val *jsonez.GoJSON
	var err error
	var key, orig_prefix, errorStr string
	var pubList []interface{}
	var i int

	if strings.Compare(msg.Rpc,
		"/telemetry.OpenConfigTelemetry/telemetrySubscribe") != 0 {
		log.Printf("%s: RPC is not telemetrySubscribe but %s", wu.FuncName(),
			msg.Rpc)
		errorStr = fmt.Sprintf("Callback %s is called for a non-telemetry RPC %s",
			wu.FuncName(), msg.Rpc)
		return nil, errors.New(errorStr)
	}

	/*
	 * Publish to the fluentd topic with value in json
	 */
	if mRoot, err = jsonez.GoJSONParse([]byte(msg.Value.(string))); err != nil {
		/*
		   * Error parsing the key, Publish the Error though
		     * the topic "_wedge.Topic_keyError"
		*/
		log.Printf("%s: Error %v while parsing message %s", wu.FuncName(),
			err, msg.Value.(string))
		errorStr = fmt.Sprintf("Error %v while parsing message %s", err,
			msg.Rpc)
		return nil, errors.New(errorStr)
	}

	fluentRoot = jsonez.AllocObject()

	systemId, err := mRoot.GetStringVal("system_id")
	if err != nil {
		log.Printf("%s: Error %v while fetching system_id", wu.FuncName(),
			err)
		errorStr = fmt.Sprintf("Error %v while fetching system_id for RPC %s",
			err, msg.Rpc)
		return nil, errors.New(errorStr)
	}
	if fluentRoot.AddVal(systemId, "system_id"); err != nil {
		log.Printf("%s: Error %v while adding system_id to fluentMsg",
			wu.FuncName(), err)
		errorStr = fmt.Sprintf("Error %v while adding system_id to "+
			"fluentMsg for RPC %s", err, msg.Rpc)
		return nil, errors.New(errorStr)
	}

	compId, err := mRoot.GetStringVal("component_id")
	if err != nil {
		compId = "0"
	}
	if err = fluentRoot.AddVal(compId, "component_id"); err != nil {
		log.Printf("%s: Error %v while adding component_id to fluentMsg",
			wu.FuncName(), err)
		errorStr = fmt.Sprintf("Error %v while adding component_id to "+
			"fluentMsg for RPC %s", err, msg.Rpc)
		return nil, errors.New(errorStr)
	}

	path, err := mRoot.GetStringVal("path")
	if err != nil {
		log.Printf("%s: Error %v while fetching system_id", wu.FuncName(),
			err)
		errorStr = fmt.Sprintf("Error %v while fetching system_id for RPC %s",
			err, msg.Rpc)
		return nil, errors.New(errorStr)
	}
	if err = fluentRoot.AddVal(path, "path"); err != nil {
		log.Printf("%s: Error %v while adding path to fluentMsg",
			wu.FuncName(), err)
		errorStr = fmt.Sprintf("Error %v while adding path to "+
			"fluentMsg for RPC %s", err, msg.Rpc)
		return nil, errors.New(errorStr)
	}

	seqno, err := mRoot.GetUIntVal("sequence_number")
	if err != nil {
		seqno = 0
	}
	if err = fluentRoot.AddVal(seqno, "sequence_number"); err != nil {
		log.Printf("%s: Error %v while adding sequence_number to fluentMsg",
			wu.FuncName(), err)
		errorStr = fmt.Sprintf("Error %v while adding sequence_number to "+
			"fluentMsg for RPC %s", err, msg.Rpc)
		return nil, errors.New(errorStr)
	}

	timestamp, err := mRoot.GetUIntVal("timestamp")
	if err != nil {
		log.Printf("%s: Error %v while fetching sequence_number", wu.FuncName(),
			err)
		errorStr = fmt.Sprintf("Error %v while fetching sequence_number for "+
			"RPC %s", err, msg.Rpc)
		return nil, errors.New(errorStr)
	}
	if err = fluentRoot.AddVal(timestamp, "time"); err != nil {
		log.Printf("%s: Error %v while adding time to fluentMsg",
			wu.FuncName(), err)
		errorStr = fmt.Sprintf("Error %v while adding time to "+
			"fluentMsg for RPC %s", err, msg.Rpc)
		return nil, errors.New(errorStr)
	}
	timestamp *= 1e6

	/*
	 * Fetch the key value pairs array
	 */
	kv, err := mRoot.Get("kv")
	if err != nil {
		log.Printf("%s: Error %v while fetching kv", wu.FuncName(),
			err)
		errorStr = fmt.Sprintf("Error %v while fetching kv for "+
			"RPC %s", err, msg.Rpc)
		return nil, errors.New(errorStr)
	}

	/*
	 * Walk the key value pairs and build the influx messages
	 * accordingly
	 */
	entry, err = kv.GetArrayElemByIndex(0)
	if err != nil {
		log.Printf("%s: Error %v while fetching first array element of kv",
			wu.FuncName(), err)
		errorStr = fmt.Sprintf("Error %v while fetching first array element "+
			"of kv for RPC %s", err, msg.Rpc)
		return nil, errors.New(errorStr)
	}

	sensors := jsonez.AllocArray()
	fluentRoot.AddEntryToObject("sensors", sensors)

	kvLength := kv.GetArraySize()

	for {
		if i >= kvLength {
			break
		} else {
			if entry != nil {
				key, err = entry.GetStringVal("key")
				orig_key := key
				key = escInfluxSpecialChars(key)

				if key != "__timestamp__" {
					/*
					 * Search for the accompanying *_value field
					 * and populate the contents accordingly
					 */
					val, err = entry.Get("double_value")
					if err != nil {
						if !strings.Contains(err.Error(), "Path not found") {
							log.Printf("%s: Error %v while fetching value "+
								"of key double_value", wu.FuncName(), err)
							errorStr = fmt.Sprintf("Error %v while "+
								"fetching value of key double_value for "+
								"RPC %s", err, msg.Rpc)
							return nil, errors.New(errorStr)
						}
					} else {
						fv := jsonez.AllocObject()
						fv.AddVal(val.Valdouble, orig_prefix+orig_key)

						if sensors.AddEntryToArray(fv); err != nil {
							log.Printf("%s: Error %v while adding %s to "+
								"sensors in fluentMsg", wu.FuncName(),
								orig_prefix+orig_key, err)
							errorStr = fmt.Sprintf("Error %v while adding "+
								"%s to sensors in fluentMsg for RPC %s",
								err, orig_prefix+orig_key, msg.Rpc)
							return nil, errors.New(errorStr)
						}

						goto ITERATE
					}

					val, err = entry.Get("int_value")
					if err != nil {
						if !strings.Contains(err.Error(), "Path not found") {
							log.Printf("%s: Error %v while fetching value "+
								"of key int_value", wu.FuncName(), err)
							errorStr = fmt.Sprintf("Error %v while "+
								"fetching value of key int_value for "+
								"RPC %s", err, msg.Rpc)
							return nil, errors.New(errorStr)
						}
					} else {
						if val.Jsontype == jsonez.JSON_UINT {

							fv := jsonez.AllocObject()
							fv.AddVal(val.Valuint, orig_prefix+orig_key)

							if sensors.AddEntryToArray(fv); err != nil {
								log.Printf("%s: Error %v while adding %s to "+
									"sensors in fluentMsg", wu.FuncName(),
									orig_prefix+orig_key, err)
								errorStr = fmt.Sprintf("Error %v while adding "+
									"%s to sensors in fluentMsg for RPC %s",
									err, orig_prefix+orig_key, msg.Rpc)
								return nil, errors.New(errorStr)
							}
						} else if val.Jsontype == jsonez.JSON_INT {
							fv := jsonez.AllocObject()
							fv.AddVal(val.Valint, orig_prefix+orig_key)

							fv.Key = orig_prefix + orig_key
							if sensors.AddEntryToArray(fv); err != nil {
								log.Printf("%s: Error %v while adding %s to "+
									"sensors in fluentMsg", wu.FuncName(),
									orig_prefix+orig_key, err)
								errorStr = fmt.Sprintf("Error %v while adding "+
									"%s to sensors in fluentMsg for RPC %s",
									err, orig_prefix+orig_key, msg.Rpc)
								return nil, errors.New(errorStr)
							}
						} else {
							log.Printf("%s: Error %v while fetching value "+
								"of key int_value", wu.FuncName(), err)
							errorStr = fmt.Sprintf("Error %v while "+
								"fetching value of key int_value for "+
								"RPC %s", err, msg.Rpc)
							return nil, errors.New(errorStr)
						}

						goto ITERATE
					}

					val, err = entry.Get("sint_value")
					if err != nil {
						if !strings.Contains(err.Error(), "Path not found") {
							log.Printf("%s: Error %v while fetching value "+
								"of key sint_value", wu.FuncName(), err)
							errorStr = fmt.Sprintf("Error %v while "+
								"fetching value of key int_value for "+
								"RPC %s", err, msg.Rpc)
							return nil, errors.New(errorStr)
						}
					} else {
						if val.Jsontype == jsonez.JSON_UINT {
							fv := jsonez.AllocObject()
							fv.AddVal(val.Valuint, orig_prefix+orig_key)

							if sensors.AddEntryToArray(fv); err != nil {
								log.Printf("%s: Error %v while adding %s to "+
									"sensors in fluentMsg", wu.FuncName(),
									orig_prefix+orig_key, err)
								errorStr = fmt.Sprintf("Error %v while adding "+
									"%s to sensors in fluentMsg for RPC %s",
									err, orig_prefix+orig_key, msg.Rpc)
								return nil, errors.New(errorStr)
							}
						} else if val.Jsontype == jsonez.JSON_INT {
							fv := jsonez.AllocObject()
							fv.AddVal(val.Valint, orig_prefix+orig_key)

							fv.Key = orig_prefix + orig_key

							if sensors.AddEntryToArray(fv); err != nil {
								log.Printf("%s: Error %v while adding %s to "+
									"sensors in fluentMsg", wu.FuncName(),
									orig_prefix+orig_key, err)
								errorStr = fmt.Sprintf("Error %v while adding "+
									"%s to sensors in fluentMsg for RPC %s",
									err, orig_prefix+orig_key, msg.Rpc)
								return nil, errors.New(errorStr)
							}
						} else {
							log.Printf("%s: Error %v while fetching value "+
								"of key sint_value", wu.FuncName(), err)
							errorStr = fmt.Sprintf("Error %v while "+
								"fetching value of key sint_value for "+
								"RPC %s", err, msg.Rpc)
							return nil, errors.New(errorStr)
						}

						goto ITERATE
					}

					val, err = entry.Get("uint_value")
					if err != nil {
						if !strings.Contains(err.Error(), "Path not found") {
							log.Printf("%s: Error %v while fetching value "+
								"of key uint_value", wu.FuncName(), err)
							errorStr = fmt.Sprintf("Error %v while "+
								"fetching value of key uint_value for "+
								"RPC %s", err, msg.Rpc)
							return nil, errors.New(errorStr)
						}
					} else {
						fv := jsonez.AllocObject()
						fv.AddVal(val.Valuint, orig_prefix+orig_key)

						fv.Key = orig_prefix + orig_key

						if sensors.AddEntryToArray(fv); err != nil {
							log.Printf("%s: Error %v while adding %s to "+
								"sensors in fluentMsg", wu.FuncName(),
								orig_prefix+orig_key, err)
							errorStr = fmt.Sprintf("Error %v while adding "+
								"%s to sensors in fluentMsg for RPC %s",
								err, orig_prefix+orig_key, msg.Rpc)
							return nil, errors.New(errorStr)
						}

						goto ITERATE
					}

					val, err = entry.Get("bool_value")
					if err != nil {
						if !strings.Contains(err.Error(), "Path not found") {
							log.Printf("%s: Error %v while fetching value "+
								"of key bool_value", wu.FuncName(), err)
							errorStr = fmt.Sprintf("Error %v while "+
								"fetching value of key bool_value for "+
								"RPC %s", err, msg.Rpc)
							return nil, errors.New(errorStr)
						}
					} else {
						fv := jsonez.AllocObject()
						fv.AddVal(val.Valbool, orig_prefix+orig_key)

						if sensors.AddEntryToArray(fv); err != nil {
							log.Printf("%s: Error %v while adding %s to "+
								"sensors in fluentMsg", wu.FuncName(),
								orig_prefix+orig_key, err)
							errorStr = fmt.Sprintf("Error %v while adding "+
								"%s to sensors in fluentMsg for RPC %s",
								err, orig_prefix+orig_key, msg.Rpc)
							return nil, errors.New(errorStr)
						}

						goto ITERATE
					}

					val, err = entry.Get("str_value")
					if err != nil {
						if !strings.Contains(err.Error(), "Path not found") {
							log.Printf("%s: Error %v while fetching value "+
								"of key string_value", wu.FuncName(), err)
							errorStr = fmt.Sprintf("Error %v while "+
								"fetching value of key string_value for "+
								"RPC %s", err, msg.Rpc)
							return nil, errors.New(errorStr)
						}
					} else {
						fv := jsonez.AllocObject()
						fv.AddVal(val.Valstr, orig_prefix+orig_key)

						if sensors.AddEntryToArray(fv); err != nil {
							log.Printf("%s: Error %v while adding %s to "+
								"sensors in fluentMsg", wu.FuncName(),
								orig_prefix+orig_key, err)
							errorStr = fmt.Sprintf("Error %v while adding "+
								"%s to sensors in fluentMsg for RPC %s",
								err, orig_prefix+orig_key, msg.Rpc)
							return nil, errors.New(errorStr)
						}

						goto ITERATE
					}

					val, err = entry.Get("bytes_value")
					if err != nil {
						if !strings.Contains(err.Error(), "Path not found") {
							log.Printf("%s: Error %v while fetching value "+
								"of key string_value", wu.FuncName(), err)
							errorStr = fmt.Sprintf("Error %v while "+
								"fetching value of key string_value for "+

								"RPC %s", err, msg.Rpc)
							return nil, errors.New(errorStr)
						}
					} else {
						fv := jsonez.AllocObject()
						fv.AddVal(val.Valstr, orig_prefix+orig_key)

						if sensors.AddEntryToArray(fv); err != nil {
							log.Printf("%s: Error %v while adding %s to "+
								"sensors in fluentMsg", wu.FuncName(),
								orig_prefix+orig_key, err)
							errorStr = fmt.Sprintf("Error %v while adding "+
								"%s to sensors in fluentMsg for RPC %s",
								err, orig_prefix+orig_key, msg.Rpc)
							return nil, errors.New(errorStr)
						}

						goto ITERATE
					}

				}
			}
		ITERATE:
			i++
			entry = entry.Next
		}
	}

	/*
	 * Add the fluentd message to pubList
	 */
	fluentdPub := new(wk.PublishMsg)
	fluentdPub.Topic = FluentdTopic
	fluentdPub.Key = []byte("")
	fluentdPub.Payload = jsonez.GoJSONPrint(fluentRoot)
	pubList = append(pubList, fluentdPub)

	return pubList, nil

}

/*
 * Callback function to process telemetry data by invoking fluent
 * and influx processing functions
 */
func TelemetryPreWriteCB(msg wu.MsgFormat) ([]interface{}, error) {
	var pubList []interface{}
	var err error

	if strings.Compare(msg.Rpc,
		"/telemetry.OpenConfigTelemetry/telemetrySubscribe") != 0 {
		return JSONPreWriteCB(msg)
	}

	jsonPub, err := JSONPreWriteCB(msg)
	if err != nil {
            return nil, err
	}
	for _, val := range jsonPub {
            pubList = append(pubList, val)
	}

	influxPub, err := InfluxPreWriteCB(msg)
	if err != nil {
            return nil, err
	}
	for _, val := range influxPub {
            pubList = append(pubList, val)
	}

	fluentPub, err := FluentdPreWriteCB(msg)
	if err != nil {
            return nil, err
	}
	for _, val := range fluentPub {
            pubList = append(pubList, val)
	}

	return pubList, nil
}
