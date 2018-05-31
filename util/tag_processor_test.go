/*
 * Copyright (c) 2018, Juniper Networks, Inc.
 * All rights reserved.
 */

package util

import (
	"fmt"
	"math"
	"testing"
)

func TestInt32(t *testing.T) {
	var x int32 = -1984
	bytes := make([]byte, 0)
	var wtype ProtobufWireType
	var tag uint32

	t.Logf("%s: val is: %d", FuncName(), x)
	bytes, _ = TagPack(35, bytes)
	bytes[0] = bytes[0] | byte(PROTOBUF_WIRE_TYPE_VARINT)

	bytes, l := Int32Pack(x, bytes)

	t.Log(FuncName()+":", "Packed length:", len(bytes), ", packed output:",
		bytes)

	used := ParseTagAndWiretype(l, bytes, &tag, &wtype)
	dbytes := bytes[used:]

	t.Log(FuncName(), ":Protobuf data obtained : ", dbytes)

	var i uint

	for i = 0; i < 10; i++ {
		if (dbytes[i] & 0x80) == 0 {
			break
		}
	}

	l = uint32(i + 1)

	val := ParseInt32(l, dbytes)

	if tag != 35 {
		t.Errorf("%s: Mismatch in tag value, Expected: 35, Got: %d",
			FuncName(), tag)
	}

	if int32(val) != x {
		t.Errorf("%s: Mismatch in output value, Expected: %d, Got: %d",
			FuncName(), x, int32(val))
	}
}

func TestUInt32(t *testing.T) {
	var x uint32 = 2018
	bytes := make([]byte, 0)
	var wtype ProtobufWireType
	var tag uint32

	t.Logf("%s: val is: %d", FuncName(), x)
	bytes, _ = TagPack(35, bytes)
	bytes[0] = bytes[0] | byte(PROTOBUF_WIRE_TYPE_VARINT)

	bytes, l := Uint32Pack(x, bytes)

	t.Log(FuncName()+":", "Packed length:", len(bytes), ", packed output:",
		bytes)

	used := ParseTagAndWiretype(l, bytes, &tag, &wtype)
	dbytes := bytes[used:]

	t.Log(FuncName(), ":Protobuf data obtained : ", dbytes)

	var i uint

	for i = 0; i < 10; i++ {
		if (dbytes[i] & 0x80) == 0 {
			break
		}
	}

	l = uint32(i + 1)

	val := ParseUint32(l, dbytes)

	if tag != 35 {
		t.Errorf("%s: Mismatch in tag value, Expected: 35, Got: %d",
			FuncName(), tag)
	}

	if val != x {
		t.Errorf("%s: Mismatch in output value, Expected: %d, Got: %d",
			FuncName(), x, val)
	}
}

func TestFixed32(t *testing.T) {
	var x uint32 = 2018
	bytes := make([]byte, 0)
	var wtype ProtobufWireType
	var tag uint32

	t.Logf("%s: val is: %d", FuncName(), x)
	bytes, _ = TagPack(35, bytes)
	bytes[0] |= PROTOBUF_WIRE_TYPE_32BIT

	bytes, l := Fixed32Pack(x, bytes)

	t.Log(FuncName()+":", "Packed length:", len(bytes), ", packed output:",
		bytes)

	used := ParseTagAndWiretype(l, bytes, &tag, &wtype)
	dbytes := bytes[used:]

	t.Log(FuncName(), ":Protobuf data obtained : ", dbytes)

	val := ParseFixedUint32(dbytes)

	if tag != 35 {
		t.Errorf("%s: Mismatch in tag value, Expected: 35, Got: %d",
			FuncName(), tag)
	}

	if val != x {
		t.Errorf("%s: Mismatch in output value, Expected: %d, Got: %d",
			FuncName(), x, val)
	}
}

func TestSint32(t *testing.T) {
	var x int32 = -123123
	bytes := make([]byte, 0)
	var wtype ProtobufWireType
	var tag uint32

	t.Logf("%s: val is: %d", FuncName(), x)
	bytes, _ = TagPack(35, bytes)
	bytes[0] = bytes[0] | byte(PROTOBUF_WIRE_TYPE_VARINT)

	bytes, l := Sint32Pack(x, bytes)

	t.Log(FuncName()+":", "Packed length:", len(bytes), ", packed output:",
		bytes)

	used := ParseTagAndWiretype(l, bytes, &tag, &wtype)
	dbytes := bytes[used:]

	t.Log(FuncName(), ":Protobuf data obtained : ", dbytes)

	var i uint

	for i = 0; i < 10; i++ {
		if (dbytes[i] & 0x80) == 0 {
			break
		}
	}

	l = uint32(i + 1)
	val := Unzigzag32(ParseUint32(l, dbytes))

	if int32(val) != x {
		t.Errorf("%s: Mismatch in output value, Expected: %d, Got: %d",
			FuncName(), x, int32(val))
	}
}

