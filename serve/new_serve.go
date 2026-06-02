package serve

import (
	"fmt"
	"strconv"

	"github.com/khanhuyy/emvgo/constants"
)

type TLV struct {
	Tag    string
	Length int
	Value  string
}

func (t *TLV) Concat() string {
	return fmt.Sprintf("%s%02d%s", t.Tag, t.Length, t.Value)
}

type instance struct {
	PayloadFormat        string
	InitiationMethod     string
	MCC                  string
	Currency             string
	Amount               string
	TipIndicator         string
	ConvenienceFee       string
	Country              string
	MerchantName         string
	MerchantCity         string
	PostalCode           string
	AdditionalData       string
	CRC                  string
	MerchantInfoLanguage string
	MerchantFields       map[string]string // Fields 02-51 (Merchant Account)
	Consumer             constants.Consumer
	tlvs                 []*TLV
}

func (qr *instance) ParseTLVs(content string) error {
	if len(content) < 4 {
		if len(content) == 0 {
			return nil
		}
		return fmt.Errorf("content too short")
	}
	tag, length, value, nextContent := qr.sliceCurrentTLV(content)
	if len(value) != length {
		return fmt.Errorf("value length error")
	}
	switch tag {
	case FieldPayloadFormat:
		qr.PayloadFormat = value
	case FieldInitiationMethod:
		qr.InitiationMethod = value
	case FieldMCC:
		qr.MCC = value
	case FieldCurrency:
		qr.Currency = value
	case FieldAmount:
		qr.Amount = value
	case FieldTipIndicator:
		qr.TipIndicator = value
	case FieldConvenienceFee:
		qr.ConvenienceFee = value
	case FieldCountry:
		qr.Country = value
	case FieldMerchantName:
		qr.MerchantName = value
	case FieldMerchantCity:
		qr.MerchantCity = value
	case FieldPostalCode:
		qr.PostalCode = value
	case FieldAdditionalData:
		qr.AdditionalData = value
	case FieldCRC:
		qr.CRC = value
	case FieldMerchantInfoLanguage:
		qr.MerchantInfoLanguage = value
	default:
		if IsMerchantAccount(tag) {
			if qr.MerchantFields == nil {
				qr.MerchantFields = make(map[string]string)
			}
			qr.MerchantFields[tag] = value
		} else if IsRFU(tag) {
			fmt.Println("implement RFU")
		} else if IsUnreserved(tag) {
			fmt.Println("implement Unreserved")
		} else {
			return fmt.Errorf("unknown tag: %s", tag)
		}
	}
	return qr.ParseTLVs(nextContent)
}

func (qr *instance) sliceCurrentTLV(content string) (tag string, length int, value, nextValue string) {
	tag = content[0:2]
	lengthString := content[2:4]
	length, _ = strconv.Atoi(lengthString)
	value = content[4 : 4+length]
	nextValue = content[4+length:]
	return tag, length, value, nextValue
}

type BuildQROptions struct {
	BankBin    string
	BankNumber string
	Amount     string
	Remark     string
}

func NewMyQRPay() *instance {
	return &instance{}
}

// addField appends a TLV. Empty values are skipped.
func (qr *instance) addField(tag, value string) {
	if value == "" {
		return
	}
	qr.tlvs = append(qr.tlvs, &TLV{
		Tag:    tag,
		Length: len(value),
		Value:  value,
	})
}

// Build produces a VietQR (field 38) from opts.
func (qr *instance) Build(options BuildQROptions) string {
	qr.tlvs = []*TLV{}

	initMethod := defaultStr(qr.InitiationMethod, "11")
	if options.Amount != "" {
		initMethod = "12"
	}

	qr.addField(FieldPayloadFormat, defaultStr(qr.PayloadFormat, "01"))
	qr.addField(FieldInitiationMethod, initMethod)

	// Field 38: 00=GUID, 01=bank info, 02=service code.
	bankInfo := buildTLV("00", options.BankBin) + buildTLV("01", options.BankNumber)
	providerContent := buildTLV("00", "A000000727") +
		buildTLV("01", bankInfo) +
		buildTLV("02", string(constants.VietQRServiceByAccount))
	qr.tlvs = append(qr.tlvs, &TLV{Tag: "38", Length: len(providerContent), Value: providerContent})

	qr.addField(FieldMCC, qr.MCC)
	qr.addField(FieldCurrency, defaultStr(qr.Currency, "704"))
	qr.addField(FieldAmount, options.Amount)
	qr.addField(FieldTipIndicator, qr.TipIndicator)
	qr.addField(FieldConvenienceFee, qr.ConvenienceFee)
	qr.addField(FieldCountry, defaultStr(qr.Country, "VN"))
	qr.addField(FieldMerchantName, qr.MerchantName)
	qr.addField(FieldMerchantCity, qr.MerchantCity)
	qr.addField(FieldPostalCode, qr.PostalCode)

	// Field 62.08 = purpose / remark.
	if options.Remark != "" {
		additional := buildTLV("08", options.Remark)
		qr.tlvs = append(qr.tlvs, &TLV{Tag: FieldAdditionalData, Length: len(additional), Value: additional})
	}

	qr.addField(FieldMerchantInfoLanguage, qr.MerchantInfoLanguage)

	var res string
	for _, tlv := range qr.tlvs {
		res += tlv.Concat()
	}

	// CRC input includes the "6304" prefix.
	withCRCPrefix := res + FieldCRC + "04"
	crc := qr.buildCRC(withCRCPrefix)
	return withCRCPrefix + crc
}

func (qr *instance) buildCRC(content string) string {
	crc := CRC16CCITT(content)
	return fmt.Sprintf("%04X", crc)
}

func buildTLV(tag, value string) string {
	return fmt.Sprintf("%s%02d%s", tag, len(value), value)
}
