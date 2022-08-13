// Code generated by "stringer -type=TagCode"; DO NOT EDIT.

package swf

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[EndTagCode-0]
	_ = x[ShowFrameTagCode-1]
	_ = x[DefineShapeTagCode-2]
	_ = x[PlaceObjectTagCode-4]
	_ = x[RemoveObjectTagCode-5]
	_ = x[DefineBitsTagCode-6]
	_ = x[DefineButtonTagCode-7]
	_ = x[JpegTablesTagCode-8]
	_ = x[SetBackgroundColorTagCode-9]
	_ = x[DefineFontTagCode-10]
	_ = x[DefineTextTagCode-11]
	_ = x[DoActionTagCode-12]
	_ = x[DefineFontInfoTagCode-13]
	_ = x[DefineSoundTagCode-14]
	_ = x[StartSoundTagCode-15]
	_ = x[DefineButtonSoundTagCode-17]
	_ = x[SoundStreamHeadTagCode-18]
	_ = x[SoundStreamBlockTagCode-19]
	_ = x[DefineBitsLosslessTagCode-20]
	_ = x[DefineBitsJpeg2TagCode-21]
	_ = x[DefineShape2TagCode-22]
	_ = x[DefineButtonCxformTagCode-23]
	_ = x[ProtectTagCode-24]
	_ = x[PlaceObject2TagCode-26]
	_ = x[RemoveObject2TagCode-28]
	_ = x[DefineShape3TagCode-32]
	_ = x[DefineText2TagCode-33]
	_ = x[DefineButton2TagCode-34]
	_ = x[DefineBitsJpeg3TagCode-35]
	_ = x[DefineBitsLossless2TagCode-36]
	_ = x[DefineEditTextTagCode-37]
	_ = x[DefineSpriteTagCode-39]
	_ = x[NameCharacterTagCode-40]
	_ = x[ProductInfoTagCode-41]
	_ = x[FrameLabelTagCode-43]
	_ = x[SoundStreamHead2TagCode-45]
	_ = x[DefineMorphShapeTagCode-46]
	_ = x[DefineFont2TagCode-48]
	_ = x[ExportAssetsTagCode-56]
	_ = x[ImportAssetsTagCode-57]
	_ = x[EnableDebuggerTagCode-58]
	_ = x[DoInitActionTagCode-59]
	_ = x[DefineVideoStreamTagCode-60]
	_ = x[VideoFrameTagCode-61]
	_ = x[DefineFontInfo2TagCode-62]
	_ = x[DebugIdTagCode-63]
	_ = x[EnableDebugger2TagCode-64]
	_ = x[ScriptLimitsTagCode-65]
	_ = x[SetTabIndexTagCode-66]
	_ = x[FileAttributesTagCode-69]
	_ = x[PlaceObject3TagCode-70]
	_ = x[ImportAssets2TagCode-71]
	_ = x[DefineFontAlignZonesTagCode-73]
	_ = x[CsmTextSettingsTagCode-74]
	_ = x[DefineFont3TagCode-75]
	_ = x[SymbolClassTagCode-76]
	_ = x[MetadataTagCode-77]
	_ = x[DefineScalingGridTagCode-78]
	_ = x[DoAbcTagCode-82]
	_ = x[DefineShape4TagCode-83]
	_ = x[DefineMorphShape2TagCode-84]
	_ = x[DefineSceneAndFrameLabelDataTagCode-86]
	_ = x[DefineBinaryDataTagCode-87]
	_ = x[DefineFontNameTagCode-88]
	_ = x[StartSound2TagCode-89]
	_ = x[DefineBitsJpeg4TagCode-90]
	_ = x[DefineFont4TagCode-91]
	_ = x[EnableTelemetryTagCode-93]
	_ = x[PlaceObject4TagCode-94]
	_ = x[UnknownTagCode-65535]
}