func TestSFixed32(t *testing.T) {
	var x int32 = -1234
	bytes := make([]byte, 0)
	var wtype ProtobufWireType
	var tag uint32

	t.Logf("%s: val is: %d", FuncName(), x)
	bytes, _ = TagPack(35, bytes)
	bytes[0] |= PROTOBUF_WIRE_TYPE_32BIT

	bytes, l := Fixed32Pack(uint32(x), bytes)

	t.Log(FuncName()+":", "Packed length:", len(bytes), ", packed output:",
		bytes)

	used := ParseTagAndWiretype(l, bytes, &tag, &wtype)
	dbytes := bytes[used:]

	t.Log(FuncName(), ":Protobuf data obtained : ", dbytes)

	val := ParseFixedUint32(dbytes)

	if tag != 35 {
		t.Errorf("%s: Mismatch in tag value, Expected: 35, Got: %d",
			FuncName(), tag)
	}

	if int32(val) != int32(x) {
		fmt.Println("uint_val is", val)
		fmt.Println("float_val is", float32(int32(val)))
		t.Errorf("%s: Mismatch in output value, Expected: %d, Got: %d",
			FuncName(), x, int32(val))
	}
}

func TestInt64(t *testing.T) {
	var x int64 = -1984123
	bytes := make([]byte, 0)
	var wtype ProtobufWireType
	var tag uint32

	t.Logf("%s: val is: %d", FuncName(), x)
	bytes, _ = TagPack(35, bytes)
	bytes[0] = bytes[0] | byte(PROTOBUF_WIRE_TYPE_VARINT)

	bytes, l := Uint64Pack(uint64(x), bytes)

	t.Log(FuncName()+":", "Packed length:", len(bytes), ", packed output:",
		bytes)

	used := ParseTagAndWiretype(l, bytes, &tag, &wtype)
	dbytes := bytes[used:]

	t.Log(FuncName(), ":Protobuf data obtained : ", dbytes)

	var i uint

	for i = 0; i < 10; i++ {
		if (dbytes[i] & 0x80) == 0 {
			break
		}
	}

	l = uint32(i + 1)
	val := ParseUint64(l, dbytes)

	if tag != 35 {
		t.Errorf("%s: Mismatch in tag value, Expected: 35, Got: %d",
			FuncName(), tag)
	}

	if int64(val) != x {
		t.Errorf("%s: Mismatch in output value, Expected: %d, Got: %d",
			FuncName(), x, int64(val))
	}
}

func TestUInt64(t *testing.T) {
	var x uint64 = 987654321
	bytes := make([]byte, 0)
	var wtype ProtobufWireType
	var tag uint32

	t.Logf("%s: val is: %d", FuncName(), x)
	bytes, _ = TagPack(35, bytes)
	bytes[0] = bytes[0] | byte(PROTOBUF_WIRE_TYPE_VARINT)

	bytes, l := Uint64Pack(x, bytes)

	t.Log(FuncName()+":", "Packed length:", len(bytes), ", packed output:",
		bytes)

	used := ParseTagAndWiretype(l, bytes, &tag, &wtype)
	dbytes := bytes[used:]

	t.Log(FuncName(), ":Protobuf data obtained : ", dbytes)

	var i uint

	for i = 0; i < 10; i++ {
		if (dbytes[i] & 0x80) == 0 {
			break
		}
	}

	l = uint32(i + 1)

	val := ParseUint64(l, dbytes)

	if tag != 35 {
		t.Errorf("%s: Mismatch in tag value, Expected: 35, Got: %d",
			FuncName(), tag)
	}

	if val != x {
		t.Errorf("%s: Mismatch in output value, Expected: %d, Got: %d",
			FuncName(), x, val)
	}
}

func TestFixed64(t *testing.T) {
	var x uint64 = 123456789
	bytes := make([]byte, 0)
	var wtype ProtobufWireType
	var tag uint32

	t.Logf("%s: val is: %d", FuncName(), x)
	bytes, _ = TagPack(35, bytes)
	bytes[0] |= PROTOBUF_WIRE_TYPE_32BIT

	bytes, l := Fixed64Pack(x, bytes)

	t.Log(FuncName()+":", "Packed length:", len(bytes), ", packed output:",
		bytes)

	used := ParseTagAndWiretype(l, bytes, &tag, &wtype)
	dbytes := bytes[used:]

	t.Log(FuncName(), ":Protobuf data obtained : ", dbytes)

	val := ParseFixedUint64(dbytes)

	if tag != 35 {
		t.Errorf("%s: Mismatch in tag value, Expected: 35, Got: %d",
			FuncName(), tag)
	}

	if val != x {
		t.Errorf("%s: Mismatch in output value, Expected: %d, Got: %d",
			FuncName(), x, val)
	}
}

