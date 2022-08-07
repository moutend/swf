package swf

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
)

type File struct {
	Header   *Header
	Contents []*Content
}

func (f *File) Serialize() ([]byte, error) {
	var result []byte

	hb, err := f.Header.Serialize()

	if err != nil {
		return nil, err
	}

	result = append(result, hb...)

	for _, content := range f.Contents {
		cb, err := content.Serialize()

		if err != nil {
			return nil, err
		}

		result = append(result, cb...)
	}

	return result, nil
}

type Header struct {
	Signature  string
	Version    uint8
	FileSize   uint32
	Rect       []uint32
	FrameRate  float32
	FrameCount uint16

	SignatureBuffer  *bytes.Buffer
	VersionBuffer    *bytes.Buffer
	FileSizeBuffer   *bytes.Buffer
	RectBuffer       *bytes.Buffer
	FrameRateBuffer  *bytes.Buffer
	FrameCountBuffer *bytes.Buffer
}

func (h *Header) Serialize() ([]byte, error) {
	var result []byte

	result = append(result, h.SignatureBuffer.Bytes()...)
	result = append(result, h.VersionBuffer.Bytes()...)
	result = append(result, h.FileSizeBuffer.Bytes()...)
	result = append(result, h.RectBuffer.Bytes()...)
	result = append(result, h.FrameRateBuffer.Bytes()...)
	result = append(result, h.FrameCountBuffer.Bytes()...)

	return result, nil
}

func (h *Header) String() string {
	return fmt.Sprintf(
		"Header{Signature: %q, Version: %d, FileSize: %d, Rect: %+v, FrameRate: %.2f, FrameCount: %d}",
		h.Signature, h.Version, h.FileSize, h.Rect, h.FrameRate, h.FrameCount,
	)
}

type Content struct {
	TagCode           TagCode
	HasExtendedLength bool
	TagCodeBuffer     *bytes.Buffer
	DataBuffer        *bytes.Buffer
}

func (c *Content) Serialize() ([]byte, error) {
	var result []byte

	result = append(result, c.TagCodeBuffer.Bytes()...)

	if c.HasExtendedLength {
		size := len(c.DataBuffer.Bytes())
		buffer := &bytes.Buffer{}

		if err := binary.Write(buffer, binary.LittleEndian, uint32(size)); err != nil {
			return nil, err
		}

		result = append(result, buffer.Bytes()...)
	}

	result = append(result, c.DataBuffer.Bytes()...)

	return result, nil
}

func (c *Content) String() string {
	return fmt.Sprintf("Content{TagCode: %s, Data: %d bytes}", c.TagCode, len(c.DataBuffer.Bytes()))
}

func Parse(name string) (*File, error) {
	input, err := os.Open(name)

	if err != nil {
		return nil, fmt.Errorf("swf: %w", err)
	}

	defer input.Close()

	file, err := parseFile(input)

	if err != nil {
		return nil, fmt.Errorf("swf: failed to parse file: %w", err)
	}

	return file, nil
}

