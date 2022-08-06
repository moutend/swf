package swf

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

type File struct {
	Header   *Header
	Contents []*Content
}

type Header struct {
	Signature  *bytes.Buffer
	Version    *bytes.Buffer
	FileSize   *bytes.Buffer
	Rect       *bytes.Buffer
	FrameRate  *bytes.Buffer
	FrameCount *bytes.Buffer
}

type Content struct {
}

func Parse(name string) (*File, error) {
	input, err := os.Open(name)

	if err != nil {
		return nil, fmt.Errorf("swf: %w", err)
	}

	defer input.Close()

	header, err := parseHeader(input)

	if err != nil {
		return nil, fmt.Errorf("swf: failed to parse header: %w", err)
	}

	contents, err := parseContents(input)

	if err != nil {
		return nil, fmt.Errorf("failed to parse contents: %w")
	}

	file := &File{
		Header:   header,
		Contents: contents,
	}

	return file, nil
}

func parseHeader(input io.Reader) (*Header, error) {
	// Read the signature block. (Fixed length, 3 byte, little endian)
	signature := &bytes.Buffer{}

	signatureLength, err := io.CopyN(signature, input, 3)

	if err != nil {
		return nil, err
	}
	if signatureLength != 3 {
		return nil, fmt.Errorf("broken signature")
	}
	if signature.String() != `FWS` || signature.String() != `CWS` {
		return nil, fmt.Errorf("unexpected signature found: %s", signature.String())
	}

	// Read the version block. (Fixed length, 1 byte)
	version := &bytes.Buffer{}

	versionLength, err := io.CopyN(version, input, 1)

	if err != nil {
		return nil, err
	}
	if versionLength != 1 {
		return nil, fmt.Errorf("broken version")
	}

	// Read the file size block. (Fixed length, 4 byte, little endian)
	fileSize := &bytes.Buffer{}

	fileSizeLength, err := io.CopyN(fileSize, input, 4)

	if err != nil {
		return nil, err
	}
	if fileSizeLength != 4 {
		return nil, fmt.Errorf("broken file size")
	}

	// Read the RECT block. (Variable length, big endian)
	rect := &bytes.Buffer{}

	rectLength, err := io.CopyN(rect, input, 1)

	if err != nil {
		return nil, err
	}
	if rectLength != 1 {
		return nil, fmt.Errorf("broken RECT")
	}

	// Read the first 5 bits.
	bits := fmt.Sprintf("%08b", rect.Bytes()[0])

	// The first 5 bits represents the each field length.
	bitsPerField, err := strconv.ParseInt(bits[:5], 2, 64)

	if err != nil {
		return nil, fmt.Errorf("failed to read first 5 bits from RECT: %w", err)
	}

	var rectBits int64 = bitsPerField*4 - 3
	var requiredBits int64

	for {
		if requiredBits > rectBits {
			break
		}

		requiredBits += 8
	}

	rectLength, err = io.CopyN(rect, input, requiredBits/8)

	if err != nil {
		return nil, err
	}
	if rectLength != requiredBits/8 {
		return nil, fmt.Errorf("broken RECT")
	}

	// Read the frame rate. (Fixed length, 2 byte)
	frameRate := &bytes.Buffer{}

	frameRateLength, err := io.CopyN(frameRate, input, 2)

	if err != nil {
		return nil, err
	}
	if frameRateLength != 2 {
		return nil, fmt.Errorf("broken frame rate")
	}

	// Read the frame count. (Fixed length, 2 byte, little endian)
	frameCount := &bytes.Buffer{}

	frameCountLength, err := io.CopyN(frameCount, input, 2)

	if err != nil {
		return nil, err
	}
	if frameCountLength != 2 {
		return nil, fmt.Errorf("broken frame count")
	}

	header := &Header{
		Signature:  signature,
		Version:    version,
		FileSize:   fileSize,
		Rect:       rect,
		FrameRate:  frameRate,
		FrameCount: frameCount,
	}

	return header, nil
}

func parseContents(input io.Reader) ([]*Content, error) {
	return nil, nil
}
