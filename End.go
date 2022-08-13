package swf

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
