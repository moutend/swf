package swf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/moutend/go-bits"
)

const (
	SignatureUncompressed = `SWF`
	SignatureCompressed   = `SWC`
)

type Uint8 struct {
	Value uint8
	value uint8
}

func (u *Uint8) String() string {
	if u == nil {
		return "<nil>"
	}

	return fmt.Sprintf("Uint8{%d}", u.Value)
}

func (u *Uint8) Bytes() []byte {
	if u == nil {
		return nil
	}

	return []byte{u.value}
}

func (u *Uint8) Serialize() ([]byte, error) {
	if u == nil {
		return nil, nil
	}

	return []byte{u.Value}, nil
}

func ReadUint8(src io.Reader) (*Uint8, error) {
	data := &bytes.Buffer{}

	dataLength, err := io.CopyN(data, src, 1)

	if err != nil {
		return nil, fmt.Errorf("failed to read Uint8: %w", err)
	}
	if dataLength != 1 {
		return nil, fmt.Errorf("failed to read Uint8: length must be 1 but got %d", dataLength)
	}

	result := &Uint8{
		Value: uint8(data.Bytes()[0]),
		value: uint8(data.Bytes()[0]),
	}

	return result, nil
}

type Uint16 struct {
	Value uint16
	data  *bytes.Buffer
}

func (u *Uint16) Contains(value uint16) bool {
	return u.Value&value == value
}

func (u *Uint16) String() string {
	if u == nil {
		return "<nil>"
	}

	return fmt.Sprintf("Uint16{%d}", u.Value)
}

func (u *Uint16) Bytes() []byte {
	if u == nil || u.data == nil {
		return nil
	}

	var data []byte

	data = append(data, u.data.Bytes()...)

	return data
}

func (u *Uint16) Serialize() ([]byte, error) {
	if u == nil {
		return nil, nil
	}

	data := &bytes.Buffer{}

	if err := binary.Write(data, binary.LittleEndian, u.Value); err != nil {
		return nil, fmt.Errorf("failed to serialize Uint16: %w", err)
	}

	return data.Bytes(), nil
}

