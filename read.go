package swf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strconv"
)

const (
	SignatureUncompressed = `SWF`
	SignatureCompressed   = `SWC`
)

type Signature struct {
	Value string
	data  *bytes.Buffer
}

func (s *Signature) String() string {
	if s == nil {
		return "<nil>"
	}

	return fmt.Sprintf("Signature{%q}", s.Value)
}

func (s *Signature) Data() []byte {
	if s == nil || s.data == nil {
		return nil
	}

	var data []byte

	data = append(data, s.data.Bytes()...)

	return data
}

func (s *Signature) Serialize() ([]byte, error) {
	if s == nil {
		return nil, fmt.Errorf("cannot serialize because Signature is nil")
	}
	switch s.Value {
	case SignatureUncompressed:
		return []byte(`FWS`)
	case SignatureCompressed:
		return []byte(`CWS`), nil
	default:
		return nil, fmt.Errorf("invalid signature: %q", s.Value)
	}
}

func ReadSignature(src io.Reader) (*Signature, error) {
	data := &bytes.Buffer{}

	dataLength, err := io.CopyN(data, src, 3)

	if err != nil {
		return nil, err
	}
	if dataLength != 3 {
		return nil, fmt.Errorf("broken signature")
	}
	if data.String() != `FWS` && data.String() != `CWS` {
		return nil, fmt.Errorf("invalid signature: %q", data.String())
	}

	var value string

	if data.String() == `FWS` {
		value = SignatureUncompressed
	} else {
		value = SignatureCompressed
	}

	signature := &Signature{
		Value: value,
		data:  data,
	}

	return signature, nil
}

type FrameRate struct {
	Value float64
	data  *bytes.Buffer
}

func (f *FrameRate) String() string {
	if f == nil {
		return "<nil>"
	}

	return fmt.Sprintf("FrameRate{%.2f}", f.Value)
}

func (f *FrameRate) Data() []byte {
	if f == nil || f.data == nil {
		return nil
	}

	var data []byte

	data = append(data, f.data.Bytes()...)

	return data
}

func (f *FrameRate) Serialize() ([]byte, error) {
	if f == nil {
		return nil, fmt.Errorf("cannot serialize because FrameRate is nil")
	}

	a := uint8(f.Value)
	b := uint8((f.Value - math.Floor(f.Value)) * 100.0)

	return []byte{b, a}, nil
}

func ReadFrameRate(src io.Reader) (*FrameRate, error) {
	data := &bytes.Buffer{}

	frameRateLength, err := io.CopyN(frameRate, src, 2)

	if err != nil {
		return nil, err
	}
	if frameRateLength != 2 {
		return nil, fmt.Errorf("broken FrameRate")
	}

	value := float64(uint8(frameRate.Bytes()[1])) + float64(uint8(frameRate.Bytes()[0]))/100.0

	result := &FrameRate{
		Value: value,
		data:  data,
	}

	return result, nil
}

type Uint8 struct {
	Value uint8
	data  *bytes.Buffer
}

func (u *Uint16) String() string {
	if u == nil {
		return "<nil>"
	}

	return fmt.Sprintf("Uint8{%d}", u.Value)
}

func (u *Uint8) Data() []byte {
	if u == nil || u.data == nil {
		return nil
	}

	return []byte{u.data.Bytes()[0]}
}

func (u *Uint8) Serialize() ([]byte, error) {
	if u == nil {
		return nil, fmt.Errorf("cannot serialize because Uint8 is nil")
	}

	return []byte{u.Value}, nil
}

func ReadUint8(src io.Reader) (*Uint8, error) {
	data := &bytes.Buffer{}

	dataLength, err := io.CopyN(data, input, 1)

	if err != nil {
		return nil, err
	}
	if dataLength != 1 {
		return nil, fmt.Errorf("broken Uint8")
	}

	result := &Uint8{
		Value: uint8(data.Bytes()[0]),
		data:  data,
	}

	return result, nil
}

type Uint16 struct {
	Value uint16
	data  *bytes.Buffer
}

func (u *Uint16) String() string {
	if u == nil {
		return "<nil>"
	}

	return fmt.Sprintf("Uint16{%d}", u.Value)
}

func (u *Uint16) Data() []byte {
	if u == nil || u.data == nil {
		return nil
	}

	var data []byte

	data = append(data, u.data.Bytes()...)

	return data
}

