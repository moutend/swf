package swf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
)

type DefineShapeContent struct {
	ID          uint16
	ShapeBounds *RectContent
	IDBuffer    *bytes.Buffer
}

func UnmarshalDefineShape(input io.Reader) (*DefineShapeContent, error) {
	// Read the ID. (Fixed length, 2 bytes, little endian)
	id := &bytes.Buffer{}

	idLength, err := io.CopyN(id, input, 2)

	if err != nil {
		return nil, err
	}
	if idLength != 2 {
		return nil, fmt.Errorf("broken id")
	}

	var idUint16 uint16

	{
		w := &bytes.Buffer{}
		r := io.TeeReader(id, w)

		if err := binary.Read(r, binary.LittleEndian, &idUint16); err != nil {
			return nil, err
		}

		id = w
	}

	shapeBounds, err := UnmarshalRect(input)

	if err != nil {
		return nil, err
	}

	result := &DefineShapeContent{
		ID:          idUint16,
		ShapeBounds: shapeBounds,
		IDBuffer:    id,
	}

	return result, nil
}

type RectContent struct {
	BitsPerField int
	MinX         uint32
	MaxX         uint32
	MinY         uint32
	MaxY         uint32
	DataBuffer   *bytes.Buffer
}

func UnmarshalRect(input io.Reader) (*RectContent, error) {
	data := &bytes.Buffer{}

	dataLength, err := io.CopyN(data, input, 1)

	if err != nil {
		return nil, err
	}
	if dataLength != 1 {
		return nil, fmt.Errorf("broken data")
	}

	bits := fmt.Sprintf("%08b", data.Bytes()[0])

	i64, err := strconv.ParseInt(bits[:5], 2, 64)

	if err != nil {
		return nil, err
	}

	bitsPerField := int(i64)
	remainingBits := bitsPerField*4 - 3
	requiredBits := 0

	for {
		if requiredBits >= remainingBits {
			break
		}

		requiredBits += 8
	}

	dataLength, err = io.CopyN(data, input, int64(requiredBits/8))

	if err != nil {
		return nil, err
	}
	if dataLength != int64(requiredBits/8) {
		return nil, fmt.Errorf("broken data")
	}

	var s string

	for _, b := range data.Bytes() {
		s += fmt.Sprintf("%08b", b)
	}

	values := make([]uint32, 4)
	start := 5

	for i := 0; i < 4; i++ {
		i64, err := strconv.ParseInt(s[start:start+bitsPerField], 2, 64)

		if err != nil {
			return nil, err
		}

		values[i] = uint32(i64)
		fmt.Println("@@@value", values[i])
		start += bitsPerField
	}

	result := &RectContent{
		BitsPerField: bitsPerField,
		MinX:         values[0],
		MaxX:         values[1],
		MinY:         values[2],
		MaxY:         values[3],
		DataBuffer:   data,
	}

	return result, nil
}
