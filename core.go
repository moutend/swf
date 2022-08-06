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

type Content struct {
	TagCode uint16
	Data    *bytes.Buffer
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

type Header struct {
	Signature  string
	Version    uint8
	FileSize   uint32
	Rect       []uint32
	FrameRate  float32
	FrameCount uint16

	rawSignature  *bytes.Buffer
	rawVersion    *bytes.Buffer
	rawFileSize   *bytes.Buffer
	rawRect       *bytes.Buffer
	rawFrameRate  *bytes.Buffer
	rawFrameCount *bytes.Buffer
}

func (h *Header) String() string {
	return fmt.Sprintf(
		"Header{Signature: %q, Version: %d, FileSize: %d, Rect: %+v, FrameRate: %.2f, FrameCount: %d}",
		h.Signature, h.Version, h.FileSize, h.Rect, h.FrameRate, h.FrameCount,
	)
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
		Signature:     signature.String(),
		Version:       uint8(version.Bytes()[0]),
		FileSize:      fileSizeUint32,
		Rect:          rectUint32Slice,
		FrameRate:     frameRateFloat32,
		FrameCount:    frameCountUint16,
		rawSignature:  signature,
		rawVersion:    version,
		rawFileSize:   fileSize,
		rawRect:       rect,
		rawFrameRate:  frameRate,
		rawFrameCount: frameCount,
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
	return nil, nil
}

func parseContent(input io.Reader) (*Content, error) {
	// Read the tag and its content length. (Fixed length, 2 byte, little endian)
	tagCode := &bytes.Buffer{}

	tagCodeLength, err := io.CopyN(tagCode, input, 2)

	if err != nil {
		return nil, err
	}
	if tagCodeLength != 2 {
		return nil, fmt.Errorf("broken tag code")
	}

	var tagCodeUint16 uint16

	if err := binary.Read(tagCode, binary.LittleEndian, &tagCodeUint16); err != nil {
		return nil, fmt.Errorf("failed to read tag code as uint16: %w", err)
	}

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

		length = int64(u32)
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
		TagCode: tagCodeUint16 >> 6,
		Data:    data,
	}

	return content, nil
}
