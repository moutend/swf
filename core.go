package swf

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
)

type File struct {
	Signature  *Signature
	Version    *Uint8
	FileSize   *Uint32
	Rectangle  *Rectangle
	FrameRate  *FrameRate
	FrameCount *Uint16
	Contents   ContentSlice
}

func (f *File) String() string {
	if f == nil {
		return "<nil>"
	}

	return fmt.Sprintf(
		"File{Compressed: %v, Version: %d, FileSize: %d, Width: %d, Height: %d, FrameRate: %.2f, FrameCount: %d}",
		f.Signature.Value == SignatureCompressed, f.Version.Value, f.FileSize.Value, f.Rectangle.MaxX/20, f.Rectangle.MaxY/20, f.FrameRate.Value, f.FrameCount.Value,
	)
}

func (f *File) Bytes() []byte {
	if f == nil {
		return nil
	}

	var data []byte

	data = append(data, f.Signature.Bytes()...)
	data = append(data, f.Version.Bytes()...)
	data = append(data, f.FileSize.Bytes()...)
	data = append(data, f.Rectangle.Bytes()...)
	data = append(data, f.FrameRate.Bytes()...)
	data = append(data, f.FrameCount.Bytes()...)
	data = append(data, f.Contents.Bytes()...)

	return data
}

func (f *File) Serialize() ([]byte, error) {
	if f == nil {
		return nil, fmt.Errorf("failed to serialize: File is nil")
	}

	signatureData, err := f.Signature.Serialize()

	if err != nil {
		return nil, fmt.Errorf("failed to serialize File.Signature: %w", err)
	}

	versionData, err := f.Version.Serialize()

	if err != nil {
		return nil, fmt.Errorf("failed to serialize File.Version: %w", err)
	}

	rectangleData, err := f.Rectangle.Serialize()

	if err != nil {
		return nil, fmt.Errorf("failed to serialize File.Rectangle: %w", err)
	}

	frameRateData, err := f.FrameRate.Serialize()

	if err != nil {
		return nil, fmt.Errorf("failed to serialize File.FrameRate: %w", err)
	}

	frameCountData, err := f.FrameCount.Serialize()

	if err != nil {
		return nil, fmt.Errorf("failed to serialize File.FrameCount: %w", err)
	}

	contentsData, err := f.Contents.Serialize()

	if err != nil {
		return nil, fmt.Errorf("failed to serialize File.Contents: %w", err)
	}

	var body []byte

	body = append(body, rectangleData...)
	body = append(body, frameRateData...)
	body = append(body, frameCountData...)
	body = append(body, contentsData...)

	if f.Signature.Value == SignatureCompressed {
		buffer := &bytes.Buffer{}
		compressed := zlib.NewWriter(buffer)

		if _, err := io.Copy(compressed, bytes.NewBuffer(body)); err != nil {
			return nil, err
		}

		defer compressed.Close()

		body = buffer.Bytes()
	}

	fileSize := &Uint32{Value: uint32(len(body) + 8)}

	fileSizeData, err := fileSize.Serialize()

	if err != nil {
		return nil, err
	}

	var header []byte

	header = append(header, signatureData...)
	header = append(header, versionData...)
	header = append(header, fileSizeData...)

	var result []byte

	result = append(result, header...)
	result = append(result, body...)

	return result, nil
}

func Parse(src io.Reader) (*File, error) {
	signature, err := ReadSignature(src)

	if err != nil {
		return nil, err
	}

	version, err := ReadUint8(src)

	if err != nil {
		return nil, err
	}

	fileSize, err := ReadUint32(src)

	if err != nil {
		return nil, err
	}
	if signature.Value == SignatureCompressed {
		reader, err := zlib.NewReader(src)

		if err != nil {
			return nil, err
		}

		defer reader.Close()

		content := &bytes.Buffer{}

		contentLength, err := io.Copy(content, reader)

		if err != nil {
			return nil, err
		}
		if contentLength != int64(fileSize.Value)-8 {
			return nil, fmt.Errorf("invalid content length: expected=%d, actual=%d", int64(fileSize.Value)-8, contentLength)
		}

		src = content
	}

	rectangle, err := ReadRectangle(src)

	if err != nil {
		return nil, err
	}

	frameRate, err := ReadFrameRate(src)

	if err != nil {
		return nil, err
	}

	frameCount, err := ReadUint16(src)

	if err != nil {
		return nil, err
	}

	contents, err := parseContents(src)

	if err != nil {
		return nil, err
	}

	file := &File{
		Signature:  signature,
		Version:    version,
		FileSize:   fileSize,
		Rectangle:  rectangle,
		FrameRate:  frameRate,
		FrameCount: frameCount,
		Contents:   contents,
	}

	return file, nil
}

