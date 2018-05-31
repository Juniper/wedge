/*
 * Copyright (c) 2018, Juniper Networks, Inc.
 * All rights reserved.
 */

package util

import "sync"

type ProtobufLabel uint
type ProtobufType uint
type ProtobufWireType uint
type ProtobufArchType uint

const (
	/** A well-formed message must have exactly one of this field. */
	PROTOBUF_LABEL_REQUIRED = iota

	/**
	 * A well-formed message can have zero or one of this field (but not
	 * more than one).
	 */
	PROTOBUF_LABEL_OPTIONAL

	/**
	 * This field can be repeated any number of times (including zero) in a
	 * well-formed message. The order of the repeated values will be
	 * preserved.
	 */
	PROTOBUF_LABEL_REPEATED
)

const (
	PROTOBUF_TYPE_INT32    = iota /**< int32 */
	PROTOBUF_TYPE_SINT32          /**< signed int32 */
	PROTOBUF_TYPE_SFIXED32        /**< signed int32 (4 bytes) */
	PROTOBUF_TYPE_INT64           /**< int64 */
	PROTOBUF_TYPE_SINT64          /**< signed int64 */
	PROTOBUF_TYPE_SFIXED64        /**< signed int64 (8 bytes) */
	PROTOBUF_TYPE_UINT32          /**< unsigned int32 */
	PROTOBUF_TYPE_FIXED32         /**< unsigned int32 (4 bytes) */
	PROTOBUF_TYPE_UINT64          /**< unsigned int64 */
	PROTOBUF_TYPE_FIXED64         /**< unsigned int64 (8 bytes) */
	PROTOBUF_TYPE_FLOAT           /**< float */
	PROTOBUF_TYPE_DOUBLE          /**< double */
	PROTOBUF_TYPE_BOOL            /**< boolean */
	PROTOBUF_TYPE_ENUM            /**< enumerated type */
	PROTOBUF_TYPE_STRING          /**< UTF-8 or ASCII string */
	PROTOBUF_TYPE_BYTES           /**< arbitrary byte sequence */
	PROTOBUF_TYPE_MESSAGE
)

const (
	PROTOBUF_WIRE_TYPE_VARINT          = 0
	PROTOBUF_WIRE_TYPE_64BIT           = 1
	PROTOBUF_WIRE_TYPE_LENGTH_PREFIXED = 2
	PROTOBUF_WIRE_TYPE_32BIT           = 5
)

const (
	PROTOBUF_LITTLE_ENDIAN = iota
	PROTOBUF_BIG_ENDIAN
)

const (
	/** Set if the field is repeated and marked with the `packed` option. */
	PROTOBUF_FIELD_FLAG_PACKED = (1 << 0)

	/** Set if the field is marked with the `deprecated` option. */
	PROTOBUF_FIELD_FLAG_DEPRECATED = (1 << 1)

	/** Set if the field is a member of a oneof (union). */
	PROTOBUF_FIELD_FLAG_ONEOF = (1 << 2)
)

var endian_type ProtobufArchType = PROTOBUF_LITTLE_ENDIAN
var mu sync.RWMutex

/*
 * Function to set the endian-ness of the machine
 */
func SetEndian(val ProtobufArchType) {
	mu.Lock()
	defer mu.Unlock()
	endian_type = val
}

/*
 * Function to set the endian-ness of the machine
 */
func GetEndian() ProtobufArchType {
	mu.RLock()
	defer mu.RUnlock()
	return endian_type

}

/*
 * Function to get the tag size
 */
func getTagSize(number uint32) uint32 {
	if number < (1 << 4) {
		return 1
	} else if number < (1 << 11) {
		return 2
	} else if number < (1 << 18) {
		return 3
	} else if number < (1 << 25) {
		return 4
	} else {
		return 5
	}
}

/*
 * Return the number of bytes required to store a variable-length unsigned
 * 32-bit integer in base-128 varint encoding.
 */
func UInt32Size(v uint32) uint32 {
	if v < (1 << 7) {
		return 1
	} else if v < (1 << 14) {
		return 2
	} else if v < (1 << 21) {
		return 3
	} else if v < (1 << 28) {
		return 4
	} else {
		return 5
	}
}