func (u *Uint16) Serialize() ([]byte, error) {
	if u == nil {
		return nil, fmt.Errorf("cannot serialize because Uint32 is nil")
	}

	data := &bytes.buffer{}

	if err := binary.Write(data, binary.LittleEndian, u.Value); err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}

func ReadUint16(src io.Reader) (*Uint16, error) {
	data := &bytes.Buffer{}

	dataLength, err := io.CopyN(data, input, 2)

	if err != nil {
		return nil, err
	}
	if dataLength != 1 {
		return nil, fmt.Errorf("broken Uint16")
	}

	var value uint16

	{
		w := &bytes.Buffer{}
		r := io.TeeReader(data, w)

		if err := binary.Read(r, binary.LittleEndian, &value); err != nil {
			return nil, fmt.Errorf("failed to read the buffer as uint16: %w", err)
		}

		data = w
	}

	result := &Uint16{
		Value: value,
		data:  data,
	}

	return result, nil
}

type Uint32 struct {
	Value uint32
	data  *bytes.Buffer
}

func (u *Uint32) String() string {
	if u == nil {
		return "<nil>"
	}

	return fmt.Sprintf("Uint32{%d}", u.Value)
}

func (u *Uint32) Data() []byte {
	if u == nil || u.data == nil {
		return nil
	}

	var data []byte

	data = append(data, u.data.Bytes()...)

	return data
}

func (u *Uint32) Serialize() ([]byte, error) {
	if u == nil {
		return nil, fmt.Errorf("cannot serialize because Uint32 is nil")
	}

	data := &bytes.buffer{}

	if err := binary.Write(data, binary.LittleEndian, u.Value); err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}

func ReadUint32(src io.Reader) (*Uint32, error) {
	data := &bytes.Buffer{}

	dataLength, err := io.CopyN(data, input, 4)

	if err != nil {
		return nil, err
	}
	if dataLength != 4 {
		return nil, fmt.Errorf("broken Uint32")
	}

	var value uint32

	{
		w := &bytes.Buffer{}
		r := io.TeeReader(data, w)

		if err := binary.Read(r, binary.LittleEndian, &value); err != nil {
			return nil, fmt.Errorf("failed to read the buffer as uint32: %w", err)
		}

		data = w
	}

	result := &Uint32{
		Value: value,
		data:  data,
	}

	return result, nil
}

type Rectangle struct {
	BitsPerField int
	MinX         uint32
	MaxX         uint32
	MinY         uint32
	MaxY         uint32
	data         *bytes.Buffer
}

func (r *Rectangle) String() string {
	if r == nil {
		return "<nil>"
	}

	return fmt.Sprintf("Rectangle{%d %d %d %d}", r.MinX, r.MaxX, r.MinY, r.MaxY)
}

func (r *Rectangle) Serialize() ([]byte, error) {
	if r == nil {
		return nil, fmt.Errorf("cannot serialize because Rectangle is nil")
	}

	var s string

	s += fmt.Sprintf("%05b", uint8(r.BitsPerField))
	s += fmt.Sprintf("%032b", r.MinX)[32-r.BitsPerField:]
	s += fmt.Sprintf("%032b", r.MaxX)[32-r.BitsPerField:]
	s += fmt.Sprintf("%032b", r.MinY)[32-r.BitsPerField:]
	s += fmt.Sprintf("%032b", r.MaxY)[32-r.BitsPerField:]

	zeroPadding := len(s) % 8

	if zeroPadding > 0 {
		zeroPadding = 8 - zeroPadding
	}
	for i := 0; i < zeroPadding; i++ {
		s += "0"
	}

	var data []byte

	for i := 0; i < len(s); i += 8 {
		i64, err := strconv.ParseInt(s[i:i+8], 2, 64)

		if err != nil {
			return nil, err
		}

		data = append(data, byte(i64))
	}

	return data, nil
}

func (r *Rectangle) Data() []byte {
	if r == nil || r.data == nil {
		return nil
	}

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
	if c == nil {
		return "<nil>"
	}
	if len(c.data.Bytes()) == 3 {
		return fmt.Sprintf("Color{0x%x 0x%x 0x%x}", c.Red, c.Green, c.Blue)
	}

	return fmt.Sprintf("Color{0x%x 0x%x 0x%x 0x%x}", c.Red, c.Green, c.Blue, c.Alpha)
}

func (c *Color) Serialize() ([]byte, error) {
	if c == nil {
		return nil, fmt.Errorf("cannot serialize because Color is nil")
	}
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
	if c == nil || c.data == nil {
		return nil
	}

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
