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

	wu "git.juniper.net/sksubra/wedge/util"
	"github.com/srikanth2212/jsonez"
)

/*
 * Fucntion to encode JSON to protobuf
 */
func JsonToProto(w *JsonCoder) ([]byte, error) {

	var input, output []byte

	input = []byte(w.jsonPayload)

	/*
	 * To Accomodate for empty messages
	 */
	if len(input) > 0 {

		root, err := jsonez.GoJSONParse(input)

		if err != nil {
			return []byte{}, err
		}

		output, _, err = w.parseJSONInput(root, output, w.inDesc.FieldDescMap,
			true)
		if err != nil {
			return []byte{}, err
		}
	}

	return output, nil
}

/*
 * Pack a prefixed field i.e a message
 */
func (w *JsonCoder) packPrefixedField(cur *jsonez.GoJSON, buf []byte,
	fMap wu.FdMap) ([]byte, uint32, error) {
	var err error
	var output []byte
	var rv, temp uint32

	output, temp, err = w.parseJSONInput(cur, output, fMap, true)
	if err != nil {
		return []byte{}, 0, err
	}

	buf, rv = wu.Uint32Pack(temp, buf)

	return append(buf, output...), rv + temp, nil

}

/*
 * Function to pack a single field into protobuf.
 * Returns the appended buf, number of uint8 chars added an error if any
 */
