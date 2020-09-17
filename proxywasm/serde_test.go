package proxywasm

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var serdeTestCases = []struct {
	maps  [][2]string
	bytes []byte
}{
	{
		maps: [][2]string{{"a", "A"}},
		bytes: []byte{
			1, 0, 0, 0,
			1, 0, 0, 0,
			1, 0, 0, 0,
			97, 0, 65, 0,
		},
	},
	{
		maps: [][2]string{{"a", "A"}, {"b", "B"}},
		bytes: []byte{
			2, 0, 0, 0,
			1, 0, 0, 0,
			1, 0, 0, 0,
			1, 0, 0, 0,
			1, 0, 0, 0,
			97, 0, 65, 0,
			98, 0, 66, 0,
		},
	},
	{
		maps: [][2]string{{"a", "ABCDEFG"}, {"@AB", "<1234"}},
		bytes: []byte{
			2, 0, 0, 0,
			1, 0, 0, 0,
			7, 0, 0, 0,
			3, 0, 0, 0,
			5, 0, 0, 0,
			97, 0,
			65, 66, 67, 68, 69, 70, 71, 0,
			64, 65, 66, 0,
			60, 49, 50, 51, 52, 0,
		},
	},
}

func Test_DeserializeMap(t *testing.T) {
	for i, c := range serdeTestCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, c.maps, DeserializeMap(c.bytes))
		})
	}
}

func Test_SerializeMap(t *testing.T) {
	for i, c := range serdeTestCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, c.bytes, SerializeMap(c.maps))
		})
	}
}
