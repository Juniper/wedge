/*
 * Copyright (c) 2018, Juniper Networks, Inc.
 * All rights reserved.
 */

package json_grpc_codec

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"log"

	wu "git.juniper.net/sksubra/wedge/util"
	"github.com/srikanth2212/jsonez"
)

func ProtoToJson(data []byte, w *JsonCoder) error {
	var root *jsonez.GoJSON
	var err error

	root = jsonez.AllocObject()

	if len(data) <= 0 {
		return nil
	}

	err = w.parseProtoInput(root, data, w.outDesc.FieldDescMap)

	if err != nil {
		return err
	}

	output := jsonez.GoJSONPrint(root)

	if len(output) == 0 {
		errorStr := fmt.Sprintf("%s: Error printing json output ", wu.FuncName())
		return errors.New(errorStr)
	}

	w.jsonPayload = string(output)

	return nil
}

/*
 * Function to process a single field and build the equivalent
 * GoJSON object
 */
func (w *JsonCoder) processSingleField(parent *jsonez.GoJSON, input []byte,
	fDesc *wu.FieldDesc, fMap wu.FdMap, vLength, prefLength uint32,
	wType wu.ProtobufWireType) error {
	var child *jsonez.GoJSON
	var err error
	var temp []byte
	var fType int

	/*
	 * create a GoJSON object based on the field type and
	 * populate the value
	 */
	switch fDesc.Ftype {
	case wu.PROTOBUF_TYPE_ENUM:
		fallthrough
	case wu.PROTOBUF_TYPE_INT32:
		if wType != wu.PROTOBUF_WIRE_TYPE_VARINT {
			errorStr := fmt.Sprintf("%s: Wire type is not "+
				"wu.PROTOBUF_WIRE_TYPE_VARINT for for field %s", wu.FuncName(),
				fDesc.Fname)
			return errors.New(errorStr)
		}

		i32Val := wu.ParseInt32(vLength, input)

		if float64(i32Val) < 0 {
			fType = jsonez.JSON_INT
		} else {
			fType = jsonez.JSON_UINT
		}
		child = jsonez.AllocNumber(float64(i32Val), fType)

	case wu.PROTOBUF_TYPE_UINT32:
		if wType != wu.PROTOBUF_WIRE_TYPE_VARINT {
			errorStr := fmt.Sprintf("%s: Wire type is not "+
				"wu.PROTOBUF_WIRE_TYPE_VARINT for for field %s", wu.FuncName(),
				fDesc.Fname)
			return errors.New(errorStr)
		}

		u32Val := wu.ParseUint32(vLength, input)
		child = jsonez.AllocNumber(float64(u32Val), jsonez.JSON_UINT)

	case wu.PROTOBUF_TYPE_SINT32:
		if wType != wu.PROTOBUF_WIRE_TYPE_VARINT {
			errorStr := fmt.Sprintf("%s: Wire type is not "+
				"wu.PROTOBUF_WIRE_TYPE_VARINT for for field %s", wu.FuncName(),
				fDesc.Fname)
			return errors.New(errorStr)
		}

		i32Val := wu.Unzigzag32(wu.ParseUint32(vLength, input))

		if float64(i32Val) < 0 {
			fType = jsonez.JSON_INT
		} else {
			fType = jsonez.JSON_UINT
		}
		child = jsonez.AllocNumber(float64(i32Val), fType)

	case wu.PROTOBUF_TYPE_SFIXED32:
		fallthrough
	case wu.PROTOBUF_TYPE_FIXED32:
		fallthrough
	case wu.PROTOBUF_TYPE_FLOAT:
		if wType != wu.PROTOBUF_WIRE_TYPE_32BIT {
			errorStr := fmt.Sprintf("%s: Wire type is not "+
				"wu.PROTOBUF_WIRE_TYPE_32BIT for for field %s", wu.FuncName(),
				fDesc.Fname)
			return errors.New(errorStr)
		}

		u32Val := wu.ParseFixedUint32(input)

		var fType int
		if fDesc.Ftype == wu.PROTOBUF_TYPE_FLOAT {
			fType = jsonez.JSON_DOUBLE
			child = jsonez.AllocNumber(float64(math.Float32frombits(u32Val)),
				fType)
		} else if fDesc.Ftype == wu.PROTOBUF_TYPE_SFIXED32 {
			if int32(u32Val) < 0 {
				fType = jsonez.JSON_INT
			} else {
				fType = jsonez.JSON_UINT
			}

			child = jsonez.AllocNumber(float64(int32(u32Val)), fType)
		} else {
			fType = jsonez.JSON_UINT
			child = jsonez.AllocNumber(float64(u32Val), fType)
		}

	case wu.PROTOBUF_TYPE_INT64:
		fallthrough
	case wu.PROTOBUF_TYPE_UINT64:
		if wType != wu.PROTOBUF_WIRE_TYPE_VARINT {
			errorStr := fmt.Sprintf("%s: Wire type is not "+
				"wu.PROTOBUF_WIRE_TYPE_VARINT for for field %s", wu.FuncName(),
				fDesc.Fname)
			return errors.New(errorStr)
		}

		u64Val := wu.ParseUint64(vLength, input)

		if float64(u64Val) < 0 {
			fType = jsonez.JSON_INT
		} else {
			fType = jsonez.JSON_UINT
		}
		child = jsonez.AllocNumber(float64(u64Val), fType)

	case wu.PROTOBUF_TYPE_SINT64:
		if wType != wu.PROTOBUF_WIRE_TYPE_VARINT {
			errorStr := fmt.Sprintf("%s: Wire type is not "+
				"wu.PROTOBUF_WIRE_TYPE_VARINT for for field %s", wu.FuncName(),
				fDesc.Fname)
			return errors.New(errorStr)
		}

		i64Val := wu.Unzigzag64(wu.ParseUint64(vLength, input))

		if float64(i64Val) < 0 {
			fType = jsonez.JSON_INT
		} else {
			fType = jsonez.JSON_UINT
		}
		child = jsonez.AllocNumber(float64(i64Val), fType)

	case wu.PROTOBUF_TYPE_SFIXED64:
		fallthrough
	case wu.PROTOBUF_TYPE_FIXED64:
		fallthrough
	case wu.PROTOBUF_TYPE_DOUBLE:
		if wType != wu.PROTOBUF_WIRE_TYPE_64BIT {
			errorStr := fmt.Sprintf("%s: Wire type is not "+
				"wu.PROTOBUF_WIRE_TYPE_64BIT for for field %s", wu.FuncName(),
				fDesc.Fname)
			return errors.New(errorStr)
		}

		u64Val := wu.ParseFixedUint64(input)

		var fType int
		if fDesc.Ftype == wu.PROTOBUF_TYPE_DOUBLE {
			fType = jsonez.JSON_DOUBLE
			child = jsonez.AllocNumber(math.Float64frombits(u64Val), fType)
		} else if fDesc.Ftype == wu.PROTOBUF_TYPE_SFIXED64 {
			if int64(u64Val) < 0 {
				fType = jsonez.JSON_INT
			} else {
				fType = jsonez.JSON_UINT
			}
			child = jsonez.AllocNumber(float64(int64(u64Val)), fType)
		} else {
			fType = jsonez.JSON_UINT
			child = jsonez.AllocNumber(float64(u64Val), fType)
		}

	case wu.PROTOBUF_TYPE_BOOL:
		child = jsonez.AllocBool(wu.ParseBoolean(vLength, input))

	case wu.PROTOBUF_TYPE_BYTES:
		fallthrough
	case wu.PROTOBUF_TYPE_STRING:
		if wType != wu.PROTOBUF_WIRE_TYPE_LENGTH_PREFIXED {
			errorStr := fmt.Sprintf("%s: Wire type is not "+
				"wu.PROTOBUF_WIRE_TYPE_LENGTH_PREFIXED for for field %s",
				wu.FuncName(), fDesc.Fname)
			return errors.New(errorStr)
		}

		child = jsonez.AllocString(string(input[prefLength:vLength]))

	case wu.PROTOBUF_TYPE_MESSAGE:
		if wType != wu.PROTOBUF_WIRE_TYPE_LENGTH_PREFIXED {
			errorStr := fmt.Sprintf("%s: Wire type is not "+
				"wu.PROTOBUF_WIRE_TYPE_LENGTH_PREFIXED for field %s",
				wu.FuncName(), fDesc.Fname)
			return errors.New(errorStr)
		} else if fDesc.SubField == nil {
			errorStr := fmt.Sprintf("%s: Field map is empty for message "+
				"field %s", wu.FuncName(), fDesc.Fname)
			return errors.New(errorStr)
		}

		child = jsonez.AllocObject()
		child.Key = fDesc.Fname
	}

	/*
	 * Add the child GoJSON object to the parent
	 */
	if parent.Jsontype == jsonez.JSON_OBJECT {
		parent.AddEntryToObject(fDesc.Fname, child)
	} else if parent.Jsontype == jsonez.JSON_ARRAY {
		parent.AddEntryToArray(child)
	} else {
		errorStr := fmt.Sprintf("%s: Parent GoJSON object with key %s "+
			"is not an object or an array", wu.FuncName(), parent.Key)
		return errors.New(errorStr)
	}

	/*
	 * If the current parsed value is a message, process the
	 * child entries
	 */
	if child.Jsontype == jsonez.JSON_OBJECT {
		temp = append(temp, input[prefLength:vLength]...)
		err = w.parseProtoInput(child, temp, fDesc.SubField)
	}

	if err != nil {
		return err
	}

	return nil
}