const _TagCode_name = "EndTagCodeShowFrameTagCodeDefineShapeTagCodePlaceObjectTagCodeRemoveObjectTagCodeDefineBitsTagCodeDefineButtonTagCodeJpegTablesTagCodeSetBackgroundColorTagCodeDefineFontTagCodeDefineTextTagCodeDoActionTagCodeDefineFontInfoTagCodeDefineSoundTagCodeStartSoundTagCodeDefineButtonSoundTagCodeSoundStreamHeadTagCodeSoundStreamBlockTagCodeDefineBitsLosslessTagCodeDefineBitsJpeg2TagCodeDefineShape2TagCodeDefineButtonCxformTagCodeProtectTagCodePlaceObject2TagCodeRemoveObject2TagCodeDefineShape3TagCodeDefineText2TagCodeDefineButton2TagCodeDefineBitsJpeg3TagCodeDefineBitsLossless2TagCodeDefineEditTextTagCodeDefineSpriteTagCodeNameCharacterTagCodeProductInfoTagCodeFrameLabelTagCodeSoundStreamHead2TagCodeDefineMorphShapeTagCodeDefineFont2TagCodeExportAssetsTagCodeImportAssetsTagCodeEnableDebuggerTagCodeDoInitActionTagCodeDefineVideoStreamTagCodeVideoFrameTagCodeDefineFontInfo2TagCodeDebugIdTagCodeEnableDebugger2TagCodeScriptLimitsTagCodeSetTabIndexTagCodeFileAttributesTagCodePlaceObject3TagCodeImportAssets2TagCodeDefineFontAlignZonesTagCodeCsmTextSettingsTagCodeDefineFont3TagCodeSymbolClassTagCodeMetadataTagCodeDefineScalingGridTagCodeDoAbcTagCodeDefineShape4TagCodeDefineMorphShape2TagCodeDefineSceneAndFrameLabelDataTagCodeDefineBinaryDataTagCodeDefineFontNameTagCodeStartSound2TagCodeDefineBitsJpeg4TagCodeDefineFont4TagCodeEnableTelemetryTagCodePlaceObject4TagCodeUnknownTagCode"

var _TagCode_map = map[TagCode]string{
	0:     _TagCode_name[0:10],
	1:     _TagCode_name[10:26],
	2:     _TagCode_name[26:44],
	4:     _TagCode_name[44:62],
	5:     _TagCode_name[62:81],
	6:     _TagCode_name[81:98],
	7:     _TagCode_name[98:117],
	8:     _TagCode_name[117:134],
	9:     _TagCode_name[134:159],
	10:    _TagCode_name[159:176],
	11:    _TagCode_name[176:193],
	12:    _TagCode_name[193:208],
	13:    _TagCode_name[208:229],
	14:    _TagCode_name[229:247],
	15:    _TagCode_name[247:264],
	17:    _TagCode_name[264:288],
	18:    _TagCode_name[288:310],
	19:    _TagCode_name[310:333],
	20:    _TagCode_name[333:358],
	21:    _TagCode_name[358:380],
	22:    _TagCode_name[380:399],
	23:    _TagCode_name[399:424],
	24:    _TagCode_name[424:438],
	26:    _TagCode_name[438:457],
	28:    _TagCode_name[457:477],
	32:    _TagCode_name[477:496],
	33:    _TagCode_name[496:514],
	34:    _TagCode_name[514:534],
	35:    _TagCode_name[534:556],
	36:    _TagCode_name[556:582],
	37:    _TagCode_name[582:603],
	39:    _TagCode_name[603:622],
	40:    _TagCode_name[622:642],
	41:    _TagCode_name[642:660],
	43:    _TagCode_name[660:677],
	45:    _TagCode_name[677:700],
	46:    _TagCode_name[700:723],
	48:    _TagCode_name[723:741],
	56:    _TagCode_name[741:760],
	57:    _TagCode_name[760:779],
	58:    _TagCode_name[779:800],
	59:    _TagCode_name[800:819],
	60:    _TagCode_name[819:843],
	61:    _TagCode_name[843:860],
	62:    _TagCode_name[860:882],
	63:    _TagCode_name[882:896],
	64:    _TagCode_name[896:918],
	65:    _TagCode_name[918:937],
	66:    _TagCode_name[937:955],
	69:    _TagCode_name[955:976],
	70:    _TagCode_name[976:995],
	71:    _TagCode_name[995:1015],
	73:    _TagCode_name[1015:1042],
	74:    _TagCode_name[1042:1064],
	75:    _TagCode_name[1064:1082],
	76:    _TagCode_name[1082:1100],
	77:    _TagCode_name[1100:1115],
	78:    _TagCode_name[1115:1139],
	82:    _TagCode_name[1139:1151],
	83:    _TagCode_name[1151:1170],
	84:    _TagCode_name[1170:1194],
	86:    _TagCode_name[1194:1229],
	87:    _TagCode_name[1229:1252],
	88:    _TagCode_name[1252:1273],
	89:    _TagCode_name[1273:1291],
	90:    _TagCode_name[1291:1313],
	91:    _TagCode_name[1313:1331],
	93:    _TagCode_name[1331:1353],
	94:    _TagCode_name[1353:1372],
	65535: _TagCode_name[1372:1386],
}

func (i TagCode) String() string {
	if str, ok := _TagCode_map[i]; ok {
		return str
	}
	return "TagCode(" + strconv.FormatInt(int64(i), 10) + ")"
}
