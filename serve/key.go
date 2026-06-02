package serve

import "strconv"

// EMVCoQRField is a 2-digit EMVCo tag.
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