/*
 * Function to process a repeated field and build the GoJSON array
 */
func (w *JsonCoder) ProcessRepField(parent *jsonez.GoJSON, input []byte,
	fDesc *wu.FieldDesc, fMap wu.FdMap, vLength, prefLength uint32,
	wType wu.ProtobufWireType) error {

	var arrObj *jsonez.GoJSON
	var err error

	/*
	 * Create the JSON array object if it doesn't exist
	 */
	arrObj, err = parent.Get(fDesc.Fname)
	if err != nil {
		arrObj = jsonez.AllocArray()
		arrObj.Key = fDesc.Fname
		/*
		 * Add the child GoJSON object to the parent
		 */
		if parent.Jsontype == jsonez.JSON_OBJECT {
			parent.AddEntryToObject(fDesc.Fname, arrObj)
		} else if parent.Jsontype == jsonez.JSON_ARRAY {
			parent.AddEntryToArray(arrObj)
		} else {
			errorStr := fmt.Sprintf("%s: Parent GoJSON object with key %s "+
				"is not an object or an array", wu.FuncName(), parent.Key)
			return errors.New(errorStr)
		}
	}

	return w.processSingleField(arrObj, input, fDesc, fMap, vLength,
		prefLength, wType)
}

/*
 * Function to process packed repeated field
 */
