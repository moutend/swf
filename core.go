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

func (f *File) Bytes() []byte {
	if f == nil {
		return nil
	}

	var data []byte

	data = append(data, f.Signature.Bytes()...)
	data = append(data, f.Version.Bytes()...)
	data = append(data, f.FileSize.Bytes()...)
	data = append(data, f.Rectangle.Bytes()...)
	data = append(data, f.FrameRate.Bytes()...)
	data = append(data, f.FrameCount.Bytes()...)
	data = append(data, f.Contents.Bytes()...)

	return data
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

	contentsData, err := f.Contents.Serialize()

	if err != nil {
		return nil, err
	}

	var header []byte

	header = append(header, signatureData...)
	header = append(header, versionData...)
	header = append(header, fileSizeData...)

	var body []byte

	body = append(body, rectangleData...)
	body = append(body, frameRateData...)
	body = append(body, frameCountData...)
	body = append(body, contentsData...)

	if f.Signature.Value == SignatureCompressed {
		buffer := &bytes.Buffer{}
		compressed := zlib.NewWriter(buffer)

		if _, err := io.Copy(compressed, bytes.NewBuffer(body)); err != nil {
			return nil, err
		}

		defer compressed.Close()

		body = buffer.Bytes()
	}

	var result []byte

	result = append(result, header...)
	result = append(result, body...)

	return result, nil
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
	Bytes() []byte
	Serialize() ([]byte, error)
}

type ContentSlice []Content

func (c ContentSlice) String() string {
	return fmt.Sprintf("ContentSlice{%d items}", len(c))
}

func (c ContentSlice) Bytes() []byte {
	var data []byte

	for i := range c {
		data = append(data, c[i].Bytes()...)
	}

	return data
}

func (c ContentSlice) Serialize() ([]byte, error) {
	var data []byte

	for i := range c {
		contentData, err := c[i].Serialize()

		if err != nil {
			return nil, err
		}

		data = append(data, contentData...)
	}

	return data, nil
}

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
		content = ParseEnd(tag)
	case ShowFrameTagCode:
		content = ParseShowFrame(tag)
	case SetBackgroundColorTagCode:
		content, err = ParseSetBackgroundColor(src, tag)
	default:
		content, err = ParseUnknown(src, tag, extended)
	}

	return content, err
}

type End struct {
	Tag *Uint16
}

func (v *End) TagCode() TagCode {
	return EndTagCode
}

func (v *End) String() string {
	if v == nil {
		return "<nil"
	}

	return "End{}"
}

func (v *End) Bytes() []byte {
	if v == nil {
		return nil
	}

	var data []byte

	data = append(data, v.Tag.Bytes()...)

	return data
}

func (v *End) Serialize() ([]byte, error) {
	return v.Bytes(), nil
}

func ParseEnd(tag *Uint16) *End {
	return &End{Tag: tag}
}

type ShowFrame struct {
	Tag *Uint16
}

func (v *ShowFrame) TagCode() TagCode {
	return ShowFrameTagCode
}

func (v *ShowFrame) String() string {
	if v == nil {
		return "<nil"
	}

	return "ShowFrame{}"
}

func (v *ShowFrame) Bytes() []byte {
	if v == nil {
		return nil
	}

	var data []byte

	data = append(data, v.Tag.Bytes()...)

	return data
}

func (v *ShowFrame) Serialize() ([]byte, error) {
	return v.Bytes(), nil
}

func ParseShowFrame(tag *Uint16) *ShowFrame {
	return &ShowFrame{Tag: tag}
}

type SetBackgroundColor struct {
	Tag   *Uint16
	Color *Color
}

func (v *SetBackgroundColor) TagCode() TagCode {
	return SetBackgroundColorTagCode
}

func (v *SetBackgroundColor) String() string {
	if v == nil {
		return "<nil>"
	}

	return fmt.Sprintf("SetBackgroundColor{Color: %s}", v.Color)
}

func (v *SetBackgroundColor) Bytes() []byte {
	if v == nil {
		return nil
	}

	var data []byte

	if v.Tag != nil {
		data = append(data, v.Tag.Bytes()...)
	}
	if v.Color != nil {
		data = append(data, v.Color.Bytes()...)
	}

	return data
}

func (v *SetBackgroundColor) Serialize() ([]byte, error) {
	if v == nil {
		return nil, fmt.Errorf("cannot serialize because SetBackgroundColor is nil")
	}

	var data []byte

	tagData, err := v.Tag.Serialize()

	if err != nil {
		return nil, err
	}

	colorData, err := v.Color.Serialize()

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

func (v *Unknown) TagCode() TagCode {
	return UnknownTagCode
}

func (v *Unknown) String() string {
	if v == nil {
		return "<nil>"
	}

	return fmt.Sprintf("Unknown{%d bytes}", len(v.data.Bytes()))
}

func (v *Unknown) Bytes() []byte {
	if v == nil || v.data == nil {
		return nil
	}

	var data []byte

	if v.Tag != nil {
		data = append(data, v.Tag.Bytes()...)
	}
	if v.Extended != nil {
		data = append(data, v.Extended.Bytes()...)
	}

	data = append(data, v.data.Bytes()...)

	return data
}

func (v *Unknown) Serialize() ([]byte, error) {
	if v == nil {
		return nil, fmt.Errorf("cannot serialize because Unknown is nil")
	}

	var data []byte

	tagData, err := v.Tag.Serialize()

	if err != nil {
		return nil, err
	}

	extendedData, err := v.Extended.Serialize()

	if err != nil {
		return nil, err
	}

	data = append(data, tagData...)
	data = append(data, extendedData...)
	data = append(data, v.data.Bytes()...)

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
