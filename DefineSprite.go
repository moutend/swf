package swf

import (
	"fmt"
	"io"
)

type DefineSprite struct {
	Tag       *Uint16
	Extended  *Uint32
	ID        *Uint16
	NumFrames *Uint16
}

func (v *DefineSprite) TagCode() TagCode {
	return DefineSpriteTagCode
}

func (v *DefineSprite) String() string {
	if v == nil {
		return "<nil>"
	}

	return fmt.Sprintf("DefineSprite{ID: %d, NumFrames: %d}", v.ID.Value, v.NumFrames.Value)
}

func (v *DefineSprite) Bytes() []byte {
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
	if v.ID != nil {
		data = append(data, v.ID.Bytes()...)
	}
	if v.NumFrames != nil {
		data = append(data, v.NumFrames.Bytes()...)
	}

	return data
}

func (v *DefineSprite) Serialize() ([]byte, error) {
	if v == nil {
		return nil, fmt.Errorf("cannot serialize because DefineSprite is nil")
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

	idData, err := v.ID.Serialize()

	if err != nil {
		return nil, err
	}

	numFramesData, err := v.NumFrames.Serialize()

	if err != nil {
		return nil, err
	}

	data = append(data, tagData...)
	data = append(data, extendedData...)
	data = append(data, idData...)
	data = append(data, numFramesData...)

	return data, nil
}

func ParseDefineSprite(src io.Reader, tag *Uint16, extended *Uint32) (*DefineSprite, error) {
	if tag == nil {
		return nil, fmt.Errorf("cannot parse because tag is nil")
	}
	if extended == nil {
		return nil, fmt.Errorf("cannot parse because extended is nil")
	}

	id, err := ReadUint16(src)

	if err != nil {
		return nil, err
	}

	numFrames, err := ReadUint16(src)

	if err != nil {
		return nil, err
	}

	result := &DefineSprite{
		Tag:       tag,
		Extended:  extended,
		ID:        id,
		NumFrames: numFrames,
	}

	return result, nil
}