func parseFile(input io.Reader) (*File, error) {
	// Read the signature block. (Fixed length, 3 byte, little endian)
	signature := &bytes.Buffer{}

	signatureLength, err := io.CopyN(signature, input, 3)

	if err != nil {
		return nil, err
	}
	if signatureLength != 3 {
		return nil, fmt.Errorf("broken signature")
	}
	if signature.String() != `FWS` && signature.String() != `CWS` {
		return nil, fmt.Errorf("unexpected signature found: %s", signature.String())
	}

	// Read the version block. (Fixed length, 1 byte)
	version := &bytes.Buffer{}

	versionLength, err := io.CopyN(version, input, 1)

	if err != nil {
		return nil, err
	}
	if versionLength != 1 {
		return nil, fmt.Errorf("broken version")
	}

	// Read the file size block. (Fixed length, 4 byte, little endian)
	fileSize := &bytes.Buffer{}

	fileSizeLength, err := io.CopyN(fileSize, input, 4)

	if err != nil {
		return nil, err
	}
	if fileSizeLength != 4 {
		return nil, fmt.Errorf("broken file size")
	}

	var fileSizeUint32 uint32

	{
		w := &bytes.Buffer{}
		r := io.TeeReader(fileSize, w)

		if err := binary.Read(r, binary.LittleEndian, &fileSizeUint32); err != nil {
			return nil, fmt.Errorf("failed to read file size as uint32: %w", err)
		}

		fileSize = w
	}

	if signature.String() == `CWS` {
		reader, err := zlib.NewReader(input)

		if err != nil {
			return nil, err
		}

		defer reader.Close()

		content := &bytes.Buffer{}

		contentLength, err := io.Copy(content, reader)

		if err != nil {
			return nil, err
		}
		if contentLength != int64(fileSizeUint32-8) {
			return nil, fmt.Errorf("broken content length")
		}

		input = content
	}

	// Read the RECT block. (Variable length, big endian)
	rect := &bytes.Buffer{}

	rectLength, err := io.CopyN(rect, input, 1)

	if err != nil {
		return nil, err
	}
	if rectLength != 1 {
		return nil, fmt.Errorf("broken RECT")
	}

	// Read the first 5 bits.
	bits := fmt.Sprintf("%08b", rect.Bytes()[0])

	// The first 5 bits represents the each field length.
	bitsPerField, err := strconv.ParseInt(bits[:5], 2, 64)

	if err != nil {
		return nil, fmt.Errorf("failed to read first 5 bits from RECT: %w", err)
	}
	if bitsPerField <= 0 {
		return nil, fmt.Errorf("unexpected RECT bits per field: %w", bitsPerField)
	}

	var rectBits int64 = bitsPerField*4 - 3
	var requiredBits int64

	for {
		if requiredBits > rectBits {
			break
		}

		requiredBits += 8
	}

	rectLength, err = io.CopyN(rect, input, requiredBits/8)

	if err != nil {
		return nil, err
	}
	if rectLength != requiredBits/8 {
		return nil, fmt.Errorf("broken RECT")
	}

	rectUint32Slice, err := parseHeaderRect(rect.Bytes(), int(bitsPerField))

	if err != nil {
		return nil, err
	}

	// Read the frame rate. (Fixed length, 2 byte)
	frameRate := &bytes.Buffer{}

	frameRateLength, err := io.CopyN(frameRate, input, 2)

	if err != nil {
		return nil, err
	}
	if frameRateLength != 2 {
		return nil, fmt.Errorf("broken frame rate")
	}

	frameRateFloat32 := float32(uint8(frameRate.Bytes()[1])) + float32(uint8(frameRate.Bytes()[0]))/100.0

	// Read the frame count. (Fixed length, 2 byte, little endian)
	frameCount := &bytes.Buffer{}

	frameCountLength, err := io.CopyN(frameCount, input, 2)

	if err != nil {
		return nil, err
	}
	if frameCountLength != 2 {
		return nil, fmt.Errorf("broken frame count")
	}

	var frameCountUint16 uint16

	{
		w := &bytes.Buffer{}
		r := io.TeeReader(frameCount, w)

		if err := binary.Read(r, binary.LittleEndian, &frameCountUint16); err != nil {
			return nil, fmt.Errorf("failed to read frame count as uint16: %w", err)
		}

		frameCount = w
	}

	contents, err := parseContents(input)

	if err != nil {
		return nil, err
	}

	header := &Header{
		Signature:        signature.String(),
		Version:          uint8(version.Bytes()[0]),
		FileSize:         fileSizeUint32,
		Rect:             rectUint32Slice,
		FrameRate:        frameRateFloat32,
		FrameCount:       frameCountUint16,
		SignatureBuffer:  signature,
		VersionBuffer:    version,
		FileSizeBuffer:   fileSize,
		RectBuffer:       rect,
		FrameRateBuffer:  frameRate,
		FrameCountBuffer: frameCount,
	}

	file := &File{
		Header:   header,
		Contents: contents,
	}

	return file, nil
}

func parseHeaderRect(input []byte, bitsPerField int) ([]uint32, error) {
	var s string

	for i := range input {
		s += fmt.Sprintf("%08b", input[i])
	}

	result := make([]uint32, 4)

	// Ignore first 5 bits.
	start := 5

	for i := 0; i < 4; i++ {
		n, err := strconv.ParseInt(s[start:start+bitsPerField], 2, 64)

		if err != nil {
			return nil, err
		}

		// Convert twip to px. // (20 twip = 1 px)
		result[i] = uint32(n) / 20

		start += bitsPerField
	}

	return result, nil
}

func parseContents(input io.Reader) ([]*Content, error) {
	contents := []*Content{}

	for {
		content, err := parseContent(input)

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

func parseContent(input io.Reader) (*Content, error) {
	// Read the tag and its content length. (Fixed length, 2 byte, little endian)
	tagCode := &bytes.Buffer{}

	tagCodeLength, err := io.CopyN(tagCode, input, 2)

	if err == io.EOF {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if tagCodeLength != 2 {
		return nil, fmt.Errorf("broken tag code")
	}

	var tagCodeUint16 uint16

	{
		w := &bytes.Buffer{}
		r := io.TeeReader(tagCode, w)

		if err := binary.Read(r, binary.LittleEndian, &tagCodeUint16); err != nil {
			return nil, fmt.Errorf("failed to read tag code as uint16: %w", err)
		}

		tagCode = w
	}

	tagCodeValue := TagCode(tagCodeUint16 >> 6)

	hasExtendedLength := false
	length := int64(0b111111 & tagCodeUint16)

	if length == 0b111111 {
		// Read the extended length. (Fixed size, 4 byte, little endian)
		extended := &bytes.Buffer{}

		extendedLength, err := io.CopyN(extended, input, 4)

		if err != nil {
			return nil, err
		}
		if extendedLength != 4 {
			return nil, fmt.Errorf("broken extended tag code length")
		}

		var u32 uint32

		if err := binary.Read(extended, binary.LittleEndian, &u32); err != nil {
			return nil, fmt.Errorf("failed to read extended tag code length as uint32: %w", err)
		}

		hasExtendedLength = true
		length = int64(u32)
	}
	switch tagCodeValue {
	case DefineSprite:
		length = 4
	default:
		// do nothing
	}

	data := &bytes.Buffer{}

	dataLength, err := io.CopyN(data, input, length)

	if err != nil {
		return nil, err
	}
	if dataLength != length {
		return nil, fmt.Errorf("broken data")
	}

	content := &Content{
		TagCode:           tagCodeValue,
		HasExtendedLength: hasExtendedLength,
		TagCodeBuffer:     tagCode,
		DataBuffer:        data,
	}

	return content, nil
}
