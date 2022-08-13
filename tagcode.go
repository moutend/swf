//go:generate stringer -type=TagCode
package swf

type TagCode uint16

const (
	EndTagCode                          TagCode = 0
	ShowFrameTagCode                    TagCode = 1
	DefineShapeTagCode                  TagCode = 2
	PlaceObjectTagCode                  TagCode = 4
	RemoveObjectTagCode                 TagCode = 5
	DefineBitsTagCode                   TagCode = 6
	DefineButtonTagCode                 TagCode = 7
	JpegTablesTagCode                   TagCode = 8
	SetBackgroundColorTagCode           TagCode = 9
	DefineFontTagCode                   TagCode = 10
	DefineTextTagCode                   TagCode = 11
	DoActionTagCode                     TagCode = 12
	DefineFontInfoTagCode               TagCode = 13
	DefineSoundTagCode                  TagCode = 14
	StartSoundTagCode                   TagCode = 15
	DefineButtonSoundTagCode            TagCode = 17
	SoundStreamHeadTagCode              TagCode = 18
	SoundStreamBlockTagCode             TagCode = 19
	DefineBitsLosslessTagCode           TagCode = 20
	DefineBitsJpeg2TagCode              TagCode = 21
	DefineShape2TagCode                 TagCode = 22
	DefineButtonCxformTagCode           TagCode = 23
	ProtectTagCode                      TagCode = 24
	PlaceObject2TagCode                 TagCode = 26
	RemoveObject2TagCode                TagCode = 28
	DefineShape3TagCode                 TagCode = 32
	DefineText2TagCode                  TagCode = 33
	DefineButton2TagCode                TagCode = 34
	DefineBitsJpeg3TagCode              TagCode = 35
	DefineBitsLossless2TagCode          TagCode = 36
	DefineEditTextTagCode               TagCode = 37
	DefineSpriteTagCode                 TagCode = 39
	NameCharacterTagCode                TagCode = 40
	ProductInfoTagCode                  TagCode = 41
	FrameLabelTagCode                   TagCode = 43
	SoundStreamHead2TagCode             TagCode = 45
	DefineMorphShapeTagCode             TagCode = 46
	DefineFont2TagCode                  TagCode = 48
	ExportAssetsTagCode                 TagCode = 56
	ImportAssetsTagCode                 TagCode = 57
	EnableDebuggerTagCode               TagCode = 58
	DoInitActionTagCode                 TagCode = 59
	DefineVideoStreamTagCode            TagCode = 60
	VideoFrameTagCode                   TagCode = 61
	DefineFontInfo2TagCode              TagCode = 62
	DebugIdTagCode                      TagCode = 63
	EnableDebugger2TagCode              TagCode = 64
	ScriptLimitsTagCode                 TagCode = 65
	SetTabIndexTagCode                  TagCode = 66
	FileAttributesTagCode               TagCode = 69
	PlaceObject3TagCode                 TagCode = 70
	ImportAssets2TagCode                TagCode = 71
	DefineFontAlignZonesTagCode         TagCode = 73
	CsmTextSettingsTagCode              TagCode = 74
	DefineFont3TagCode                  TagCode = 75
	SymbolClassTagCode                  TagCode = 76
	MetadataTagCode                     TagCode = 77
	DefineScalingGridTagCode            TagCode = 78
	DoAbcTagCode                        TagCode = 82
	DefineShape4TagCode                 TagCode = 83
	DefineMorphShape2TagCode            TagCode = 84
	DefineSceneAndFrameLabelDataTagCode TagCode = 86
	DefineBinaryDataTagCode             TagCode = 87
	DefineFontNameTagCode               TagCode = 88
	StartSound2TagCode                  TagCode = 89
	DefineBitsJpeg4TagCode              TagCode = 90
	DefineFont4TagCode                  TagCode = 91
	EnableTelemetryTagCode              TagCode = 93
	PlaceObject4TagCode                 TagCode = 94
	UnknownTagCode                      TagCode = 0xffff
)
