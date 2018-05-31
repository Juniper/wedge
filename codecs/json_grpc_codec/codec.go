/*
 * Copyright (c) 2018, Juniper Networks, Inc.
 * All rights reserved.
 */

package json_grpc_codec

import wu "git.juniper.net/sksubra/wedge/util"

type JsonGrpcCodec struct{}

type JsonCoder struct {
	jsonPayload  string
	inDesc       *wu.MsgDesc
	outDesc      *wu.MsgDesc
	inStreamApi  bool
	outStreamAPi bool
	target       string
}

/**
 * Method to encode JSON to protobuf
 */
func (JsonGrpcCodec) Marshal(v interface{}) ([]byte, error) {
	return JsonToProto(v.(*JsonCoder))
}

/**
 * Method to decode protobuf to JSON
 */
func (JsonGrpcCodec) Unmarshal(data []byte, v interface{}) error {
	return ProtoToJson(data, v.(*JsonCoder))
}

/**
 * Method to return the name of codec
 */
func (JsonGrpcCodec) String() string {
	return "JsonGrpcCodec"
}

/*
 * Method to build a web coder object based
 * on the rpc name and input json string
 */
func (JsonGrpcCodec) EncodeInput(rpc string, input interface{}) (interface{},
	error) {

	rdesc, err := wu.GetRpcDesc(rpc)

	if err != nil {
		return nil, err
	}

	w := new(JsonCoder)

	w.jsonPayload = input.(string)
	w.inDesc = rdesc.InMsgDescriptor
	w.outDesc = rdesc.OutMsgDescriptor
	w.target = rdesc.Target
	w.inStreamApi = rdesc.InStreamAPI
	w.outStreamAPi = rdesc.OutStreamAPI

	return w, nil
}

/*
 * Method to Get the JSON output of an RPC execution
 * with a JsonCoder object as input
 */
func (JsonGrpcCodec) DecodeOutput(codedOutput interface{}) interface{} {
	val := codedOutput.(*JsonCoder)
	return val.jsonPayload
}

/*
 * Method to Create an output codec for an object
 */
func (JsonGrpcCodec) CreateRPCOutputObj(input interface{}) interface{} {
	in := input.(*JsonCoder)
	out := new(JsonCoder)

	out.inDesc = in.inDesc
	out.outDesc = in.outDesc
	out.inStreamApi = in.inStreamApi
	out.outStreamAPi = in.outStreamAPi
	out.target = in.target

	return out
}

/*
 * Method to indicate if the payload is empty as per the
 * codec format
 */
func (JsonGrpcCodec) IsEmpty(codedVal interface{}) bool {
	val := codedVal.(*JsonCoder)
	return len(val.jsonPayload) == 0
}

/*
 * Function to create JsonGrpcCodec object
 */
func CreateJsonGrpcCodecObject() wu.GenCodec {
	codec := new(JsonGrpcCodec)
	return codec
}
