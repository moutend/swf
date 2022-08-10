package swf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

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