func ReadUint16(src io.Reader) (*Uint16, error) {
	data := &bytes.Buffer{}

	dataLength, err := io.CopyN(data, src, 2)

	if err != nil {
		return nil, fmt.Errorf("failed to read Uint16: %w", err)
	}
	if dataLength != 2 {
		return nil, fmt.Errorf("failed to read Uint16: length must be 2 but got %d", dataLength)
	}

	var value uint16

	{
		w := &bytes.Buffer{}
		r := io.TeeReader(data, w)

		if err := binary.Read(r, binary.LittleEndian, &value); err != nil {
			return nil, fmt.Errorf("failed to read Uint16: cannot read the buffer as uint16: %w", err)
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

func (u *Uint32) Bytes() []byte {
	if u == nil || u.data == nil {
		return nil
	}

	var data []byte

	data = append(data, u.data.Bytes()...)

	return data
}

func (u *Uint32) Serialize() ([]byte, error) {
	if u == nil {
		return nil, nil
	}

	data := &bytes.Buffer{}

	if err := binary.Write(data, binary.LittleEndian, u.Value); err != nil {
		return nil, fmt.Errorf("failed to serialize Uint32: %w", err)
	}

	return data.Bytes(), nil
}

func ReadUint32(src io.Reader) (*Uint32, error) {
	data := &bytes.Buffer{}

	dataLength, err := io.CopyN(data, src, 4)

	if err != nil {
		return nil, fmt.Errorf("failed to read Uint32: %w", err)
	}
	if dataLength != 4 {
		return nil, fmt.Errorf("failed to read Uint32: length must be 4 but got %d", dataLength)
	}

	var value uint32

	{
		w := &bytes.Buffer{}
		r := io.TeeReader(data, w)

		if err := binary.Read(r, binary.LittleEndian, &value); err != nil {
			return nil, fmt.Errorf("failed to read Uint32: cannot read the buffer as uint32: %w", err)
		}

		data = w
	}

	result := &Uint32{
		Value: value,
		data:  data,
	}

	return result, nil
}

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

func (s *Signature) Bytes() []byte {
	if s == nil || s.data == nil {
		return nil
	}

	var data []byte

	data = append(data, s.data.Bytes()...)

	return data
}

func (s *Signature) Serialize() ([]byte, error) {
	if s == nil {
		return nil, nil
	}
	switch s.Value {
	case SignatureUncompressed:
		return []byte(`FWS`), nil
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

func (f *FrameRate) Bytes() []byte {
	if f == nil || f.data == nil {
		return nil
	}

	var data []byte

	data = append(data, f.data.Bytes()...)

	return data
}

func (f *FrameRate) Serialize() ([]byte, error) {
	if f == nil {
		return nil, nil
	}

	a := uint8(f.Value)
	b := uint8((f.Value - math.Floor(f.Value)) * 100.0)

	return []byte{b, a}, nil
}

func ReadFrameRate(src io.Reader) (*FrameRate, error) {
	data := &bytes.Buffer{}

	dataLength, err := io.CopyN(data, src, 2)

	if err != nil {
		return nil, err
	}
	if dataLength != 2 {
		return nil, fmt.Errorf("broken FrameRate")
	}

	value := float64(uint8(data.Bytes()[1])) + float64(uint8(data.Bytes()[0]))/100.0

	result := &FrameRate{
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

func (r *Rectangle) Bytes() []byte {
	if r == nil || r.data == nil {
		return nil
	}

	var data []byte

	data = append(data, r.data.Bytes()...)

	return data
}

func (r *Rectangle) Serialize() ([]byte, error) {
	if r == nil {
		return nil, nil
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
		return fmt.Sprintf("RGB{%d, %d, %d}", c.Red, c.Green, c.Blue)
	}

	return fmt.Sprintf("RGBA{%d, %d, %d, %d}", c.Red, c.Green, c.Blue, c.Alpha)
}

func (c *Color) Bytes() []byte {
	if c == nil || c.data == nil {
		return nil
	}

	var data []byte

	data = append(data, c.data.Bytes()...)

	return data
}

func (c *Color) Serialize() ([]byte, error) {
	if c == nil {
		return nil, nil
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

type Matrix struct {
	// Scale
	HasScale     bool
	NumScaleBits uint8
	A            uint32
	D            uint32
	// Rotate/Skew
	HasRotate     bool
	NumRotateBits uint8
	B             uint32
	C             uint32
	// Translate (always present)
	NumTranslateBits uint8
	TX               uint32
	TY               uint32
}

func ReadMatrix(src io.Reader) (*Matrix, error) {
	matrix := &Matrix{}
	buffer := &bits.Buffer{}

	hasScale, err := buffer.Scan(src, 1)

	if err != nil {
		return nil, fmt.Errorf("failed to read Matrix.HasScale: %w", err)
	}
	if hasScale == 1 {
		numScaleBits, err := buffer.Scan(src, 5)

		if err != nil {
			return nil, fmt.Errorf("failed to read Matrix.NumScaleBits: %w", err)
		}

		a, err := buffer.Scan(src, int(numScaleBits))

		if err != nil {
			return nil, fmt.Errorf("failed to read Matrix.A: %w", err)
		}

		d, err := buffer.Scan(src, int(numScaleBits))

		if err != nil {
			return nil, fmt.Errorf("failed to read Matrix.D: %w", err)
		}

		matrix.HasScale = true
		matrix.NumScaleBits = uint8(numScaleBits)
		matrix.A = uint32(a)
		matrix.D = uint32(d)
	}

	hasRotate, err := buffer.Scan(src, 1)

	if err != nil {
		return nil, fmt.Errorf("failed to read Matrix.HasRotate: %w", err)
	}
	if hasRotate == 1 {
		numRotateBits, err := buffer.Scan(src, 5)

		if err != nil {
			return nil, fmt.Errorf("failed to read Matrix.NumRotateBits: %w", err)
		}

		b, err := buffer.Scan(src, int(numRotateBits))

		if err != nil {
			return nil, fmt.Errorf("failed to read Matrix.B: %w", err)
		}

		c, err := buffer.Scan(src, int(numRotateBits))

		if err != nil {
			return nil, fmt.Errorf("failed to read Matrix.C: %w", err)
		}

		matrix.HasRotate = true
		matrix.NumRotateBits = uint8(numRotateBits)
		matrix.B = uint32(b)
		matrix.C = uint32(c)
	}

	numTranslateBits, err := buffer.Scan(src, 5)

	if err != nil {
		return nil, fmt.Errorf("failed to read Matrix.NumTranslateBits: %w", err)
	}

	tx, err := buffer.Scan(src, int(numTranslateBits))

	if err != nil {
		return nil, fmt.Errorf("failed to read Matrix.TX: %w", err)
	}

	ty, err := buffer.Scan(src, int(numTranslateBits))

	if err != nil {
		return nil, fmt.Errorf("failed to read Matrix.TY: %w", err)
	}

	matrix.NumTranslateBits = uint8(numTranslateBits)
	matrix.TX = uint32(tx)
	matrix.TY = uint32(ty)

	return matrix, nil
}

type GradientRecord struct {
	Ratio *Uint8
	Color *Color
}

func ReadGradientRecord(src io.Reader, shapeVersion int) (*GradientRecord, error) {
	ratio, err := ReadUint8(src)

	if err != nil {
		return nil, fmt.Errorf("failed to read GradientRecord: %w", err)
	}

	var color *Color

	if shapeVersion >= 3 {
		color, err = ReadRGBA(src)
	} else {
		color, err = ReadRGB(src)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read GradientRecord: %w", err)
	}

	result := &GradientRecord{
		Ratio: ratio,
		Color: color,
	}

	return result, nil
}

type GradientFlags struct {
	NumRecords    uint8
	Spread        uint8
	Interporation uint8
}

func ReadGradientFlags(src io.Reader) (*GradientFlags, error) {
	flags, err := ReadUint8(src)

	if err != nil {
		return nil, fmt.Errorf("failed to read Gradient flags: %w", err)
	}

	spread := (flags.Value >> 6) & 0b11
	interporation := (flags.Value >> 4) & 0b11
	numRecords := flags.Value & 0b1111

	result := &GradientFlags{
		NumRecords:    numRecords,
		Spread:        spread,
		Interporation: interporation,
	}

	return result, nil
}

type Gradient struct {
	Matrix  *Matrix
	Flags   *GradientFlags
	Records []*GradientRecord
}

func ReadGradient(src io.Reader, shapeVersion int) (*Gradient, error) {
	matrix, err := ReadMatrix(src)

	if err != nil {
		return nil, fmt.Errorf("failed to read Gradient.Matrix: %w", err)
	}

	flags, err := ReadGradientFlags(src)

	if err != nil {
		return nil, fmt.Errorf("failed to read Gradient flags: %w", err)
	}

	records := make([]*GradientRecord, int(flags.NumRecords))

	for i := range records {
		record, err := ReadGradientRecord(src, shapeVersion)

		if err != nil {
			return nil, fmt.Errorf("failed to read Gradient.Records[%d]: %w", i, err)
		}

		records[i] = record
	}

	result := &Gradient{
		Matrix:  matrix,
		Flags:   flags,
		Records: records,
	}

	return result, nil
}

type MorphGradient struct {
	Flags *GradientFlags
	Start *Gradient
	End   *Gradient
}

func ReadMorphGradient(src io.Reader, shapeVersion int) (*MorphGradient, error) {
	startMatrix, err := ReadMatrix(src)

	if err != nil {
		return nil, fmt.Errorf("failed to read MorphGradient.Start.Matrix: %w", err)
	}

	endMatrix, err := ReadMatrix(src)

	if err != nil {
		return nil, fmt.Errorf("failed to read MorphGradient.End.Matrix: %w", err)
	}

	flags, err := ReadGradientFlags(src)

	if err != nil {
		return nil, fmt.Errorf("failed to read MorphGradient.Flags: %w", err)
	}

	start := &Gradient{Matrix: startMatrix}
	end := &Gradient{Matrix: endMatrix}

	startRecords := make([]*GradientRecord, int(flags.NumRecords))
	endRecords := make([]*GradientRecord, int(flags.NumRecords))

	for i := 0; i < int(flags.NumRecords); i++ {
		startRecord, err := ReadGradientRecord(src, shapeVersion)

		if err != nil {
			return nil, fmt.Errorf("failed to read MorphGradient.Start.Records[%d]: %w", i, err)
		}

		startRecords[i] = startRecord

		endRecord, err := ReadGradientRecord(src, shapeVersion)

		if err != nil {
			return nil, fmt.Errorf("failed to read MorphGradient.End.Records[%d]: %w", i, err)
		}

		endRecords[i] = endRecord
	}

	result := &MorphGradient{
		Flags: flags,
		Start: start,
		End:   end,
	}

	return result, nil
}

type FillStyle struct {
	Type       *Uint8
	Color      *Color
	Gradient   *Gradient
	FocalPoint *Uint16
	ID         *Uint16
	Matrix     *Matrix
}

func ReadFillStyle(src io.Reader, shapeVersion int) (*FillStyle, error) {
	fillStyleType, err := ReadUint8(src)

	if err != nil {
		return nil, fmt.Errorf("failed to read FillStyle.Type: %w", err)
	}

	result := &FillStyle{
		Type: fillStyleType,
	}

	switch fillStyleType.Value {
	case 0x00:
		var color *Color

		if shapeVersion >= 3 {
			color, err = ReadRGBA(src)
		} else {
			color, err = ReadRGB(src)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read FillStyle.Color: %w", err)
		}

		result.Color = color
	case 0x10, 0x12:
		gradient, err := ReadGradient(src, shapeVersion)

		if err != nil {
			return nil, fmt.Errorf("failed to read FillStyle.Gradient: %w", err)
		}

		result.Gradient = gradient
	case 0x13:
		gradient, err := ReadGradient(src, shapeVersion)

		if err != nil {
			return nil, fmt.Errorf("failed to read FillStyle.Gradient: %w", err)
		}

		focalPoint, err := ReadUint16(src)

		if err != nil {
			return nil, fmt.Errorf("failed to read FillStyle.FocalPoint: %w", err)
		}

		result.Gradient = gradient
		result.FocalPoint = focalPoint
	case 0x40, 0x41, 0x42, 0x43:
		id, err := ReadUint16(src)

		if err != nil {
			return nil, fmt.Errorf("failed to read FillStyle.ID: %w", err)
		}

		matrix, err := ReadMatrix(src)

		if err != nil {
			return nil, fmt.Errorf("failed to read FillStyle.Matrix: %w", err)
		}

		result.ID = id
		result.Matrix = matrix
	default:
		return nil, fmt.Errorf("failed to read FillStyle: invalid type: %d", fillStyleType.Value)
	}

	return result, nil
}

type MorphFillStyle struct {
	Type            *Uint8
	Start           *FillStyle
	End             *FillStyle
	MorphGradient   *MorphGradient
	StartFocalPoint *Uint16
	EndFocalPoint   *Uint16
	ID              *Uint16
}

func ReadMorphFillStyle(src io.Reader, shapeVersion int) (*MorphFillStyle, error) {
	fillStyleType, err := ReadUint8(src)

	if err != nil {
		return nil, fmt.Errorf("failed to read MorphFillStyle.Type: %w", err)
	}

	result := &MorphFillStyle{Type: fillStyleType}

	switch fillStyleType.Value {
	case 0x00:
		var startColor *Color

		if shapeVersion >= 3 {
			startColor, err = ReadRGBA(src)
		} else {
			startColor, err = ReadRGB(src)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read MorphFillStyle.Start.Color: %w", err)
		}

		var endColor *Color

		if shapeVersion >= 3 {
			endColor, err = ReadRGBA(src)
		} else {
			endColor, err = ReadRGB(src)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read MorphFillStyle.End.Color: %w", err)
		}

		result.Start = &FillStyle{Color: startColor}
		result.End = &FillStyle{Color: endColor}
	case 0x10, 0x12:
		morphGradient, err := ReadMorphGradient(src, shapeVersion)

		if err != nil {
			return nil, fmt.Errorf("failed to read MorphFillStyle.MorphGradient: %w", err)
		}

		result.MorphGradient = morphGradient
	case 0x13:
		morphGradient, err := ReadMorphGradient(src, shapeVersion)

		if err != nil {
			return nil, fmt.Errorf("failed to read MorphFillStyle.MorphGradient: %w", err)
		}

		startFocalPoint, err := ReadUint16(src)

		if err != nil {
			return nil, fmt.Errorf("failed to read MorphFillStyle.StartFocalPoint: %w", err)
		}

		endFocalPoint, err := ReadUint16(src)

		if err != nil {
			return nil, fmt.Errorf("failed to read MorphFillStyle.EndFocalPoint: %w", err)
		}

		result.MorphGradient = morphGradient
		result.StartFocalPoint = startFocalPoint
		result.EndFocalPoint = endFocalPoint
	case 0x40, 0x41, 0x42, 0x43:
		id, err := ReadUint16(src)

		if err != nil {
			return nil, fmt.Errorf("failed to read MorphFillStyle.ID: %w", err)
		}

		startMatrix, err := ReadMatrix(src)

		if err != nil {
			return nil, fmt.Errorf("failed to read MorphFillStyle.Start.Matrix: %w", err)
		}

		endMatrix, err := ReadMatrix(src)

		if err != nil {
			return nil, fmt.Errorf("failed to read MorphFillStyle.End.Matrix: %w", err)
		}

		result.ID = id
		result.Start = &FillStyle{Matrix: startMatrix}
		result.End = &FillStyle{Matrix: endMatrix}
	default:
		return nil, fmt.Errorf("failed to read MorphFillStyle: invalid type: %d", fillStyleType.Value)
	}

	return result, nil
}

const (
	// First byte.
	LineStyleFlagPixelHinting      = 1 << 0
	LineStyleFlagNoVerticalScale   = 1 << 1
	LineStyleFlagNoHorizontalScale = 1 << 2
	LineStyleFlagHasFill           = 1 << 3
	LineStyleFlagJoinStyle         = 0b11 << 4
	LineStyleFlagStartCapStyle     = 0b11 << 6
	// Second byte.
	LineStyleFlagEndCapStyle = 0b11 << 8
	LineStyleFlagNoClose     = 1 << 10
	// JoinStyle mask values.
	JoinStyleRound = 0b00 << 4
	JoinStyleBevel = 0b01 << 4
	JoinStyleMiter = 0b10 << 4
)

type LineStyle struct {
	Width     *Uint16
	Color     *Color
	Flags     *Uint16
	FillStyle *FillStyle
	Miter     *Uint16
}

func ReadLineStyle(src io.Reader, shapeVersion int) (*LineStyle, error) {
	width, err := ReadUint16(src)

	if err != nil {
		return nil, fmt.Errorf("failed to read LineStyle.Width: %w", err)
	}

	result := &LineStyle{Width: width}

	if shapeVersion < 4 {
		var color *Color

		if shapeVersion >= 3 {
			color, err = ReadRGBA(src)
		} else {
			color, err = ReadRGB(src)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read LineStyle.Color: %w", err)
		}

		result.Color = color

		return result, nil
	}

	flags, err := ReadUint16(src)

	if err != nil {
		return nil, fmt.Errorf("failed to read LineStyle.Flags: %w", err)
	}

	result.Flags = flags

	if flags.Contains(LineStyleFlagJoinStyle) {
		flags.Value -= LineStyleFlagJoinStyle
	}
	if flags.Contains(LineStyleFlagStartCapStyle) {
		flags.Value -= LineStyleFlagStartCapStyle
	}
	if flags.Contains(LineStyleFlagEndCapStyle) {
		flags.Value -= LineStyleFlagEndCapStyle
	}
	if flags.Value&LineStyleFlagJoinStyle == JoinStyleMiter {
		miter, err := ReadUint16(src)

		if err != nil {
			return nil, fmt.Errorf("failed to read LineStyle.Miter: %w", err)
		}

		result.Miter = miter
	}
	if flags.Contains(LineStyleFlagHasFill) {
		fillStyle, err := ReadFillStyle(src, shapeVersion)

		if err != nil {
			return nil, fmt.Errorf("failed to read LineStyle.FillStyle: %w", err)
		}

		result.FillStyle = fillStyle
	} else {
		color, err := ReadRGBA(src)

		if err != nil {
			return nil, fmt.Errorf("failed tor ead LineStyle.Color: %w", err)
		}

		result.FillStyle = &FillStyle{Color: color}
	}

	return result, nil
}

type ShapeStyles struct {
	NumFillStyles  *Uint8
	NumFillStyles2 *Uint16
	FillStyles     []*FillStyle
	NumLineStyles  *Uint8
	NumLineStyles2 *Uint16
	LineStyles     []*LineStyle
	NumBits        *Uint8
}

func ReadShapeStyles(src io.Reader, shapeVersion int) (*ShapeStyles, error) {
	numFillStyles, err := ReadUint8(src)

	if err != nil {
		return nil, fmt.Errorf("failed to read ShapeStyles.NumFillStyles: %w", err)
	}

	result := &ShapeStyles{NumFillStyles: numFillStyles}

	var numFillStylesInt int

	if numFillStyles.Value == 0xff && shapeVersion >= 2 {
		numFillStyles2, err := ReadUint16(src)

		if err != nil {
			return nil, fmt.Errorf("failed to read ShapeStyles.NumFillStyles2: %w", err)
		}

		result.NumFillStyles2 = numFillStyles2

		numFillStylesInt = int(numFillStyles2.Value)
	} else {
		numFillStylesInt = int(numFillStyles.Value)
	}

	fillStyles := make([]*FillStyle, numFillStylesInt)

	for i := range fillStyles {
		fillStyle, err := ReadFillStyle(src, shapeVersion)

		if err != nil {
			return nil, fmt.Errorf("failed to read ShapeStyles.FillStyles[%d]: %w", i, err)
		}

		fillStyles[i] = fillStyle
	}

	result.FillStyles = fillStyles

	numLineStyles, err := ReadUint8(src)

	if err != nil {
		return nil, fmt.Errorf("failed to read ShapeStyles.NumLineStyles: %w", err)
	}

	result.NumLineStyles = numLineStyles

	var numLineStylesInt int

	if numLineStyles.Value == 0xff && shapeVersion >= 2 {
		numLineStyles2, err := ReadUint16(src)

		if err != nil {
			return nil, fmt.Errorf("failed to read ShapeStyles.NumLineStyles2: %w", err)
		}

		result.NumLineStyles2 = numLineStyles2

		numLineStylesInt = int(numLineStyles2.Value)
	} else {
		numLineStylesInt = int(numLineStyles.Value)
	}

	lineStyles := make([]*LineStyle, numLineStylesInt)

	for i := range lineStyles {
		lineStyle, err := ReadLineStyle(src, shapeVersion)

		if err != nil {
			return nil, fmt.Errorf("failed to read ShapeStyles.LineStyles[%d]: %w", err)
		}

		lineStyles[i] = lineStyle
	}

	result.LineStyles = lineStyles

	numBits, err := ReadUint8(src)

	if err != nil {
		return nil, fmt.Errorf("failed to read ShapeStyles.NumBits: %w", err)
	}

	result.NumBits = numBits

	return result, nil
}

type ShapeContext struct {
	SWFVersion   int
	ShapeVersion int
	NumFillBits  uint8
	NumLineBits  uint8
}

type StyleChangeData struct {
	NumBitsValue    *uint64
	MoveToValue1    *uint64
	MoveToValue2    *uint64
	FillStyle0Value *uint64
	FillStyle1Value *uint64
	LineStyleValue  *uint64
	ShapeStyles     *ShapeStyles
}

type ShapeRecord struct {
	IsEdgeRecordValue   *uint64
	IsStraightEdgeValue *uint64
	NumBitsValue        *uint64
	IsAxisAlignedValue  *uint64
	IsVerticalValue     *uint64
	DeltaXValue         *uint64
	DeltaYValue         *uint64
	ControlDeltaXValue  *uint64
	ControlDeltaYValue  *uint64
	AnchorDeltaXValue   *uint64
	AnchorDeltaYValue   *uint64
	FlagsValue          *uint64
	StyleChangeData     *StyleChangeData
}

func ReadShapeRecord(src io.Reader, shapeContext *ShapeContext) (*ShapeRecord, error) {
	buffer := &bits.Buffer{}

	isEdgeRecordValue, err := buffer.Scan(src, 1)

	if err != nil {
		return nil, fmt.Errorf("failed to read ShapeRecord.IsEdgeRecordValue: %w", err)
	}

	isEdgeRecord := isEdgeRecordValue == 1
	result := &ShapeRecord{IsEdgeRecordValue: &isEdgeRecordValue}

	if isEdgeRecord {
		isStraightEdgeValue, err := buffer.Scan(src, 1)

		if err != nil {
			return nil, fmt.Errorf("failed to read ShapeRecord.IsStraightEdgeValue: %w", err)
		}

		isStraightEdge := isStraightEdgeValue == 1
		result.IsStraightEdgeValue = &isStraightEdgeValue

		numBitsValue, err := buffer.Scan(src, 4)

		if err != nil {
			return nil, fmt.Errorf("failed to read ShapeRecord.NumBitsValue: %w", err)
		}

		numBitsValue += 2
		result.NumBitsValue = &numBitsValue

		if isStraightEdge {
			// StraightEdge
			isAxisAlignedValue, err := buffer.Scan(src, 1)

			if err != nil {
				return nil, fmt.Errorf("failed to read ShapeRecord.IsAxisAlignedValue: %w", err)
			}

			isAxisAligned := isAxisAlignedValue == 1
			result.IsAxisAlignedValue = &isAxisAlignedValue

			isVerticalValue, err := buffer.Scan(src, 1)

			if err != nil {
				return nil, fmt.Errorf("failed to read ShapeRecord.IsVerticalValue: %w", err)
			}

			isVertical := isAxisAligned && isVerticalValue == 1
			result.IsVerticalValue = &isVerticalValue

			if !isAxisAligned || !isVertical {
				deltaXValue, err := buffer.Scan(src, int(numBitsValue))

				if err != nil {
					return nil, fmt.Errorf("failed to read ShapeRecord.DeltaXValue: %w", err)
				}

				result.DeltaXValue = &deltaXValue
			}
			if !isAxisAligned || isVertical {
				deltaYValue, err := buffer.Scan(src, int(numBitsValue))

				if err != nil {
					return nil, fmt.Errorf("failed to read ShapeRecord.DeltaYValue: %w", err)
				}

				result.DeltaYValue = &deltaYValue
			}
		} else {
			// CurvedEdge
			controlDeltaXValue, err := buffer.Scan(src, int(numBitsValue))

			if err != nil {
				return nil, fmt.Errorf("failed to read ShapeRecord.ControlDeltaXValue: %w", err)
			}

			result.ControlDeltaXValue = &controlDeltaXValue

			controlDeltaYValue, err := buffer.Scan(src, int(numBitsValue))

			if err != nil {
				return nil, fmt.Errorf("failed to read ShapeRecord.ControlDeltaYValue: %w", err)
			}

			result.ControlDeltaYValue = &controlDeltaYValue

			anchorDeltaXValue, err := buffer.Scan(src, int(numBitsValue))

			if err != nil {
				return nil, fmt.Errorf("failed to read ShapeRecord.AnchorDeltaXValue: %w", err)
			}

			result.AnchorDeltaXValue = &anchorDeltaXValue

			anchorDeltaYValue, err := buffer.Scan(src, int(numBitsValue))

			if err != nil {
				return nil, fmt.Errorf("failed to read ShapeRecord.AnchorDeltaYValue: %w", err)
			}

			result.AnchorDeltaYValue = &anchorDeltaYValue
		}
	} else {
		flagsValue, err := buffer.Scan(src, 5)

		if err != nil {
			return nil, fmt.Errorf("failed to read ShapeRecord.FlagsValue: %w", err)
		}

		result.FlagsValue = &flagsValue

		if flagsValue != 0 {
			// StyleChange
			newStyle := &StyleChangeData{}

			if (flagsValue & 0b1) != 0 {
				// move
				numBitsValue, err := buffer.Scan(src, 5)

				if err != nil {
					return nil, fmt.Errorf("failed to read StyleChangeData.NumBitsValue: %w", err)
				}

				newStyle.NumBitsValue = &numBitsValue

				moveToValue1, err := buffer.Scan(src, int(numBitsValue))

				if err != nil {
					return nil, fmt.Errorf("failed to read StyleChangeData.MoveToValue1: %w", err)
				}

				newStyle.MoveToValue1 = &moveToValue1

				moveToValue2, err := buffer.Scan(src, int(numBitsValue))

				if err != nil {
					return nil, fmt.Errorf("failed to read StyleChangeData.MoveToValue2: %w", err)
				}

				newStyle.MoveToValue2 = &moveToValue2
			}
			if (flagsValue & 0b10) != 0 {
				fillStyle0Value, err := buffer.Scan(src, int(shapeContext.NumFillBits))

				if err != nil {
					return nil, fmt.Errorf("failed to read StyleChangeData.FillStyle0Value: %w", err)
				}

				newStyle.FillStyle0Value = &fillStyle0Value
			}
			if (flagsValue & 0b100) != 0 {
				fillStyle1Value, err := buffer.Scan(src, int(shapeContext.NumFillBits))

				if err != nil {
					return nil, fmt.Errorf("failed to read StyleChangeData.FillStyle1Value: %w", err)
				}

				newStyle.FillStyle1Value = &fillStyle1Value
			}
			if (flagsValue & 0b1000) != 0 {
				lineStyleValue, err := buffer.Scan(src, int(shapeContext.NumLineBits))

				if err != nil {
					return nil, fmt.Errorf("failed to read StyleChangeData.LineStyleValue: %w", err)
				}

				newStyle.LineStyleValue = &lineStyleValue
			}
			if (flagsValue & 0b10000) != 0 {
				newShapeStyles, err := ReadShapeStyles(src, shapeContext.ShapeVersion)

				if err != nil {
					return nil, fmt.Errorf("failed to read ShapeRecord: %w", err)
				}

				newStyle.ShapeStyles = newShapeStyles
				shapeContext.NumFillBits = newShapeStyles.NumBits.Value >> 4
				shapeContext.NumLineBits = newShapeStyles.NumBits.Value & 0b1111
			}

			result.StyleChangeData = newStyle
		} else {
			return nil, nil
		}
	}

	return result, nil
}

type Shape struct {
	ID           *Uint16
	ShapeBounds  *Rectangle
	EdgeBounds   *Rectangle
	Flags        *Uint8
	ShapeStyles  *ShapeStyles
	ShapeRecords []*ShapeRecord
}

func ReadDefineShape(src io.Reader, swfVersion, shapeVersion int) (*Shape, error) {
	id, err := ReadUint16(src)

	if err != nil {
		return nil, fmt.Errorf("failed to read Shape.ID: %w", err)
	}

	shapeBounds, err := ReadRectangle(src)

	if err != nil {
		return nil, fmt.Errorf("failed to read Shape.ShapeBounds: %w", err)
	}

	result := &Shape{
		ID:          id,
		ShapeBounds: shapeBounds,
	}

	if shapeVersion >= 4 {
		edgeBounds, err := ReadRectangle(src)

		if err != nil {
			return nil, fmt.Errorf("failed to read Shape.EdgeBounds: %w", err)
		}

		flags, err := ReadUint8(src)

		if err != nil {
			return nil, fmt.Errorf("failed to read Shape.Flags: %w", err)
		}

		result.EdgeBounds = edgeBounds
		result.Flags = flags
	}

	shapeStyles, err := ReadShapeStyles(src, shapeVersion)

	if err != nil {
		return nil, fmt.Errorf("failed to read Shape.ShapeStyles: %w", err)
	}

	result.ShapeStyles = shapeStyles

	shapeContext := &ShapeContext{
		SWFVersion:   swfVersion,
		ShapeVersion: shapeVersion,
		NumFillBits:  shapeStyles.NumBits.Value >> 4,
		NumLineBits:  shapeStyles.NumBits.Value & 0b1111,
	}

	shapeRecords := []*ShapeRecord{}
	i := 0

	for {
		shapeRecord, err := ReadShapeRecord(src, shapeContext)

		if err != nil {
			return nil, fmt.Errorf("failed to read Shape.ShapeRecords[%d]: %w", i, err)
		}
		if shapeRecord == nil {
			break
		}

		shapeRecords = append(shapeRecords, shapeRecord)

		i += 1
	}

	result.ShapeRecords = shapeRecords

	return result, nil
}