func (w *JsonCoder) packSingleField(cur *jsonez.GoJSON, buf []byte,
	fDesc *wu.FieldDesc, fid uint32, fMap wu.FdMap, isNotPacked bool) ([]byte,
	uint32, error) {
	var rv, temp uint32
	var output []byte

	if isNotPacked {
		output, rv = wu.TagPack(fid, output)
	}

	/*
	 * Pack the field based on type.
	 */
	switch fDesc.Ftype {
	case wu.PROTOBUF_TYPE_SINT32:
		if isNotPacked {
			output[0] |= wu.PROTOBUF_WIRE_TYPE_VARINT
		}

		if cur.Jsontype == jsonez.JSON_INT {
			output, temp = wu.Sint32Pack(int32(cur.Valint), output)
		} else if cur.Jsontype == jsonez.JSON_UINT {
			output, temp = wu.Sint32Pack(int32(cur.Valuint), output)
		} else {
			s := fmt.Sprintf("Unexpected json type for field %s", cur.Key)
			return buf, 0, errors.New(s)
		}

		return append(buf, output...), rv + temp, nil

	case wu.PROTOBUF_TYPE_INT32:
		if isNotPacked {
			output[0] |= wu.PROTOBUF_WIRE_TYPE_VARINT
		}
		if cur.Jsontype == jsonez.JSON_INT {
			output, temp = wu.Int32Pack(int32(cur.Valint), output)
		} else if cur.Jsontype == jsonez.JSON_UINT {
			output, temp = wu.Int32Pack(int32(cur.Valuint), output)
		} else {
			s := fmt.Sprintf("Unexpected json type for field %s", cur.Key)
			return buf, 0, errors.New(s)
		}

		return append(buf, output...), rv + temp, nil

	case wu.PROTOBUF_TYPE_UINT32:
		fallthrough
	case wu.PROTOBUF_TYPE_ENUM:
		if isNotPacked {
			output[0] |= wu.PROTOBUF_WIRE_TYPE_VARINT
		}

		if cur.Jsontype == jsonez.JSON_INT {
			output, temp = wu.Uint32Pack(uint32(cur.Valint), output)
		} else if cur.Jsontype == jsonez.JSON_UINT {
			output, temp = wu.Uint32Pack(uint32(cur.Valuint), output)
		} else {
			s := fmt.Sprintf("Unexpected json type for field %s", cur.Key)
			return buf, 0, errors.New(s)
		}

		return append(buf, output...), rv + temp, nil

	case wu.PROTOBUF_TYPE_SINT64:
		if isNotPacked {
			output[0] |= wu.PROTOBUF_WIRE_TYPE_VARINT
		}

		if cur.Jsontype == jsonez.JSON_INT {
			output, temp = wu.Sint64Pack(int64(cur.Valint), output)
		} else if cur.Jsontype == jsonez.JSON_UINT {
			output, temp = wu.Sint64Pack(int64(cur.Valuint), output)
		} else {
			s := fmt.Sprintf("Unexpected json type for field %s", cur.Key)
			return buf, 0, errors.New(s)
		}

		return append(buf, output...), rv + temp, nil

	case wu.PROTOBUF_TYPE_INT64:
		fallthrough
	case wu.PROTOBUF_TYPE_UINT64:
		if isNotPacked {
			output[0] |= wu.PROTOBUF_WIRE_TYPE_VARINT
		}

		if cur.Jsontype == jsonez.JSON_INT {
			output, temp = wu.Uint64Pack(uint64(cur.Valint), output)
		} else if cur.Jsontype == jsonez.JSON_UINT {
			output, temp = wu.Uint64Pack(uint64(cur.Valuint), output)
		} else {
			s := fmt.Sprintf("Unexpected json type for field %s", cur.Key)
			return buf, 0, errors.New(s)
		}

		return append(buf, output...), rv + temp, nil

	case wu.PROTOBUF_TYPE_SFIXED32:
		fallthrough
	case wu.PROTOBUF_TYPE_FIXED32:
		if isNotPacked {
			output[0] |= wu.PROTOBUF_WIRE_TYPE_32BIT
		}

		if cur.Jsontype == jsonez.JSON_INT {
			output, temp = wu.Fixed32Pack(uint32(cur.Valint), output)
		} else if cur.Jsontype == jsonez.JSON_UINT {
			output, temp = wu.Fixed32Pack(uint32(cur.Valuint), output)
		} else {
			s := fmt.Sprintf("Unexpected json type for field %s", cur.Key)
			return buf, 0, errors.New(s)
		}

		return append(buf, output...), rv + temp, nil

	case wu.PROTOBUF_TYPE_FLOAT:
		if isNotPacked {
			output[0] |= wu.PROTOBUF_WIRE_TYPE_32BIT
		}
		output, temp = wu.Fixed32Pack(math.Float32bits(float32(cur.Valdouble)),
			output)
		return append(buf, output...), rv + temp, nil

	case wu.PROTOBUF_TYPE_SFIXED64:
		fallthrough
	case wu.PROTOBUF_TYPE_FIXED64:
		if isNotPacked {
			output[0] |= wu.PROTOBUF_WIRE_TYPE_64BIT
		}

		if cur.Jsontype == jsonez.JSON_INT {
			output, temp = wu.Fixed64Pack(uint64(cur.Valint), output)
		} else if cur.Jsontype == jsonez.JSON_UINT {
			output, temp = wu.Fixed64Pack(uint64(cur.Valuint), output)
		} else {
			s := fmt.Sprintf("Unexpected json type for field %s", cur.Key)
			return buf, 0, errors.New(s)
		}

		return append(buf, output...), rv + temp, nil

	case wu.PROTOBUF_TYPE_DOUBLE:
		if isNotPacked {
			output[0] |= wu.PROTOBUF_WIRE_TYPE_64BIT
		}
		output, temp = wu.Fixed64Pack(math.Float64bits(cur.Valdouble), output)
		return append(buf, output...), rv + temp, nil

	case wu.PROTOBUF_TYPE_BOOL:
		if isNotPacked {
			output[0] |= wu.PROTOBUF_WIRE_TYPE_VARINT
		}
		output, temp := wu.BooleanPack(cur.Valbool, output)
		return append(buf, output...), rv + temp, nil

	case wu.PROTOBUF_TYPE_STRING:
		if isNotPacked {
			output[0] |= wu.PROTOBUF_WIRE_TYPE_LENGTH_PREFIXED
		}
		output, temp = wu.StringPack(cur.Valstr, output)
		return append(buf, output...), rv + temp, nil

	case wu.PROTOBUF_TYPE_BYTES:
		if isNotPacked {
			output[0] |= wu.PROTOBUF_WIRE_TYPE_LENGTH_PREFIXED
		}
		output, temp = wu.BytesDataPack([]byte(cur.Valstr), output)
		return append(buf, output...), rv + temp, nil

	case wu.PROTOBUF_TYPE_MESSAGE:
		subMap := fDesc.SubField
		if isNotPacked {
			output[0] |= wu.PROTOBUF_WIRE_TYPE_LENGTH_PREFIXED
		}

		output, temp, err := w.packPrefixedField(cur, output, subMap)
		if err != nil {
			return []byte{}, 0, err
		}

		return append(buf, output...), rv + temp, nil
	default:
		s := fmt.Sprintf("Unknown field type %d for field of type %d",
			fDesc.Ftype, fid)
		return buf, 0, errors.New(s)
	}

	return buf, 0, nil
}