func TestSint64(t *testing.T) {
	var x int64 = -19214912
	bytes := make([]byte, 0)
	var wtype ProtobufWireType
	var tag uint32

	t.Logf("%s: val is: %d", FuncName(), x)
	bytes, _ = TagPack(35, bytes)
	bytes[0] = bytes[0] | byte(PROTOBUF_WIRE_TYPE_VARINT)

	bytes, l := Sint64Pack(x, bytes)

	t.Log(FuncName()+":", "Packed length:", len(bytes), ", packed output:",
		bytes)

	used := ParseTagAndWiretype(l, bytes, &tag, &wtype)
	dbytes := bytes[used:]

	t.Log(FuncName(), ":Protobuf data obtained : ", dbytes)

	var i uint

	for i = 0; i < 10; i++ {
		if (dbytes[i] & 0x80) == 0 {
			break
		}
	}

	l = uint32(i + 1)
	val := Unzigzag64(ParseUint64(l, dbytes))

	if tag != 35 {
		t.Errorf("%s: Mismatch in tag value, Expected: 35, Got: %d",
			FuncName(), tag)
	}

	if int64(val) != x {
		t.Errorf("%s: Mismatch in output value, Expected: %d, Got: %d",
			FuncName(), x, int64(val))
	}
}

func TestSFixed64(t *testing.T) {
	var x int64 = -987654321
	bytes := make([]byte, 0)
	var wtype ProtobufWireType
	var tag uint32

	t.Logf("%s: val is: %d", FuncName(), x)
	bytes, _ = TagPack(35, bytes)
	bytes[0] |= PROTOBUF_WIRE_TYPE_32BIT

	bytes, l := Fixed64Pack(uint64(x), bytes)

	t.Log(FuncName()+":", "Packed length:", len(bytes), ", packed output:",
		bytes)

	used := ParseTagAndWiretype(l, bytes, &tag, &wtype)
	dbytes := bytes[used:]

	t.Log(FuncName(), ":Protobuf data obtained : ", dbytes)

	val := ParseFixedUint64(dbytes)

	if tag != 35 {
		t.Errorf("%s: Mismatch in tag value, Expected: 35, Got: %d",
			FuncName(), tag)
	}

	if int64(val) != int64(x) {

		fmt.Println("float value is ", float64(int64(val)))
		t.Errorf("%s: Mismatch in output value, Expected: %d, Got: %d",
			FuncName(), x, int64(val))
	}
}

func TestFloat(t *testing.T) {
	var x float32 = 2018.5678
	bytes := make([]byte, 0)
	var wtype ProtobufWireType
	var tag uint32

	t.Logf("%s: val is: %f", FuncName(), x)
	bytes, _ = TagPack(35, bytes)
	bytes[0] |= PROTOBUF_WIRE_TYPE_32BIT

	bytes, l := Fixed32Pack(math.Float32bits(x), bytes)

	t.Log(FuncName()+":", "Packed length:", len(bytes), ", packed output:",
		bytes)

	used := ParseTagAndWiretype(l, bytes, &tag, &wtype)
	dbytes := bytes[used:]

	t.Log(FuncName(), ":Protobuf data obtained : ", dbytes)

	val := ParseFixedUint32(dbytes)

	if tag != 35 {
		t.Errorf("%s: Mismatch in tag value, Expected: 35, Got: %d",
			FuncName(), tag)
	}

	if math.Float32frombits(val) != x {
		t.Errorf("%s: Mismatch in output value, Expected: %f, Got: %f",
			FuncName(), x, math.Float32frombits(val))
	}
}

func TestDouble(t *testing.T) {
	var x float64 = 123456789.123456789
	bytes := make([]byte, 0)
	var wtype ProtobufWireType
	var tag uint32

	t.Logf("%s: val is: %f", FuncName(), x)
	bytes, _ = TagPack(35, bytes)
	bytes[0] |= PROTOBUF_WIRE_TYPE_32BIT

	bytes, l := Fixed64Pack(math.Float64bits(x), bytes)

	t.Log(FuncName()+":", "Packed length:", len(bytes), ", packed output:",
		bytes)

	used := ParseTagAndWiretype(l, bytes, &tag, &wtype)
	dbytes := bytes[used:]

	t.Log(FuncName(), ":Protobuf data obtained : ", dbytes)

	val := ParseFixedUint64(dbytes)

	if tag != 35 {
		t.Errorf("%s: Mismatch in tag value, Expected: 35, Got: %d",
			FuncName(), tag)
	}

	if math.Float64frombits(val) != x {
		t.Errorf("%s: Mismatch in output value, Expected: %f, Got: %f",
			FuncName(), x, math.Float64frombits(val))
	}
}
