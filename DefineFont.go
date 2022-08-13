package swf

import (
	"bytes"
	"fmt"
	"io"
)

type DefineFont struct {
	Tag      *Uint16
	Extended *Uint32
	data     *bytes.Buffer
}

func (v *DefineFont) TagCode() TagCode {
	return DefineFontTagCode
}

func (v *DefineFont) String() string {
	if v == nil {
		return "<nil>"
	}

	return fmt.Sprintf("DefineFont{%d bytes}", len(v.data.Bytes()))
}

func (v *DefineFont) Bytes() []byte {
	if v == nil {
		return nil
	}

	var data []byte

	if v.Tag != nil {
		data = append(data, v.Tag.Bytes()...)
	}
	if v.Extended != nil {
		data = append(data, v.Extended.Bytes()...)
	}
	if v.data != nil {
		data = append(data, v.data.Bytes()...)
	}

	return data
}

func (v *DefineFont) Serialize() ([]byte, error) {
	if v == nil {
		return nil, fmt.Errorf("cannot serialize because DefineFont is nil")
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

func ParseDefineFont(src io.Reader, tag *Uint16, extended *Uint32) (*DefineFont, error) {
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
		return nil, fmt.Errorf("broken DefineFont")
	}

	result := &DefineFont{
		Tag:      tag,
		Extended: extended,
		data:     data,
	}

	return result, nil
}
