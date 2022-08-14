package swf

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadRectangle(t *testing.T) {
	data := []byte{0x78, 0x00, 0x03, 0xe8, 0x00, 0x00, 0x13, 0x88, 0x00}
	buffer := bytes.NewBuffer(data)

	rectangle, err := ReadRectangle(buffer)

	require.NoError(t, err)
	require.NotNil(t, rectangle)

	require.Equal(t, 15, rectangle.BitsPerField)
	require.Equal(t, uint32(0), rectangle.MinX)
	require.Equal(t, uint32(8000), rectangle.MaxX)
	require.Equal(t, uint32(0), rectangle.MinY)
	require.Equal(t, uint32(10000), rectangle.MaxY)

	require.Equal(t, data, rectangle.Bytes())
	require.Equal(t, "Rectangle{0 8000 0 10000}", rectangle.String())

	d1, err := rectangle.Serialize()

	require.NoError(t, err)
	require.Equal(t, data, d1)
}

func TestReadRGB(t *testing.T) {
	input := bytes.NewBuffer([]byte{0x12, 0x34, 0x56, 0x78, 0x90})

	color, err := ReadRGB(input)

	require.NoError(t, err)
	require.NotNil(t, color)

	require.Equal(t, uint8(0x12), color.Red)
	require.Equal(t, uint8(0x34), color.Green)
	require.Equal(t, uint8(0x56), color.Blue)
	require.Equal(t, uint8(0xff), color.Alpha)

	require.Equal(t, "RGB{18, 52, 86}", color.String())

	require.Len(t, color.Bytes(), 3)
	require.Equal(t, []byte{0x12, 0x34, 0x56}, color.Bytes())

	color.Red = 0x11
	color.Green = 0x22
	color.Blue = 0x33

	data, err := color.Serialize()

	require.NoError(t, err)
	require.NotNil(t, data)

	require.Len(t, data, 3)
	require.Equal(t, data, []byte{0x11, 0x22, 0x33})
}
