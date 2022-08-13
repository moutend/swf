package swf

import (
	"bytes"
	"fmt"
	"io"
)

type DefineFontInfo struct {
	Tag      *Uint16
	Extended *Uint32
	data     *bytes.Buffer
}

func (v *DefineFontInfo) TagCode() TagCode {
	return DefineFontInfoTagCode
}

func (v *DefineFontInfo) String() string {
	if v == nil {
		return "<nil>"
	}

	return fmt.Sprintf("DefineFontInfo{%d bytes}", len(v.data.Bytes()))
}

func (v *DefineFontInfo) Bytes() []byte {
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

func (v *DefineFontInfo) Serialize() ([]byte, error) {
	if v == nil {
		return nil, fmt.Errorf("cannot serialize because DefineFontInfo is nil")
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

func ParseDefineFontInfo(src io.Reader, tag *Uint16, extended *Uint32) (*DefineFontInfo, error) {
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
		return nil, fmt.Errorf("broken DefineFontInfo")
	}

	result := &DefineFontInfo{
		Tag:      tag,
		Extended: extended,
		data:     data,
	}

	return result, nil
}
