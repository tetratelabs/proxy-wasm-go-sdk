package properties

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"time"
	"unsafe"
)

// serializeBool converts a boolean value to a byte slice representation.
func serializeBool(value bool) []byte {
	if value {
		return []byte{1}
	}
	return []byte{0}
}

// deserializeBool converts a byte slice back to a boolean value.
func deserializeBool(bs []byte) (bool, error) {
	if len(bs) == 0 {
		return false, nil
	}
	if len(bs) != 1 {
		return false, fmt.Errorf("invalid byte slice length for boolean deserialization")
	}
	return bs[0] != 0, nil
}

// serializeByteMap serializes a map where keys are strings and values are raw byte slices.
// The resulting byte slice can be used for efficient storage or transmission.
//   - keys are always string
//   - values are raw byte slices
func serializeByteSliceMap(data map[string][]byte) []byte {
	if len(data) == 0 {
		return []byte{}
	}

	totalSize := 4
	for key, value := range data {
		totalSize += 4 + len(key) + 1 + 4 + len(value) + 1
	}
	bs := make([]byte, totalSize)
	binary.LittleEndian.PutUint32(bs[0:4], uint32(len(data)))
	var sizeIndex = 4
	var dataIndex = 4 + 4*2*len(data)
	for key, value := range data {
		binary.LittleEndian.PutUint32(bs[sizeIndex:sizeIndex+4], uint32(len(key)))
		sizeIndex += 4
		copy(bs[dataIndex:dataIndex+len(key)], key)
		dataIndex += len(key) + 1
		binary.LittleEndian.PutUint32(bs[sizeIndex:sizeIndex+4], uint32(len(value)))
		sizeIndex += 4
		copy(bs[dataIndex:dataIndex+len(value)], value)
		dataIndex += len(value) + 1
	}
	return bs
}

// deserializeByteMap deserializes the byte slice to key value map, used for mixed type maps
//   - keys are always string
//   - value are raw byte strings that need further parsing
func deserializeByteSliceMap(bs []byte) map[string][]byte { //nolint:unused
	ret := make(map[string][]byte)
	if len(bs) == 0 {
		return ret
	}

	numHeaders := binary.LittleEndian.Uint32(bs[0:4])
	var sizeIndex = 4
	var dataIndex = 4 + 4*2*int(numHeaders)
	for i := 0; i < int(numHeaders); i++ {
		keySize := int(binary.LittleEndian.Uint32(bs[sizeIndex : sizeIndex+4]))
		sizeIndex += 4
		keyPtr := bs[dataIndex : dataIndex+keySize]
		key := *(*string)(unsafe.Pointer(&keyPtr))
		dataIndex += keySize + 1

		valueSize := int(binary.LittleEndian.Uint32(bs[sizeIndex : sizeIndex+4]))
		sizeIndex += 4
		valuePtr := bs[dataIndex : dataIndex+valueSize]
		value := *(*[]byte)(unsafe.Pointer(&valuePtr))
		dataIndex += valueSize + 1
		ret[key] = value
	}
	return ret
}

// serializeByteSliceSlice serializes a slice of byte slices into a single byte slice.
// The resulting byte slice can be used for efficient storage or transmission.
// Each byte slice in the input is prefixed with its length, allowing for efficient deserialization.
func serializeByteSliceSlice(slices [][]byte) []byte {
	if len(slices) == 0 {
		return []byte{}
	}

	totalSize := 4
	for _, slice := range slices {
		totalSize += 8 + len(slice) + 2
	}
	bs := make([]byte, totalSize)
	binary.LittleEndian.PutUint32(bs[:4], uint32(len(slices)))
	idx := 4
	dataIdx := 4 + 8*len(slices)
	for _, slice := range slices {
		binary.LittleEndian.PutUint64(bs[idx:idx+8], uint64(len(slice)))
		idx += 8
		copy(bs[dataIdx:dataIdx+len(slice)], slice)
		dataIdx += len(slice) + 2
	}
	return bs
}

// deserializeByteSliceSlice deserializes the given bytes to string slice.
func deserializeByteSliceSlice(bs []byte) [][]byte {
	if len(bs) == 0 {
		return [][]byte{}
	}
	numStrings := int(binary.LittleEndian.Uint32(bs[:4]))
	ret := make([][]byte, numStrings)
	idx := 4
	dataIdx := 4 + 8*numStrings
	for i := 0; i < numStrings; i++ {
		strLen := int(binary.LittleEndian.Uint64(bs[idx : idx+8]))
		idx += 8
		ret[i] = bs[dataIdx : dataIdx+strLen]
		dataIdx += strLen + 2
	}
	return ret
}

// serializeFloat64 serializes the given float64 to bytes.
// The resulting byte slice can be used for efficient storage or transmission.
func serializeFloat64(value float64) []byte {
	bits := math.Float64bits(value)
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, bits)
	return bs
}

// deserializeFloat64 deserializes the given bytes to float64.
func deserializeFloat64(bs []byte) float64 {
	bits := binary.LittleEndian.Uint64(bs)
	float := math.Float64frombits(bits)
	return float
}

// serializeProtoStringSlice serializes a slice of strings into a protobuf-like encoded byte slice.
// The resulting byte slice can be used for efficient storage or transmission.
// Each string in the slice is prefixed with its length, allowing for efficient deserialization.
func serializeProtoStringSlice(strs []string) []byte {
	var bs []byte
	if len(strs) == 0 {
		return bs
	}

	for _, str := range strs {
		if len(str) > 255 {
			panic("string length exceeds 255 characters")
		}
		bs = append(bs, 0x00)
		bs = append(bs, byte(len(str)))
		bs = append(bs, []byte(str)...)
	}
	return bs
}

