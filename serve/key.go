package serve

import "strconv"

// EMVCoQRField Example
// 000201
// parse
//
// 000201
// 010211
// 38580010A000000727012800069704070114190338681560110208QRIBFTTA
// 5303704
// 5802VN
// 8300
// 8400
// 63043C31
//
// ----------------------
// 000201
// 010212
// 38260010A0000007270208QRIBFTTA
// 530370454061000005802VN62080804test6304B031
//
// 00020101021138570010A00000072701270006970425011354564645645670208QRIBFTTA53037045802VN6304CF4D
// 000201
// 010211
// 3857
//    0010A000000727
//    0127000697042501135456464564567
//    0208QRIBFTTA
// 5303704
// 5802VN
// 6304CF4D

type EMVCoQRField string

const (
	FieldPayloadFormat        = "00"
	FieldInitiationMethod     = "01"
	FieldMCC                  = "52"
	FieldCurrency             = "53"
	FieldAmount               = "54"
	FieldTipIndicator         = "55"
	FieldConvenienceFee       = "56"
	FieldCountry              = "58"
	FieldMerchantName         = "59"
	FieldMerchantCity         = "60"
	FieldPostalCode           = "61"
	FieldAdditionalData       = "62"
	FieldCRC                  = "63"
	FieldMerchantInfoLanguage = "64"
)

func (f EMVCoQRField) String() string {
	return string(f)
}

func IsMerchantAccount(tag string) bool {
	n, _ := strconv.Atoi(tag)
	return n >= 2 && n <= 51
}

func IsRFU(tag string) bool {
	n, _ := strconv.Atoi(tag)
	return n >= 65 && n <= 79
}

func IsUnreserved(tag string) bool {
	n, _ := strconv.Atoi(tag)
	return n >= 80 && n <= 99
}
