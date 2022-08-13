package swf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
)

type Content struct {
	TagCode            TagCode
	HasExtendedLength  bool
	DefineSpriteLength int64
	TagCodeBuffer      *bytes.Buffer
	DataBuffer         *bytes.Buffer
}

type ContentSlice []*Content

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

func (c *Content) String() string {
	return fmt.Sprintf("Content{TagCode: %s, Data: %d bytes}", c.TagCode, len(c.DataBuffer.Bytes()))
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
	var defineSpriteLength int64

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
		defineSpriteLength = length
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
		TagCode:            tagCodeValue,
		HasExtendedLength:  hasExtendedLength,
		DefineSpriteLength: defineSpriteLength,
		TagCodeBuffer:      tagCode,
		DataBuffer:         data,
	}

	return content, nil
}