// deserializeProtoStringSlice deserializes a protobuf encoded string slice
func deserializeProtoStringSlice(bs []byte) []string {
	ret := make([]string, 0)
	if len(bs) == 0 {
		return ret
	}
	i := 0
	for i < len(bs) {
		i++
		length := int(bs[i])
		i++
		str := string(bs[i : i+length])
		ret = append(ret, str)
		i += length
	}
	return ret
}

// serializeStringMap serializes a map of strings to a byte slice.
//   - keys are always string
//   - values are always string
//
// The resulting byte slice starts with a 4-byte representation of the number of key-value pairs.
// This is followed by a series of 4-byte representations of the sizes of the keys and values.
// Finally, the actual key and value data are appended.
func serializeStringMap(m map[string]string) []byte {
	headerBytes := make([]byte, 4)
	if len(m) == 0 {
		return headerBytes
	}
	var buf bytes.Buffer
	numHeaders := uint32(len(m))
	binary.LittleEndian.PutUint32(headerBytes, numHeaders)
	buf.Write(headerBytes)
	var sizeData bytes.Buffer
	var data bytes.Buffer
	for key, value := range m {
		keySize := uint32(len(key))
		keySizeBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(keySizeBytes, keySize)
		sizeData.Write(keySizeBytes)
		keyData := *(*[]byte)(unsafe.Pointer(&key))
		data.Write(keyData)
		data.WriteByte(0)
		valueSize := uint32(len(value))
		valueSizeBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(valueSizeBytes, valueSize)
		sizeData.Write(valueSizeBytes)
		valueData := *(*[]byte)(unsafe.Pointer(&value))
		data.Write(valueData)
		data.WriteByte(0)
	}
	buf.Write(sizeData.Bytes())
	buf.Write(data.Bytes())
	return buf.Bytes()
}

// deserializeStringMap deserializes the bytes to key value map, used for string only type maps
//   - keys are always string
//   - value are always string
func deserializeStringMap(bs []byte) map[string]string {
	numHeaders := binary.LittleEndian.Uint32(bs[0:4])
	if numHeaders == 0 {
		return map[string]string{}
	}

	var sizeIndex = 4
	var dataIndex = 4 + 4*2*int(numHeaders)
	ret := make(map[string]string, numHeaders)
	for i := 0; i < int(numHeaders); i++ {
		keySize := int(binary.LittleEndian.Uint32(bs[sizeIndex : sizeIndex+4]))
		sizeIndex += 4
		keyPtr := bs[dataIndex : dataIndex+keySize]
		key := *(*string)(unsafe.Pointer(&keyPtr))
		dataIndex += keySize + 1

		valueSize := int(binary.LittleEndian.Uint32(bs[sizeIndex : sizeIndex+4]))
		sizeIndex += 4
		valuePtr := bs[dataIndex : dataIndex+valueSize]
		value := *(*string)(unsafe.Pointer(&valuePtr))
		dataIndex += valueSize + 1
		ret[key] = value
	}
	return ret
}

// serializeStringSlice serializes a slice of strings into a single byte slice.
// The resulting byte slice can be used for efficient storage or transmission.
// Each string in the input is prefixed with its length, allowing for efficient deserialization.
func serializeStringSlice(strings []string) []byte {
	if len(strings) == 0 {
		return make([]byte, 4)
	}
	totalSize := 4
	for _, str := range strings {
		totalSize += 8 + len(str) + 2
	}
	bs := make([]byte, totalSize)
	binary.LittleEndian.PutUint32(bs[:4], uint32(len(strings)))
	idx := 4
	dataIdx := 4 + 8*len(strings)
	for _, str := range strings {
		binary.LittleEndian.PutUint64(bs[idx:idx+8], uint64(len(str)))
		idx += 8
		copy(bs[dataIdx:dataIdx+len(str)], str)
		dataIdx += len(str) + 2
	}
	return bs
}

// deserializeStringSlice deserializes the given byte slice to string slice.
func deserializeStringSlice(bs []byte) []string {
	numStrings := int(binary.LittleEndian.Uint32(bs[:4]))
	if numStrings == 0 {
		return []string{}
	}
	ret := make([]string, numStrings)
	idx := 4
	dataIdx := 4 + 8*numStrings
	for i := 0; i < numStrings; i++ {
		strLen := int(binary.LittleEndian.Uint64(bs[idx : idx+8]))
		idx += 8
		ret[i] = string(bs[dataIdx : dataIdx+strLen])
		dataIdx += strLen + 2
	}
	return ret
}

// serializeTimestamp serializes the given timestamp to bytes.
// The resulting byte slice can be used for efficient storage or transmission.
func serializeTimestamp(timestamp time.Time) []byte {
	nanos := timestamp.UnixNano()
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(nanos))
	return bs
}

// deserializeTimestamp deserializes the given bytes to timestamp.
func deserializeTimestamp(bs []byte) time.Time {
	nanos := int64(binary.LittleEndian.Uint64(bs))
	return time.Unix(0, nanos)
}

// serializeUint64 serializes the given uint64 to bytes.
// The resulting byte slice can be used for efficient storage or transmission.
func serializeUint64(value uint64) []byte {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, value)
	return bs
}

// deserializeUint64 deserializes  the given bytes to uint64.
func deserializeUint64(bs []byte) uint64 {
	return binary.LittleEndian.Uint64(bs)
}
