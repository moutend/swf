package swf

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
)

type File struct {
	Signature  *Signature
	Version    *Uint8
	FileSize   *Uint32
	Rectangle  *Rectangle
	FrameRate  *FrameRate
	FrameCount *Uint16
	Contents   ContentSlice
}

func (f *File) String() string {
	if f == nil {
		return "<nil>"
	}

	return fmt.Sprintf(
		"File{Compressed: %v, Version: %d, FileSize: %d, Width: %d, Height: %d, FrameRate: %.2f, FrameCount: %d}",
		f.Signature.Value == SignatureCompressed, f.Version.Value, f.FileSize.Value, f.Rectangle.MaxX/20, f.Rectangle.MaxY/20, f.FrameRate.Value, f.FrameCount.Value,
	)
}

func (f *File) Data() []byte {
	if f == nil {
		return nil
	}

	return nil
}

func (f *File) Serialize() ([]byte, error) {
	if f == nil {
		return nil, fmt.Errorf("cannot serialize because File is nil")
	}

	signatureData, err := f.Signature.Serialize()

	if err != nil {
		return nil, err
	}

	versionData, err := f.Version.Serialize()

	if err != nil {
		return nil, err
	}

	fileSizeData, err := f.FileSize.Serialize()

	if err != nil {
		return nil, err
	}

	rectangleData, err := f.Rectangle.Serialize()

	if err != nil {
		return nil, err
	}

	frameRateData, err := f.FrameRate.Serialize()

	if err != nil {
		return nil, err
	}

	frameCountData, err := f.FrameCount.Serialize()

	if err != nil {
		return nil, err
	}

	var data []byte

	data = append(data, signatureData...)
	data = append(data, versionData...)
	data = append(data, fileSizeData...)

	data = append(data, rectangleData...)
	data = append(data, frameRateData...)
	data = append(data, frameCountData...)

	return data, nil
}

func Parse(src io.Reader) (*File, error) {
	signature, err := ReadSignature(src)

	if err != nil {
		return nil, err
	}

	version, err := ReadUint8(src)

	if err != nil {
		return nil, err
	}

	fileSize, err := ReadUint32(src)

	if err != nil {
		return nil, err
	}
	if signature.Value == SignatureCompressed {
		reader, err := zlib.NewReader(src)

		if err != nil {
			return nil, err
		}

		defer reader.Close()

		content := &bytes.Buffer{}

		contentLength, err := io.Copy(content, reader)

		if err != nil {
			return nil, err
		}
		if contentLength != int64(fileSize.Value)-8 {
			return nil, fmt.Errorf("invalid content length: expected=%d, actual=%d", int64(fileSize.Value)-8, contentLength)
		}

		src = content
	}

	rectangle, err := ReadRectangle(src)

	if err != nil {
		return nil, err
	}

	frameRate, err := ReadFrameRate(src)

	if err != nil {
		return nil, err
	}

	frameCount, err := ReadUint16(src)

	if err != nil {
		return nil, err
	}

	contents, err := parseContents(src)

	if err != nil {
		return nil, err
	}

	file := &File{
		Signature:  signature,
		Version:    version,
		FileSize:   fileSize,
		Rectangle:  rectangle,
		FrameRate:  frameRate,
		FrameCount: frameCount,
		Contents:   contents,
	}

	return file, nil
}

type Content interface {
	TagCode() TagCode
	String() string
	Data() []byte
	Serialize() ([]byte, error)
}

type ContentSlice []Content

/*
func (c ContentSlice) Serialize() ([]byte, error) {
	var result []byte

	for i := range c {
		result = append(result, c[i].TagCodeBuffer.Bytes()...)

		if c[i].HasExtendedLength {
			length := len(c[i].DataBuffer.Bytes())

			if c[i].TagCode == DefineSprite {
				length = int(c[i].DefineSpriteLength)
			}

			buffer := &bytes.Buffer{}

			if err := binary.Write(buffer, binary.LittleEndian, uint32(length)); err != nil {
				return nil, err
			}

			result = append(result, buffer.Bytes()...)
		}

		result = append(result, c[i].DataBuffer.Bytes()...)
	}

	return result, nil
}
*/

func parseContents(src io.Reader) (ContentSlice, error) {
	var contents ContentSlice

	for {
		content, err := parseContent(src)

		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		contents = append(contents, content)
	}

	return contents, nil
}