/*
 * Return the number of bytes required to store a variable-length signed 32-bit
 * integer in base-128 varint encoding.
 *
 */
func Int32Size(v int32) uint32 {
	if v < 0 {
		return 10
	} else if v < (1 << 7) {
		return 1
	} else if v < (1 << 14) {
		return 2
	} else if v < (1 << 21) {
		return 3
	} else if v < (1 << 28) {
		return 4
	} else {
		return 5
	}
}

/*
 * Return the ZigZag-encoded 32-bit unsigned integer form of a 32-bit signed
 * integer.
 */
func zigzag32(v int32) uint32 {
	if v < 0 {
		return ((uint32)(-v))*2 - 1
	} else {
		return uint32(v * 2)
	}
}

/*
 * Return the number of bytes required to store a signed 32-bit integer,
 * converted to an unsigned 32-bit integer with ZigZag encoding, using base-128
 * varint encoding.
 */
func sInt32Size(v int32) uint32 {
	return UInt32Size(zigzag32(v))
}

/*
 * Return the number of bytes required to store a 64-bit unsigned integer in
 * base-128 varint encoding.
 */
func uint64Size(v uint64) uint32 {
	var upper_v uint32 = uint32(v >> 32)

	if upper_v == 0 {
		return UInt32Size(uint32(v))
	} else if upper_v < (1 << 3) {
		return 5
	} else if upper_v < (1 << 10) {
		return 6
	} else if upper_v < (1 << 17) {
		return 7
	} else if upper_v < (1 << 24) {
		return 8
	} else if upper_v < (1 << 31) {
		return 9
	} else {
		return 10
	}
}

/*
 * Return the ZigZag-encoded 64-bit unsigned integer form of a 64-bit signed
 * integer.
 */
func zigzag64(v int64) uint64 {
	if v < 0 {
		return uint64(-v)*2 - 1
	} else {
		return uint64(v) * 2
	}
}

/*
 * Return the number of bytes required to store a signed 64-bit integer.
 */
func sint64Size(v int64) uint32 {
	return uint64Size(zigzag64(v))
}

/*
 * Pack an unsigned 32-bit integer in base-128 varint encoding and return the
 * number of bytes written, which must be 5 or less.
 */
func Uint32Pack(value uint32, buf []byte) ([]byte, uint32) {
	var rv uint32 = 0

	if value >= 0x80 {
		buf = append(buf, byte(value|0x80))
		value = value >> 7
		rv++
		if value >= 0x80 {
			buf = append(buf, byte(value|0x80))
			value = value >> 7
			if value >= 0x80 {
				buf = append(buf, byte(value|0x80))
				value = value >> 7
				rv++
				if value >= 0x80 {
					buf = append(buf, byte(value|0x80))
					value = value >> 7
					rv++
				}
			}
		}
	}

	buf = append(buf, byte(value))
	rv++

	return buf, rv
}

/*
 * Pack a signed 32-bit integer and return the number of bytes written.
 * Negative numbers are encoded as two's complement 64-bit integers.
 */
func Int32Pack(value int32, buf []byte) ([]byte, uint32) {
	if value < 0 {
		buf = append(buf, uint8(value)|0x80,
			uint8(value>>7)|0x80,
			uint8(value>>14)|0x80,
			uint8(value>>21)|0x80,
			uint8(value>>28)|0x80,
			0xFF, 0xFF, 0xFF, 0xFF,
			0x01)
		return buf, 10
	} else {
		return Uint32Pack(uint32(value), buf)
	}
}

/*
 * Pack a signed 32-bit integer using ZigZag encoding and return the number of
 * bytes written.
 */
func Sint32Pack(value int32, buf []byte) ([]byte, uint32) {
	return Uint32Pack(zigzag32(value), buf)
}

/*
 * Pack a 64-bit unsigned integer using base-128 varint encoding and return the
 * number of bytes written.
 */