func (w *JsonCoder) ProcessPackedRepField(parent *jsonez.GoJSON, input []byte,
	fDesc *wu.FieldDesc, fMap wu.FdMap, vLength, prefLength uint32,
	wType wu.ProtobufWireType) error {

	var child *jsonez.GoJSON
	var count, i, shift uint32
	var fType int

	switch fDesc.Ftype {
	case wu.PROTOBUF_TYPE_SFIXED32:
		fallthrough
	case wu.PROTOBUF_TYPE_FIXED32:
		fallthrough
	case wu.PROTOBUF_TYPE_FLOAT:
		count = (vLength - prefLength) / 4
		input = input[prefLength:]

		for i = 0; i < count; i++ {
			u32Val := wu.ParseFixedUint32(input)

			if fDesc.Ftype == wu.PROTOBUF_TYPE_FLOAT {
				fType = jsonez.JSON_DOUBLE
			} else {
				if float64(u32Val) < 0 {
					fType = jsonez.JSON_INT
				} else {
					fType = jsonez.JSON_UINT
				}
			}
			child = jsonez.AllocNumber(float64(u32Val), fType)

			/*
			 * Add the child GoJSON object to the parent
			 */
			if parent.Jsontype == jsonez.JSON_ARRAY {
				parent.AddEntryToArray(child)
			} else {
				errorStr := fmt.Sprintf("%s: Parent GoJSON object with "+
					"key %s is not an object or an array", wu.FuncName(),
					parent.Key)
				return errors.New(errorStr)
			}

			input = input[4:]
		}

	case wu.PROTOBUF_TYPE_SFIXED64:
		fallthrough
	case wu.PROTOBUF_TYPE_FIXED64:
		fallthrough
	case wu.PROTOBUF_TYPE_DOUBLE:
		count = (vLength - prefLength) / 8
		input = input[prefLength:]

		for i = 0; i < count; i++ {
			u64Val := wu.ParseFixedUint64(input)

			if fDesc.Ftype == wu.PROTOBUF_TYPE_DOUBLE {
				fType = jsonez.JSON_DOUBLE
			} else {
				if float64(u64Val) < 0 {
					fType = jsonez.JSON_INT
				} else {
					fType = jsonez.JSON_UINT
				}
			}
			child = jsonez.AllocNumber(float64(u64Val), fType)

			/*
			 * Add the child GoJSON object to the parent
			 */
			if parent.Jsontype == jsonez.JSON_ARRAY {
				parent.AddEntryToArray(child)
			} else {
				errorStr := fmt.Sprintf("%s: Parent GoJSON object with "+
					"key %s is not an object or an array", wu.FuncName(),
					parent.Key)
				return errors.New(errorStr)
			}

			input = input[8:]
		}

	case wu.PROTOBUF_TYPE_ENUM:
		fallthrough
	case wu.PROTOBUF_TYPE_INT32:
		input = input[prefLength:]
		for {
			length := uint32(len(input))
			if length > 0 {
				shift = wu.ScanVarint(length, input)
				if shift == 0 {
					errorStr := fmt.Sprintf("%s Bad packed-repeated enum/int32 "+
						"value for field %s", wu.FuncName(), fDesc.Fname)
					return errors.New(errorStr)
				}

				i32Val := wu.ParseInt32(shift, input)
				if float64(i32Val) < 0 {
					fType = jsonez.JSON_INT
				} else {
					fType = jsonez.JSON_UINT
				}
				child = jsonez.AllocNumber(float64(i32Val), fType)

				/*
				 * Add the child GoJSON object to the parent
				 */
				if parent.Jsontype == jsonez.JSON_ARRAY {
					parent.AddEntryToArray(child)
				} else {
					errorStr := fmt.Sprintf("%s: Parent GoJSON object with "+
						"key %s is not an object or an array", wu.FuncName(),
						parent.Key)
					return errors.New(errorStr)
				}

				input = input[shift:]
			} else {
				break
			}
		}

	case wu.PROTOBUF_TYPE_SINT32:
		input = input[prefLength:]
		for {
			length := uint32(len(input))
			if length > 0 {
				shift = wu.ScanVarint(length, input)
				if shift == 0 {
					errorStr := fmt.Sprintf("%s Bad packed-repeated sint32 "+
						"value for field %s", wu.FuncName(), fDesc.Fname)
					return errors.New(errorStr)
				}

				i32Val := wu.Unzigzag32(wu.ParseUint32(shift, input))

				if float64(i32Val) < 0 {
					fType = jsonez.JSON_INT
				} else {
					fType = jsonez.JSON_UINT
				}

				child = jsonez.AllocNumber(float64(i32Val), fType)

				/*
				 * Add the child GoJSON object to the parent
				 */
				if parent.Jsontype == jsonez.JSON_ARRAY {
					parent.AddEntryToArray(child)
				} else {
					errorStr := fmt.Sprintf("%s: Parent GoJSON object with "+
						"key %s is not an object or an array", wu.FuncName(),
						parent.Key)
					return errors.New(errorStr)
				}

				input = input[shift:]
			} else {
				break
			}
		}

	case wu.PROTOBUF_TYPE_UINT32:
		input = input[prefLength:]
		for {
			length := uint32(len(input))
			if length > 0 {
				shift = wu.ScanVarint(length, input)
				if shift == 0 {
					errorStr := fmt.Sprintf("%s Bad packed-repeated uint32 "+
						"value for field %s", wu.FuncName(), fDesc.Fname)
					return errors.New(errorStr)
				}

				u32Val := wu.ParseUint32(shift, input)

				if float64(u32Val) < 0 {
					fType = jsonez.JSON_INT
				} else {
					fType = jsonez.JSON_UINT
				}
				child = jsonez.AllocNumber(float64(u32Val), fType)

				/*
				 * Add the child GoJSON object to the parent
				 */
				if parent.Jsontype == jsonez.JSON_ARRAY {
					parent.AddEntryToArray(child)
				} else {
					errorStr := fmt.Sprintf("%s: Parent GoJSON object with "+
						"key %s is not an object or an array", wu.FuncName(),
						parent.Key)
					return errors.New(errorStr)
				}

				input = input[shift:]
			} else {
				break
			}
		}

	case wu.PROTOBUF_TYPE_SINT64:
		input = input[prefLength:]
		for {
			length := uint32(len(input))
			if length > 0 {
				shift = wu.ScanVarint(length, input)
				if shift == 0 {
					errorStr := fmt.Sprintf("%s Bad packed-repeated sint64 "+
						"value for field %s", wu.FuncName(), fDesc.Fname)
					return errors.New(errorStr)
				}

				i64Val := wu.Unzigzag64(wu.ParseUint64(shift, input))

				if float64(i64Val) < 0 {
					fType = jsonez.JSON_INT
				} else {
					fType = jsonez.JSON_UINT
				}
				child = jsonez.AllocNumber(float64(i64Val), fType)

				/*
				 * Add the child GoJSON object to the parent
				 */
				if parent.Jsontype == jsonez.JSON_ARRAY {
					parent.AddEntryToArray(child)
				} else {
					errorStr := fmt.Sprintf("%s: Parent GoJSON object with "+
						"key %s is not an object or an array", wu.FuncName(),
						parent.Key)
					return errors.New(errorStr)
				}

				input = input[shift:]
			} else {
				break
			}
		}

	case wu.PROTOBUF_TYPE_INT64:
		fallthrough
	case wu.PROTOBUF_TYPE_UINT64:
		input = input[prefLength:]
		for {
			length := uint32(len(input))
			if length > 0 {
				shift = wu.ScanVarint(length, input)
				if shift == 0 {
					errorStr := fmt.Sprintf("%s Bad packed-repeated "+
						"int64/uint64 value for field %s", wu.FuncName(),
						fDesc.Fname)
					return errors.New(errorStr)
				}

				u64Val := wu.ParseUint64(shift, input)

				if float64(u64Val) < 0 {
					fType = jsonez.JSON_INT
				} else {
					fType = jsonez.JSON_UINT
				}
				child = jsonez.AllocNumber(float64(u64Val), fType)

				/*
				 * Add the child GoJSON object to the parent
				 */
				if parent.Jsontype == jsonez.JSON_ARRAY {
					parent.AddEntryToArray(child)
				} else {
					errorStr := fmt.Sprintf("%s: Parent GoJSON object with "+
						"key %s is not an object or an array", wu.FuncName(),
						parent.Key)
					return errors.New(errorStr)
				}

				input = input[shift:]
			} else {
				break
			}
		}

	case wu.PROTOBUF_TYPE_BOOL:
		input = input[prefLength:]
		length := uint32(len(input))

		for i = 0; i < length; i++ {
			if input[i] > 1 {
				errorStr := fmt.Sprintf("%s Bad packed-repeated "+
					"bool value for field %s", wu.FuncName(),
					fDesc.Fname)
				return errors.New(errorStr)
			} else if input[i] == 1 {
				child = jsonez.AllocBool(true)
			} else {
				child = jsonez.AllocBool(false)
			}

			/*
			 * Add the child GoJSON object to the parent
			 */
			if parent.Jsontype == jsonez.JSON_ARRAY {
				parent.AddEntryToArray(child)
			} else {
				errorStr := fmt.Sprintf("%s: Parent GoJSON object with "+
					"key %s is not an object or an array", wu.FuncName(),
					parent.Key)
				return errors.New(errorStr)
			}
		}

	default:
		errorStr := fmt.Sprintf("%s: Incompatible packed-repeated type"+
			"%d for field %s", wu.FuncName(), fDesc.Ftype, parent.Key)
		return errors.New(errorStr)
	}

	return nil
}

