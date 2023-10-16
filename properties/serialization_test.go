package properties

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSerializeBool(t *testing.T) {
	tests := []struct {
		input    bool
		expected []byte
	}{
		{true, []byte{1}},
		{false, []byte{0}},
	}

	for _, test := range tests {
		result := serializeBool(test.input)
		require.Equal(t, test.expected, result)
	}
}

func TestDeserializeBool(t *testing.T) {
	tests := []struct {
		input    []byte
		expected bool
		err      error
	}{
		{[]byte{1}, true, nil},
		{[]byte{0}, false, nil},
		{[]byte{}, false, nil},
		{[]byte{1, 0}, false, fmt.Errorf("invalid byte slice length for boolean deserialization")},
	}

	for _, test := range tests {
		result, err := deserializeBool(test.input)
		require.Equal(t, test.expected, result)

		if test.err != nil {
			require.EqualError(t, err, test.err.Error())
		} else {
			require.NoError(t, err)
		}
	}
}

func TestSerializeAndDeserializeByteSliceMap(t *testing.T) {
	tests := []struct {
		input map[string][]byte
	}{
		{
			map[string][]byte{
				"key1": []byte("value1"),
				"key2": []byte("value2"),
			},
		},
		{
			map[string][]byte{
				"hello": []byte("world"),
				"foo":   []byte("bar"),
			},
		},
		{
			map[string][]byte{},
		},
	}

	for _, test := range tests {
		serialized := serializeByteSliceMap(test.input)
		deserialized := deserializeByteSliceMap(serialized)
		require.Equal(t, test.input, deserialized)
	}
}

func TestSerializeAndDeserializeByteSliceSlice(t *testing.T) {
	tests := []struct {
		input [][]byte
	}{
		{
			[][]byte{
				[]byte("slice1"),
				[]byte("slice2"),
			},
		},
		{
			[][]byte{
				[]byte("hello"),
				[]byte("world"),
				[]byte("foo"),
				[]byte("bar"),
			},
		},
		{
			[][]byte{},
		},
	}

	for _, test := range tests {
		serialized := serializeByteSliceSlice(test.input)
		deserialized := deserializeByteSliceSlice(serialized)
		require.Equal(t, test.input, deserialized)
	}
}

func TestSerializeAndDeserializeFloat64(t *testing.T) {
	tests := []struct {
		input float64
	}{
		{input: 123.456},
		{input: 0.0},
		{input: -987.654},
		{input: math.Pi},
		{input: math.MaxFloat64},
		{input: math.SmallestNonzeroFloat64},
	}

	for _, test := range tests {
		serialized := serializeFloat64(test.input)
		deserialized := deserializeFloat64(serialized)
		require.Equal(t, test.input, deserialized)
	}
}

func TestSerializeAndDeserializeProtoStringSlice(t *testing.T) {
	tests := []struct {
		input []string
	}{
		{input: []string{"hello", "world"}},
		{input: []string{"envoy.formatter.metadata", "envoy.formatter", "envoy.extensions.formatter.metadata.v3.Metadata"}},
		{input: []string{"a", "ab", "abc", "abcd"}},
		{input: []string{}},
		{input: []string{"", "empty", ""}},
	}

	for _, test := range tests {
		serialized := serializeProtoStringSlice(test.input)
		deserialized := deserializeProtoStringSlice(serialized)
		require.Equal(t, test.input, deserialized)
	}
}

func TestSerializeAndDeserializeStringMap(t *testing.T) {
	tests := []struct {
		input map[string]string
	}{
		{input: map[string]string{"hello": "world", "foo": "bar"}},
		{input: map[string]string{"key": "value"}},
		{input: map[string]string{}},
		{input: map[string]string{"empty": "", "": "empty"}},
	}

	for _, test := range tests {
		serialized := serializeStringMap(test.input)
		deserialized := deserializeStringMap(serialized)
		require.Equal(t, test.input, deserialized)
	}
}

func TestSerializeAndDeserializeStringSlice(t *testing.T) {
	tests := []struct {
		input []string
	}{
		{input: []string{"hello", "world", "foo", "bar"}},
		{input: []string{"key", "value"}},
		{input: []string{}},
		{input: []string{"", "empty", ""}},
	}

	for _, test := range tests {
		serialized := serializeStringSlice(test.input)
		deserialized := deserializeStringSlice(serialized)
		require.Equal(t, test.input, deserialized)
	}
}

func TestSerializeAndDeserializeTimestamp(t *testing.T) {
	tests := []struct {
		input time.Time
	}{
		{input: time.Now().UTC()},
		{input: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
		{input: time.Date(1990, 5, 15, 5, 5, 5, 5, time.UTC)},
	}

	for _, test := range tests {
		serialized := serializeTimestamp(test.input)
		deserialized := deserializeTimestamp(serialized)
		if !test.input.Equal(deserialized) {
			t.Errorf("Expected %v, got %v", test.input, deserialized)
		}
	}
}

func TestSerializeAndDeserializeUint64(t *testing.T) {
	tests := []struct {
		input uint64
	}{
		{input: 1234567890},
		{input: 0},
		{input: 9876543210},
		{input: 18446744073709551615},
	}

	for _, test := range tests {
		serialized := serializeUint64(test.input)
		deserialized := deserializeUint64(serialized)
		require.Equal(t, test.input, deserialized)
	}
}
