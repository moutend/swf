package swf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
)

type Rectangle struct {
	BitsPerField int
	MinX         uint32
	MaxX         uint32
	MinY         uint32
	MaxY         uint32

	data *bytes.Buffer
}

func (r *Rectangle) String() string {
	return fmt.Sprintf("Rectangle{%d %d %d %d}", r.MinX, r.MaxX, r.MinY, r.MaxY)
}

func (r *Rectangle) Serialize() ([]byte, error) {
	s := fmt.Sprintf("%05b", uint8(r.BitsPerField))
	s += fmt.Sprintf("%032b", r.MinX)[32-r.BitsPerField:]
	s += fmt.Sprintf("%032b", r.MaxX)[32-r.BitsPerField:]
	s += fmt.Sprintf("%032b", r.MinY)[32-r.BitsPerField:]
	s += fmt.Sprintf("%032b", r.MaxY)[32-r.BitsPerField:]

	padding := len(s) % 8

	for i := 0; i < padding; i++ {
		s += "0"
	}

	var data []byte

	for i := 0; i <= len(s); i += 8 {
		i64, err := strconv.ParseInt(s[i:i+8], 2, 64)

		if err != nil {
			return nil, err
		}

		data = append(data, byte(i64))
	}

	return data, nil
}

func (r *Rectangle) Data() []byte {
	var data []byte

	data = append(data, r.data.Bytes()...)

	return data
}

func ReadRectangle(src io.Reader) (*Rectangle, error) {
	data := &bytes.Buffer{}

	dataLength, err := io.CopyN(data, src, 1)

	if err != nil {
		return nil, err
	}
	if dataLength != 1 {
		return nil, fmt.Errorf("ReadRectangle: broken data")
	}

	// Read the first 5 bits.
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

	dataLength, err = io.CopyN(data, src, int64(requiredBits/8))

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

		start += bitsPerField
	}

	rectangle := &Rectangle{
		BitsPerField: bitsPerField,
		MinX:         values[0],
		MaxX:         values[1],
		MinY:         values[2],
		MaxY:         values[3],
		data:         data,
	}

	return rectangle, nil
}

type Color struct {
	Red   uint8
	Green uint8
	Blue  uint8
	Alpha uint8

	data *bytes.Buffer
}

func (c *Color) String() string {
	if len(c.data.Bytes()) == 3 {
		return fmt.Sprintf("Color{0x%x 0x%x 0x%x}", c.Red, c.Green, c.Blue)
	}

	return fmt.Sprintf("Color{0x%x 0x%x 0x%x 0x%x}", c.Red, c.Green, c.Blue, c.Alpha)
}

func (c *Color) Serialize() ([]byte, error) {
	switch len(c.data.Bytes()) {
	case 3:
		return []byte{c.Red, c.Green, c.Blue}, nil
	case 4:
		return []byte{c.Red, c.Green, c.Blue, c.Alpha}, nil
	default:
		return nil, fmt.Errorf("broken Color")
	}
}

func (c *Color) Data() []byte {
	var data []byte

	data = append(data, c.data.Bytes()...)

	return data
}

func ReadRGB(src io.Reader) (*Color, error) {
	w := &bytes.Buffer{}
	r := io.TeeReader(src, w)

	var red, green, blue uint8

	if err := binary.Read(r, binary.LittleEndian, &red); err != nil {
		return nil, fmt.Errorf("failed to read RGB color: %w", err)
	}
	if err := binary.Read(r, binary.LittleEndian, &green); err != nil {
		return nil, fmt.Errorf("failed to read RGB color: %w", err)
	}
	if err := binary.Read(r, binary.LittleEndian, &blue); err != nil {
		return nil, fmt.Errorf("failed to read RGB color: %w", err)
	}

	return &Color{red, green, blue, 0xFF, w}, nil
}

func ReadRGBA(src io.Reader) (*Color, error) {
	w := &bytes.Buffer{}
	r := io.TeeReader(src, w)

	var red, green, blue, alpha uint8

	if err := binary.Read(r, binary.LittleEndian, &red); err != nil {
		return nil, fmt.Errorf("failed to read RGBA color: %w", err)
	}
	if err := binary.Read(r, binary.LittleEndian, &green); err != nil {
		return nil, fmt.Errorf("failed to read RGBA color: %w", err)
	}
	if err := binary.Read(r, binary.LittleEndian, &blue); err != nil {
		return nil, fmt.Errorf("failed to read RGBA color: %w", err)
	}

	return &Color{red, green, blue, alpha, w}, nil
}

type ShapeStyles struct {
	data *bytes.Buffer
}

func ReadShapeStyles(input io.Reader) (*ShapeStyles, error) {
	// TODO: implement me!
	return nil, nil
}

type FillStyle struct {
}

func ReadFillStyle(src io.Reader, shapeVersion uint8) (*FillStyle, error) {
	return nil, nil
}
