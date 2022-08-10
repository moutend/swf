package swf

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadRGB(t *testing.T) {
	input := bytes.NewBuffer([]byte{0x12, 0x34, 0x56, 0x78, 0x90})

	color, err := ReadRGB(input)

	require.NoError(t, err)
	require.NotNil(t, color)

	require.Equal(t, uint8(0x12), color.Red)
	require.Equal(t, uint8(0x34), color.Green)
	require.Equal(t, uint8(0x56), color.Blue)
	require.Equal(t, uint8(0xff), color.Alpha)

	require.Equal(t, "Color{0x12 0x34 0x56}", color.String())

	require.Len(t, color.Data(), 3)
	require.Equal(t, []byte{0x12, 0x34, 0x56}, color.Data())

	color.Red = 0x11
	color.Green = 0x22
	color.Blue = 0x33

	data, err := color.Serialize()

	require.NoError(t, err)
	require.NotNil(t, data)

	require.Len(t, data, 3)
	require.Equal(t, data, []byte{0x11, 0x22, 0x33})
}
