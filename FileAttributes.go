package swf

import (
	"fmt"
	"io"
)

type FileAttributes struct {
	Tag      *Uint16
	Extended *Uint32
	Flags    *Uint32
}

func (v *FileAttributes) TagCode() TagCode {
	return FileAttributesTagCode
}

func (v *FileAttributes) String() string {
	if v == nil {
		return "<nil>"
	}

	return fmt.Sprintf("FileAttributes{}")
}

func (v *FileAttributes) Bytes() []byte {
	if v == nil {
		return nil
	}

	var data []byte

	if v.Tag != nil {
		data = append(data, v.Tag.Bytes()...)
	}
	if v.Flags != nil {
		data = append(data, v.Flags.Bytes()...)
	}

	return data
}

func (v *FileAttributes) Serialize() ([]byte, error) {
	if v == nil {
		return nil, fmt.Errorf("failed to serialize: FileAttributes is nil")
	}

	var data []byte

	tagData, err := v.Tag.Serialize()

	if err != nil {
		return nil, fmt.Errorf("failed to serialize FileAttributes.Tag: %w", err)
	}

	flagsData, err := v.Flags.Serialize()

	if err != nil {
		return nil, fmt.Errorf("failed to serialize FileAttributes.Flags: %w", err)
	}

	data = append(data, tagData...)
	data = append(data, flagsData...)

	return data, nil
}

func ParseFileAttributes(src io.Reader, tag *Uint16) (*FileAttributes, error) {
	if tag == nil {
		return nil, fmt.Errorf("failed to parse FileAttributes.Tag: tag is nil")
	}

	length := int64(tag.Value & 0b111111)

	if length != 4 {
		return nil, fmt.Errorf("failed to parse FileAttributes.Tag: content length must be 4 bytes")
	}

	flags, err := ReadUint32(src)

	if err != nil {
		return nil, fmt.Errorf("failed to parse FileAttributes.Flags: %w", err)
	}

	result := &FileAttributes{
		Tag:   tag,
		Flags: flags,
	}

	return result, nil
}