func Uint64Pack(value uint64, out []byte) ([]byte, uint32) {
	hi := uint32(value >> 32)
	lo := uint32(value)
	var rv uint32

	if hi == 0 {
		return Uint32Pack(uint32(lo), out)
	}

	out = append(out, uint8(lo)|0x80,
		uint8(lo>>7)|0x80,
		uint8(lo>>14)|0x80,
		uint8(lo>>21)|0x80)

	if hi < 8 {
		out = append(out, uint8(hi<<4)|uint8(lo>>28))
		return out, 5
	} else {
		out = append(out, uint8(hi&7<<4)|uint8(lo>>28)|0x80)
		hi = hi >> 3
	}
	rv = 5
	for hi >= 128 {
		out = append(out, uint8(hi|0x80))
		hi >>= 7
		rv++
	}
	out = append(out, uint8(hi))
	rv++
	return out, rv
}

/*
 * Pack a 64-bit signed integer in ZigZag encoding and return the number of
 * bytes written.
 */
func Sint64Pack(value int64, out []byte) ([]byte, uint32) {
	return Uint64Pack(zigzag64(value), out)
}

/*
 * Pack a 32-bit quantity in little-endian byte order. Used for protobuf wire
 * types fixed32, sfixed32, float.
 */
func Fixed32Pack(value uint32, out []byte) ([]byte, uint32) {

	out = append(out, uint8(value),
		uint8(value>>8),
		uint8(value>>16),
		uint8(value>>24))

	return out, 4
}

/*
 * Pack a 64-bit quantity in little-endian byte order. Used for protobuf wire
 * types fixed64, sfixed64, double.
 */
func Fixed64Pack(value uint64, out []byte) ([]byte, uint32) {

	out, _ = Fixed32Pack(uint32(value), out)
	out, _ = Fixed32Pack(uint32(value>>32), out)

	return out, 8
}

/*
 * Pack a boolean value as an integer and return the number of bytes written.
 */
func BooleanPack(value bool, out []byte) ([]byte, uint32) {
	var b byte
	if value == true {
		b = 1
	} else {
		b = 0
	}
	out = append(out, b)
	return out, 1
}

/*
 * Pack a string and return the number of bytes written. The
 */
func StringPack(str string, out []byte) ([]byte, uint32) {
	if str == "" {
		out = append(out, 0)
		return out, 1
	} else {
		var length uint32 = uint32(len(str))
		var rv uint32
		out, rv = Uint32Pack(length, out)

		for _, c := range str {
			out = append(out, byte(c))
		}

		return out, rv + length
	}
}

/*
 * Pack a sequence of bytes
 */
func BytesDataPack(bytes []byte, out []byte) ([]byte, uint32) {
	var length uint32 = uint32(len(bytes))
	var rv uint32
	out, rv = Uint32Pack(length, out)

	for _, c := range bytes {
		out = append(out, c)
	}

	return out, rv + length
}

/*
 * Pack a field tag.
 */
func TagPack(id uint32, out []byte) ([]byte, uint32) {
	if id < (1 << (32 - 3)) {
		return Uint32Pack(id<<3, out)
	} else {
		return Uint64Pack(uint64(id)<<3, out)
	}
}

/*
 * Get the minimum number of bytes required to pack a field value of a
 * particular type.
 */
func getTypeMinSize(t ProtobufType) uint32 {
	if t == PROTOBUF_TYPE_SFIXED32 ||
		t == PROTOBUF_TYPE_FIXED32 ||
		t == PROTOBUF_TYPE_FLOAT {
		return 4
	}
	if t == PROTOBUF_TYPE_SFIXED64 ||
		t == PROTOBUF_TYPE_FIXED64 ||
		t == PROTOBUF_TYPE_DOUBLE {
		return 8
	}

	return 1
}

/*
 * Parse the wire data and get the tag, type
 */
