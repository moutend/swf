package swf

import (
	"fmt"
	"io"
)

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