/*
 * Pack a repeated field. Handle cases for both packed and unpacked
 * fields
 */
func (w *JsonCoder) packRepeatedField(cur *jsonez.GoJSON, buf []byte,
	fDesc *wu.FieldDesc, fMap wu.FdMap) ([]byte, uint32, error) {
	var child *jsonez.GoJSON
	var err error
	var subMap wu.FdMap
	var isPacked bool
	var output1, output2 []byte
	var rv, lLen, temp, fid uint32

	f, err := strconv.ParseInt(fDesc.Fid, 10, 32)
	fid = uint32(f)

	if err != nil {
		return []byte{}, 0, err
	}

	if fDesc.Ftype != wu.PROTOBUF_TYPE_MESSAGE &&
		(fDesc.Flags&wu.PROTOBUF_FIELD_FLAG_PACKED) != 0 {
		output1, rv = wu.TagPack(uint32(fid), output1)
		output1[0] |= byte(wu.PROTOBUF_WIRE_TYPE_LENGTH_PREFIXED)
		isPacked = true
	} else {
		isPacked = false
	}

	if fDesc.Ftype == wu.PROTOBUF_TYPE_MESSAGE {
		subMap = fDesc.SubField
		if subMap == nil {
			s := fmt.Sprintf("%s Sub filed map not found for field %s"+
				"of type wu.PROTOBUF_TYPE_MESSAGE", wu.FuncName(), fDesc.Fname)
			return []byte{}, 0, errors.New(s)
		}
	} else {
		subMap = fMap
	}

	for i := 0; i < cur.GetArraySize(); i++ {
		child, err = cur.GetArrayElemByIndex(i)
		if err != nil {
			return []byte{}, 0, err
		}

		output2, temp, err = w.packSingleField(child, output2, fDesc, fid,
			subMap, !isPacked)
		lLen += temp
	}

	if (fDesc.Flags & wu.PROTOBUF_FIELD_FLAG_PACKED) != 0 {
		/*
		 * Pack the list length
		 */
		output1, temp = wu.Uint32Pack(uint32(lLen), output1)

	}

	output1 = append(output1, output2...)
	return append(buf, output1...), rv + lLen, nil
}

/*
 * Function to parse JSON input and build proto
 */
func (w *JsonCoder) parseJSONInput(cur *jsonez.GoJSON, buf []byte, fMap wu.FdMap,
	isNotPacked bool) ([]byte, uint32, error) {
	var child *jsonez.GoJSON
	var fDesc *wu.FieldDesc
	var fid uint32
	var ok bool
	var rv, temp uint32

	/*
	 * Walk the children i.e fields of this GoJSON object and build the
	 * equivalent proto message
	 */

	child = cur.Child
	for {
		if child != nil {
			/*
			 * Get the field descriptor for this field
			 */
			if fDesc, ok = fMap[child.Key]; !ok {
				errorStr := fmt.Sprintf("%s: Field %s not found for message %s",
					wu.FuncName(), child.Key, cur.Key)
				return []byte{}, 0, errors.New(errorStr)
			}

			/*
			 * If the child is an array, pack as a repeated field
			 */
			v, err := strconv.ParseInt(fDesc.Fid, 10, 32)
			if err != nil {
				return []byte{}, 0, err
			}
			fid = uint32(v)

			if fDesc.Flabel == wu.PROTOBUF_LABEL_REPEATED {
				buf, temp, err = w.packRepeatedField(child, buf, fDesc, fMap)
			} else {
				buf, temp, err = w.packSingleField(child, buf, fDesc, fid,
					fMap, isNotPacked)

				if err != nil {
					return []byte{}, 0, err
				}
			}

			rv += temp
			child = child.Next
		} else {
			break
		}
	}

	return buf, rv, nil
}