/*
 * Function to process the current proto input
 */
func (w *JsonCoder) processProtoInput(cur *jsonez.GoJSON, input []byte,
	fDesc *wu.FieldDesc, fMap wu.FdMap, vLength, prefLength uint32,
	wType wu.ProtobufWireType) error {
	var arrObj *jsonez.GoJSON
	var err error

	switch fDesc.Flabel {
	case wu.PROTOBUF_LABEL_REQUIRED:
		fallthrough
	case wu.PROTOBUF_LABEL_OPTIONAL:
		return w.processSingleField(cur, input, fDesc, fMap, vLength,
			prefLength, wType)

	case wu.PROTOBUF_LABEL_REPEATED:
		/*
		 * Create the JSON array object if it doesn't exist
		 */
		arrObj, err = cur.Get(fDesc.Fname)
		if err != nil {
			arrObj = jsonez.AllocArray()
			arrObj.Key = fDesc.Fname
			/*
			 * Add the child GoJSON object to the parent
			 */
			if cur.Jsontype == jsonez.JSON_OBJECT {
				cur.AddEntryToObject(fDesc.Fname, arrObj)
			} else if cur.Jsontype == jsonez.JSON_ARRAY {
				cur.AddEntryToArray(arrObj)
			} else {
				errorStr := fmt.Sprintf("%s: Parent GoJSON object with key %s "+
					"is not an object or an array", wu.FuncName(), cur.Key)
				return errors.New(errorStr)
			}
		}

		if (fDesc.Flags & wu.PROTOBUF_FIELD_FLAG_PACKED) != 0 {
			return w.ProcessPackedRepField(arrObj, input, fDesc, fMap, vLength,
				prefLength, wType)

		} else {
			return w.processSingleField(arrObj, input, fDesc, fMap, vLength,
				prefLength, wType)
		}

	default:
		errorStr := fmt.Sprintf("%s: Unrecognized label type for child "+
			"tag %s of parent %s ", wu.FuncName(), fDesc.Fid, cur.Key)
		return errors.New(errorStr)
	}
}

