package properties

import (
	"encoding/binary"
	"math"
	"time"
	"unsafe"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

// getPropertyString returns a string property.
func getPropertyString(path []string) (string, error) {
	b, err := proxywasm.GetProperty(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// getPropertyUint64 returns a uint64 property.
func getPropertyUint64(path []string) (uint64, error) {
	b, err := proxywasm.GetProperty(path)
	if err != nil {
		return 0, err
	}

	return deserializeToUint64(b), nil
}

// getPropertyFloat64 returns a float64 property.
func getPropertyFloat64(path []string) (float64, error) {
	b, err := proxywasm.GetProperty(path)
	if err != nil {
		return 0, err
	}

	return deserializeToFloat64(b), nil
}

// getPropertyBool returns a bool property.
func getPropertyBool(path []string) (bool, error) {
	b, err := proxywasm.GetProperty(path)
	if err != nil {
		return false, err
	}

	return b[0] != 0, nil
}

// getPropertyTimestamp returns a timestamp property.
func getPropertyTimestamp(path []string) (time.Time, error) {
	b, err := proxywasm.GetProperty(path)
	if err != nil {
		return time.Now(), err
	}

	return deserializeToTimestamp(b), nil
}

// getPropertyByteSliceMap retrieves a complex property object as a map of byte slices.
// to be used when dealing with mixed type properties
func getPropertyByteSliceMap(path []string) (map[string][]byte, error) {
	b, err := proxywasm.GetProperty(path)
	if err != nil {
		return nil, err
	}

	return deserializeToByteMap(b), nil
}

// getPropertyStringMap retrieves a complex property object as a map of string
// to be used when dealing with string only type properties.
func getPropertyStringMap(path []string) (map[string]string, error) {
	b, err := proxywasm.GetProperty(path)
	if err != nil {
		return nil, err
	}

	return deserializeToStringMap(b), nil
}

// getPropertyStringSlice retrieves a  complex property object as a string slice.
func getPropertyStringSlice(path []string) ([]string, error) {
	b, err := proxywasm.GetProperty(path)
	if err != nil {
		return nil, err
	}

	return deserializeToStringSlice(b), nil
}

// deserializeToStringSlice deserializes the given byte slice to string slice.
func deserializeToStringSlice(bs []byte) []string {
	numStrings := int(binary.LittleEndian.Uint32(bs[:4]))
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

// getPropertyByteSliceSlice retrieves a complex property object as a string slice.
func getPropertyByteSliceSlice(path []string) ([][]byte, error) {
	b, err := proxywasm.GetProperty(path)
	if err != nil {
		return nil, err
	}
	return deserializeToByteSliceSlice(b), nil
}

// deserializeToByteSliceSlice deserializes the given bytes to string slice.
func deserializeToByteSliceSlice(bs []byte) [][]byte {
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

// deserializeToUint64 deserializes  the given bytes  to uint64.
func deserializeToUint64(bytes []byte) uint64 {
	return binary.LittleEndian.Uint64(bytes)
}

// deserializeToFloat64 deserializes the given bytes to float64.
func deserializeToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

// deserializeToTimestamp deserializes the given bytes to timestamp.
func deserializeToTimestamp(data []byte) time.Time {
	nanos := int64(binary.LittleEndian.Uint64(data))
	return time.Unix(0, nanos)
}

// deserializeProtobufToStringSlice deserializes the given bytes to string slice.
func deserializeProtobufToStringSlice(data []byte) []string {
	var ret []string
	i := 0
	for i < len(data) {
		i++
		length := int(data[i])
		i++
		str := string(data[i : i+length])
		ret = append(ret, str)
		i += length
	}
	return ret
}

// deserializeToByteMap deserializes the byte slice to key value map, used for mixed type maps
//   - keys are always string
//   - value are raw byte strings that need further parsing
func deserializeToByteMap(bs []byte) map[string][]byte {
	numHeaders := binary.LittleEndian.Uint32(bs[0:4])
	var sizeIndex = 4
	var dataIndex = 4 + 4*2*int(numHeaders)
	ret := make(map[string][]byte)
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

// deserializeToStringMap deserializes the bytes to key value map, used for string only type maps
//   - keys are always string
//   - value are always string
func deserializeToStringMap(bs []byte) map[string]string {
	numHeaders := binary.LittleEndian.Uint32(bs[0:4])
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