func ParseTagAndWiretype(length uint32, data []byte, tag_out *uint32,
	wiretype_out *ProtobufWireType) uint32 {
	var max_rv uint32

	if length > 5 {
		max_rv = 5
	} else {
		max_rv = length
	}

	var tag uint32 = uint32((uint8(data[0]) & 0x7f) >> 3)
	var shift uint = 4
	var rv uint32

	*wiretype_out = ProtobufWireType(data[0] & 7)
	if (data[0] & 0x80) == 0 {
		*tag_out = tag
		return 1
	}

	for rv = 1; rv < max_rv; rv++ {
		if r := uint8(data[rv]) & 0x80; r != 0 {
			tag = tag | uint32((uint8(data[rv])&0x7f)<<shift)
			shift += 7
		} else {
			tag = tag | uint32(uint8(data[rv])<<shift)
			*tag_out = tag
			return rv + 1
		}
	}
	return 0 /* error: bad header */
}

/*
 *  get prefix data length
 */
func ScanLengthPrefixData(length uint32, data []byte,
	prefix_len_out *uint32) uint32 {
	var hdr_max uint32

	if length < 5 {
		hdr_max = length
	} else {
		hdr_max = 5
	}

	var hdr_len uint32
	var val uint32 = 0
	var shift uint32 = 0
	var i uint32

	for i = 0; i < hdr_max; i++ {
		val = val | uint32((uint8(data[i])&0x7f)<<shift)
		shift += 7
		if (uint8(data[i]) & 0x80) == 0 {
			break
		}

	}

	if i == hdr_max {
		return 0
	}

	hdr_len = i + 1
	*prefix_len_out = hdr_len

	if hdr_len+val > length {
		return 0
	}

	return hdr_len + val
}

/*
 * parse uint32 integer
 */
func ParseUint32(length uint32, data []byte) uint32 {
	var rv uint32 = uint32(data[0]) & 0x7f

	if length > 1 {
		rv = rv | (uint32(data[1])&0x7f)<<7
		if length > 2 {
			rv = rv | (uint32(data[2])&0x7f)<<14
			if length > 3 {
				rv = rv | (uint32(data[3])&0x7f)<<21
				if length > 4 {
					rv = rv | uint32(data[4])<<28
				}
			}
		}
	}

	return rv
}

/*
 * parse int32 integer
 */
func ParseInt32(length uint32, data []byte) uint32 {
	return ParseUint32(length, data)
}

/*
 * unzigzag the integer
 */
func Unzigzag32(v uint32) int32 {
	if b := v & 1; b != 0 {
		return -int32(v>>1) - 1
	} else {
		return int32(v >> 1)
	}
}

/*
 * parse fixed uint32 integer
 */
func ParseFixedUint32(data []byte) uint32 {

	return uint32(data[0]) |
		(uint32(data[1]) << 8) |
		(uint32(data[2]) << 16) |
		(uint32(data[3]) << 24)

}

/*
 * parse uint64 integer
 */
func ParseUint64(length uint32, data []byte) uint64 {
	var shift, i uint32
	var rv uint64

	if length < 5 {
		return uint64(ParseUint32(length, data))
	}

	rv = (uint64(data[0] & 0x7f)) |
		(uint64(data[1]&0x7f) << 7) |
		(uint64(data[2]&0x7f) << 14) |
		(uint64(data[3]&0x7f) << 21)

	shift = 28
	for i = 4; i < length; i++ {
		rv = rv | ((uint64(data[i] & 0x7f)) << shift)
		shift += 7
	}

	return rv
}

func Unzigzag64(v uint64) int64 {
	if b := v & 1; b != 0 {
		return -int64(v>>1) - 1
	} else {
		return int64(v >> 1)
	}
}

func ParseFixedUint64(data []byte) uint64 {

	return uint64(ParseFixedUint32(data)) |
		(uint64(ParseFixedUint32(data[4:])) << 32)

}

func ParseBoolean(length uint32, data []byte) bool {
	var i uint32
	for i = 0; i < length; i++ {
		if b := uint8(data[i]) & 0x7f; b != 0 {
			return true
		}
	}

	return false
}

func ScanVarint(length uint32, data []byte) uint32 {
	var i uint32

	if length > 10 {
		length = 10
	}

	for i = 0; i < length; i++ {
		if (uint8(data[i]) & 0x80) == 0 {
			break
		}
	}

	if i == length {
		return 0
	}

	return i + 1
}
