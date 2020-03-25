package runtime

import (
	"encoding/binary"
)

func deserializeMap(bs []byte) [][2]string {
	numHeaders := binary.LittleEndian.Uint32(bs[0:4])
	sizes := make([]int, numHeaders*2)
	for i := 0; i < len(sizes); i++ {
		s := 4 + i*4
		sizes[i] = int(binary.LittleEndian.Uint32(bs[s : s+4]))
	}

	var sizeIndex int
	var dataIndex = 4 * (1 + 2*int(numHeaders))
	ret := make([][2]string, numHeaders)
	for i := range ret {
		keySize := sizes[sizeIndex]
		sizeIndex++
		key := string(bs[dataIndex : dataIndex+keySize]) // TODO: zero alloc
		dataIndex += keySize + 1

		valueSize := sizes[sizeIndex]
		sizeIndex++
		value := string(bs[dataIndex : dataIndex+valueSize]) // TODO: zero alloc
		dataIndex += valueSize + 1
		ret[i] = [2]string{key, value}
	}
	return ret
}

func serializeMap(ms [][2]string) []byte {
	size := 4
	for _, m := range ms {
		// key/value's bytes + len * 2 (8 bytes) + nil * 2 (2 bytes)
		size += len(m[0]) + len(m[1]) + 10
	}

	ret := make([]byte, size)
	binary.LittleEndian.PutUint32(ret[0:4], uint32(len(ms)))

	var base = 4
	for _, m := range ms {
		binary.LittleEndian.PutUint32(ret[base:base+4], uint32(len(m[0])))
		base += 4
		binary.LittleEndian.PutUint32(ret[base:base+4], uint32(len(m[1])))
		base += 4
	}

	for _, m := range ms {
		for i := 0; i < len(m[0]); i++ {
			ret[base] = m[0][i]
			base++
		}
		base++ // nil

		for i := 0; i < len([]byte(m[1])); i++ {
			ret[base] = m[1][i]
			base++
		}
		base++ // nil
	}
	return ret
}
