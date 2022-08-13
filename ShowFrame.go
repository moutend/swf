package swf

type ShowFrame struct {
	Tag *Uint16
}

func (v *ShowFrame) TagCode() TagCode {
	return ShowFrameTagCode
}

func (v *ShowFrame) String() string {
	if v == nil {
		return "<nil"
	}

	return "ShowFrame{}"
}

func (v *ShowFrame) Bytes() []byte {
	if v == nil {
		return nil
	}

	var data []byte

	data = append(data, v.Tag.Bytes()...)

	return data
}

func (v *ShowFrame) Serialize() ([]byte, error) {
	return v.Bytes(), nil
}

func ParseShowFrame(tag *Uint16) *ShowFrame {
	return &ShowFrame{Tag: tag}
}