/*
 * Function to process the current wiretype
 */
func (w *JsonCoder) processWireType(length uint32, buf []byte) (uint32,
	uint32, uint32, uint32, wu.ProtobufWireType, error) {
	var used, vLength, prefLength, tag, max, i uint32
	var fDesc *wu.FieldDesc
	var wType wu.ProtobufWireType

	used = wu.ParseTagAndWiretype(length, buf, &tag, &wType)

	if used == 0 {
		errorStr := fmt.Sprintf("%s: Error parsing tag/wiretype",
			wu.FuncName())

		return 0, 0, 0, 0, 0, errors.New(errorStr)
	}

	input := buf[used:]

	switch wType {
	case wu.PROTOBUF_WIRE_TYPE_VARINT:
		if length < 10 {
			max = length
		} else {
			max = 10
		}

		for i = 0; i < max; i++ {
			if (input[i] & 0x80) == 0 {
				break
			}
		}

		if i == max {
			errorStr := fmt.Sprintf("%s: Unterminated varint "+
				"found for field %s", wu.FuncName(), fDesc.Fname)

			return 0, 0, 0, 0, 0, errors.New(errorStr)
		}

		vLength = uint32(i + 1)

	case wu.PROTOBUF_WIRE_TYPE_32BIT:
		if length < 4 {
			errorStr := fmt.Sprintf("%s Too short for fixed 32 bit wiretype"+
				" when parsing field %s", wu.FuncName(), fDesc.Fname)

			return 0, 0, 0, 0, 0, errors.New(errorStr)
		}
		vLength = 4

	case wu.PROTOBUF_WIRE_TYPE_64BIT:
		if length < 4 {
			errorStr := fmt.Sprintf("%s Too short for fixed 64 bit wiretype"+
				" when parsing field %s", wu.FuncName(), fDesc.Fname)

			return 0, 0, 0, 0, 0, errors.New(errorStr)
		}
		vLength = 8

	case wu.PROTOBUF_WIRE_TYPE_LENGTH_PREFIXED:
		vLength = wu.ScanLengthPrefixData(length, input, &prefLength)

	default:
		errorStr := fmt.Sprintf("%s Unknown wire type obtained"+
			" when parsing field %s", wu.FuncName(), fDesc.Fname)

		return 0, 0, 0, 0, 0, errors.New(errorStr)

	}

	return used, vLength, prefLength, tag, wType, nil
}

/*
 * Function to parse the proto input and build a JSON
 * output
 */
func (w *JsonCoder) parseProtoInput(cur *jsonez.GoJSON, input []byte,
	fMap wu.FdMap) error {
	var ok bool
	var fDesc *wu.FieldDesc
	var used, tag, vLength, prefLength uint32
	var wType wu.ProtobufWireType
	var err error

	for {
		length := uint32(len(input))
		if length > 0 {
			used, vLength, prefLength, tag, wType, err =
				w.processWireType(length, input)

			if err != nil {
				return err
			}

			input = input[used:]

			/*
			 * Get the field descriptor corresponding to the tag
			 * received
			 */
			if fDesc, ok = fMap[strconv.Itoa(int(tag))]; !ok {
				log.Println(wu.FuncName(), ":Field id", tag, "not found")
				input = input[vLength:]
				continue
			}

			err = w.processProtoInput(cur, input, fDesc, fMap, vLength,
				prefLength, wType)

			if err != nil {
				return err
			}

			input = input[vLength:]
		} else {
			break
		}
	}

	return nil
}
