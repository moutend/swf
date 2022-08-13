package swf

import (
	"bytes"
	"fmt"
	"io"
)

type DefineButtonCxform struct {
	Tag      *Uint16
	Extended *Uint32
	data     *bytes.Buffer
}

func (v *DefineButtonCxform) TagCode() TagCode {
	return DefineButtonCxformTagCode
}

func (v *DefineButtonCxform) String() string {
	if v == nil {
		return "<nil>"
	}

	return fmt.Sprintf("DefineButtonCxform{%d bytes}", len(v.data.Bytes()))
}

func (v *DefineButtonCxform) Bytes() []byte {
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

func (v *DefineButtonCxform) Serialize() ([]byte, error) {
	if v == nil {
		return nil, fmt.Errorf("cannot serialize because DefineButtonCxform is nil")
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

func ParseDefineButtonCxform(src io.Reader, tag *Uint16, extended *Uint32) (*DefineButtonCxform, error) {
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
		return nil, fmt.Errorf("broken DefineButtonCxform")
	}

	result := &DefineButtonCxform{
		Tag:      tag,
		Extended: extended,
		data:     data,
	}

	return result, nil
}