type Content interface {
	TagCode() TagCode
	String() string
	Bytes() []byte
	Serialize() ([]byte, error)
}

type ContentSlice []Content

func (c ContentSlice) String() string {
	return fmt.Sprintf("ContentSlice{%d items}", len(c))
}

func (c ContentSlice) Bytes() []byte {
	var data []byte

	for i := range c {
		data = append(data, c[i].Bytes()...)
	}

	return data
}

func (c ContentSlice) Serialize() ([]byte, error) {
	var data []byte

	for i := range c {
		contentData, err := c[i].Serialize()

		if err != nil {
			return nil, fmt.Errorf("failed to serialize %s (File.Contents[%d]): %w", c[i].TagCode(), i, err)
		}

		data = append(data, contentData...)
	}

	return data, nil
}

func parseContents(src io.Reader) (ContentSlice, error) {
	var contents ContentSlice

	for {
		content, err := parseContent(src)

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

func parseContent(src io.Reader) (Content, error) {
	tag, err := ReadUint16(src)

	if err != nil {
		return nil, err
	}

	tagCode := TagCode(tag.Value >> 6)
	length := int64(tag.Value & 0b111111)

	var extended *Uint32

	if length == 0b111111 {
		extended, err = ReadUint32(src)

		if err != nil {
			return nil, err
		}
	}

	var content Content

	switch tagCode {
	case EndTagCode:
		content = ParseEnd(tag)
	case ShowFrameTagCode:
		content = ParseShowFrame(tag)
	case DefineShapeTagCode:
		content, err = ParseDefineShape(src, tag, extended)
	case PlaceObjectTagCode:
		content, err = ParsePlaceObject(src, tag, extended)
	case RemoveObjectTagCode:
		content, err = ParseRemoveObject(src, tag, extended)
	case DefineBitsTagCode:
		content, err = ParseDefineBits(src, tag, extended)
	case DefineButtonTagCode:
		content, err = ParseDefineButton(src, tag, extended)
	case JpegTablesTagCode:
		content, err = ParseJpegTables(src, tag, extended)
	case SetBackgroundColorTagCode:
		content, err = ParseSetBackgroundColor(src, tag)
	case DefineFontTagCode:
		content, err = ParseDefineFont(src, tag, extended)
	case DefineTextTagCode:
		content, err = ParseDefineText(src, tag, extended)
	case DoActionTagCode:
		content, err = ParseDoAction(src, tag, extended)
	case DefineFontInfoTagCode:
		content, err = ParseDefineFontInfo(src, tag, extended)
	case DefineSoundTagCode:
		content, err = ParseDefineSound(src, tag, extended)
	case StartSoundTagCode:
		content, err = ParseStartSound(src, tag, extended)
	case DefineButtonSoundTagCode:
		content, err = ParseDefineButtonSound(src, tag, extended)
	case SoundStreamHeadTagCode:
		content, err = ParseSoundStreamHead(src, tag, extended)
	case SoundStreamBlockTagCode:
		content, err = ParseSoundStreamBlock(src, tag, extended)
	case DefineBitsLosslessTagCode:
		content, err = ParseDefineBitsLossless(src, tag, extended)
	case DefineBitsJpeg2TagCode:
		content, err = ParseDefineBitsJpeg2(src, tag, extended)
	case DefineShape2TagCode:
		content, err = ParseDefineShape2(src, tag, extended)
	case DefineButtonCxformTagCode:
		content, err = ParseDefineButtonCxform(src, tag, extended)
	case ProtectTagCode:
		content, err = ParseProtect(src, tag, extended)
	case PlaceObject2TagCode:
		content, err = ParsePlaceObject2(src, tag, extended)
	case RemoveObject2TagCode:
		content, err = ParseRemoveObject2(src, tag, extended)
	case DefineShape3TagCode:
		content, err = ParseDefineShape3(src, tag, extended)
	case DefineText2TagCode:
		content, err = ParseDefineText2(src, tag, extended)
	case DefineButton2TagCode:
		content, err = ParseDefineButton2(src, tag, extended)
	case DefineBitsJpeg3TagCode:
		content, err = ParseDefineBitsJpeg3(src, tag, extended)
	case DefineBitsLossless2TagCode:
		content, err = ParseDefineBitsLossless2(src, tag, extended)
	case DefineEditTextTagCode:
		content, err = ParseDefineEditText(src, tag, extended)
	case DefineSpriteTagCode:
		content, err = ParseDefineSprite(src, tag, extended)
	case NameCharacterTagCode:
		content, err = ParseNameCharacter(src, tag, extended)
	case ProductInfoTagCode:
		content, err = ParseProductInfo(src, tag, extended)
	case FrameLabelTagCode:
		content, err = ParseFrameLabel(src, tag, extended)
	case SoundStreamHead2TagCode:
		content, err = ParseSoundStreamHead2(src, tag, extended)
	case DefineMorphShapeTagCode:
		content, err = ParseDefineMorphShape(src, tag, extended)
	case DefineFont2TagCode:
		content, err = ParseDefineFont2(src, tag, extended)
	case ExportAssetsTagCode:
		content, err = ParseExportAssets(src, tag, extended)
	case ImportAssetsTagCode:
		content, err = ParseImportAssets(src, tag, extended)
	case EnableDebuggerTagCode:
		content, err = ParseEnableDebugger(src, tag, extended)
	case DoInitActionTagCode:
		content, err = ParseDoInitAction(src, tag, extended)
	case DefineVideoStreamTagCode:
		content, err = ParseDefineVideoStream(src, tag, extended)
	case VideoFrameTagCode:
		content, err = ParseVideoFrame(src, tag, extended)
	case DefineFontInfo2TagCode:
		content, err = ParseDefineFontInfo2(src, tag, extended)
	case DebugIdTagCode:
		content, err = ParseDebugId(src, tag, extended)
	case EnableDebugger2TagCode:
		content, err = ParseEnableDebugger2(src, tag, extended)
	case ScriptLimitsTagCode:
		content, err = ParseScriptLimits(src, tag, extended)
	case SetTabIndexTagCode:
		content, err = ParseSetTabIndex(src, tag, extended)
	case FileAttributesTagCode:
		content, err = ParseFileAttributes(src, tag)
	case PlaceObject3TagCode:
		content, err = ParsePlaceObject3(src, tag, extended)
	case ImportAssets2TagCode:
		content, err = ParseImportAssets2(src, tag, extended)
	case DefineFontAlignZonesTagCode:
		content, err = ParseDefineFontAlignZones(src, tag, extended)
	case CsmTextSettingsTagCode:
		content, err = ParseCsmTextSettings(src, tag, extended)
	case DefineFont3TagCode:
		content, err = ParseDefineFont3(src, tag, extended)
	case SymbolClassTagCode:
		content, err = ParseSymbolClass(src, tag, extended)
	case MetadataTagCode:
		content, err = ParseMetadata(src, tag, extended)
	case DefineScalingGridTagCode:
		content, err = ParseDefineScalingGrid(src, tag, extended)
	case DoAbcTagCode:
		content, err = ParseDoAbc(src, tag, extended)
	case DefineShape4TagCode:
		content, err = ParseDefineShape4(src, tag, extended)
	case DefineMorphShape2TagCode:
		content, err = ParseDefineMorphShape2(src, tag, extended)
	case DefineSceneAndFrameLabelDataTagCode:
		content, err = ParseDefineSceneAndFrameLabelData(src, tag, extended)
	case DefineBinaryDataTagCode:
		content, err = ParseDefineBinaryData(src, tag, extended)
	case DefineFontNameTagCode:
		content, err = ParseDefineFontName(src, tag, extended)
	case StartSound2TagCode:
		content, err = ParseStartSound2(src, tag, extended)
	case DefineBitsJpeg4TagCode:
		content, err = ParseDefineBitsJpeg4(src, tag, extended)
	case DefineFont4TagCode:
		content, err = ParseDefineFont4(src, tag, extended)
	case EnableTelemetryTagCode:
		content, err = ParseEnableTelemetry(src, tag, extended)
	case PlaceObject4TagCode:
		content, err = ParsePlaceObject4(src, tag, extended)
	default:
		content, err = ParseUnknown(src, tag, extended)
	}

	return content, err
}