func parseContent(src io.Reader) (Content, error) {
	tag, err := ReadUint16(src)

	if err == io.EOF {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	tagCode := TagCode(tag.Value >> 6)
	length := int64(tag.Value & 0b111111)

	var extended *Uint32

	if length == 0b111111 {
		extended, err = ReadUint32(src)

		if err != nil {
			return nil, err
		}
	}

	var content Content

	switch tagCode {
	case EndTagCode:
		content = &End{}
	case SetBackgroundColorTagCode:
		content, err = ParseSetBackgroundColor(src, tag)
	default:
		content, err = ParseUnknown(src, tag, extended)
	}
	if err != nil {
		return nil, err
	}

	return content, nil
}

type End struct {
	// No fields
}

func (e *End) TagCode() TagCode {
	return EndTagCode
}

func (e *End) String() string {
	if e == nil {
		return "<nil"
	}

	return "End{}"
}

func (e *End) Data() []byte {
	return []byte{0x00}
}

func (e *End) Serialize() ([]byte, error) {
	if e == nil {
		return nil, fmt.Errorf("cannot serialize because End is nil")
	}

	return []byte{0x00}, nil
}

type SetBackgroundColor struct {
	Tag   *Uint16
	Color *Color
}

func (s *SetBackgroundColor) TagCode() TagCode {
	return SetBackgroundColorTagCode
}

func (s *SetBackgroundColor) String() string {
	if s == nil {
		return "<nil>"
	}

	return fmt.Sprintf("SetBackgroundColor{Color: %s}", s.Color)
}

func (s *SetBackgroundColor) Data() []byte {
	if s == nil {
		return nil
	}

	var data []byte

	if s.Tag != nil {
		data = append(data, s.Tag.Data()...)
	}
	if s.Color != nil {
		data = append(data, s.Color.Data()...)
	}

	return data
}

func (s *SetBackgroundColor) Serialize() ([]byte, error) {
	if s == nil {
		return nil, fmt.Errorf("cannot serialize because SetBackgroundColor is nil")
	}

	var data []byte

	tagData, err := s.Tag.Serialize()

	if err != nil {
		return nil, err
	}

	colorData, err := s.Color.Serialize()

	if err != nil {
		return nil, err
	}

	data = append(data, tagData...)
	data = append(data, colorData...)

	return data, nil
}

func ParseSetBackgroundColor(src io.Reader, tag *Uint16) (*SetBackgroundColor, error) {
	if tag == nil {
		return nil, fmt.Errorf("cannot parse because tag is nil")
	}

	color, err := ReadRGB(src)

	if err != nil {
		return nil, err
	}

	result := &SetBackgroundColor{
		Tag:   tag,
		Color: color,
	}

	return result, nil
}

type Unknown struct {
	Tag      *Uint16
	Extended *Uint32
	data     *bytes.Buffer
}

func (u *Unknown) TagCode() TagCode {
	return UnknownTagCode
}

func (u *Unknown) String() string {
	if u == nil {
		return "<nil>"
	}

	return fmt.Sprintf("Unknown{%d bytes}", len(u.data.Bytes()))
}

func (u *Unknown) Data() []byte {
	if u == nil || u.data == nil {
		return nil
	}

	var data []byte

	if u.Tag != nil {
		data = append(data, u.Tag.Data()...)
	}
	if u.Extended != nil {
		data = append(data, u.Extended.Data()...)
	}

	data = append(data, u.data.Bytes()...)

	return data
}

func (u *Unknown) Serialize() ([]byte, error) {
	if u == nil {
		return nil, fmt.Errorf("cannot serialize because Unknown is nil")
	}

	var data []byte

	tagData, err := u.Tag.Serialize()

	if err != nil {
		return nil, err
	}

	extendedData, err := u.Extended.Serialize()

	if err != nil {
		return nil, err
	}

	data = append(data, tagData...)
	data = append(data, extendedData...)
	data = append(data, u.data.Bytes()...)

	return data, nil
}

func ParseUnknown(src io.Reader, tag *Uint16, extended *Uint32) (*Unknown, error) {
	if tag == nil {
		return nil, fmt.Errorf("cannot parse because tag is nil")
	}

	length := int64(tag.Value & 0b111111)

	if extended != nil {
		length = int64(extended.Value)
	}

	data := &bytes.Buffer{}

	dataLength, err := io.CopyN(data, src, length)

	if err != nil {
		return nil, err
	}
	if dataLength != length {
		return nil, fmt.Errorf("broken Unknown")
	}

	result := &Unknown{
		Tag:      tag,
		Extended: extended,
		data:     data,
	}

	return result, nil
}
