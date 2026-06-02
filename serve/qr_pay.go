package serve

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/khanhuyy/emvgo/constants"
)

type QRPay interface {
	Parse(content string)
	Build() string
	SetEVMCoField(id, value string)
	SetUnreservedField(id, value string)
	GenCRCCode(content string) string
	parseRootContent(content string)
	parseProviderInfo(content string)
	parseVietQRConsumer(content string)
	parseAdditionalData(content string)
	verifyCRC(content string) bool
	sliceContent(content string) (string, int, string, string)
	invalid()
	genFieldData(id string, value string) string
}

// Pay is the low-level Vietnam QR container. Most callers want bankqr instead.
type Pay struct {
	IsValid          bool
	Version          string
	InitMethod       string
	Provider         constants.Provider
	Merchant         constants.Merchant
	Consumer         constants.Consumer
	Category         string
	Currency         string
	Amount           string
	TipAndFeeType    string
	TipAndFeeAmount  string
	TipAndFeePercent string
	Nation           string
	City             string
	ZipCode          string
	AdditionalData   constants.AdditionalData
	CRC              string

	EVMCo      map[string]string
	Unreserved map[string]string
}

func (qr *Pay) InitVNPayQR() {
	//TODO implement me
	panic("implement me")
}

func NewQRPay(content string) *Pay {
	qr := &Pay{
		IsValid:        true,
		Provider:       constants.Provider{},
		Consumer:       constants.Consumer{},
		Merchant:       constants.Merchant{},
		AdditionalData: constants.AdditionalData{},
		EVMCo:          map[string]string{},
		Unreserved:     map[string]string{},
	}
	qr.Parse(content)
	return qr
}

func (qr *Pay) Parse(content string) {
	// Empty = fresh build mode.
	if content == "" {
		return
	}
	if len(content) < 4 {
		qr.invalid()
		return
	}
	if !qr.verifyCRC(content) {
		qr.invalid()
		return
	}
	qr.parseRootContent(content)
}

func (qr *Pay) Build() string {
	version := qr.genFieldData(constants.FieldIDVersion.String(), defaultStr(qr.Version, "01"))
	initMethod := qr.genFieldData(constants.FieldIDInitMethod.String(), defaultStr(qr.InitMethod, "11"))

	guid := qr.genFieldData(constants.ProviderFieldIDGUID.String(), qr.Provider.GUID)

	var providerDataContent string
	switch qr.Provider.GUID {
	case constants.QRProviderGUIDVIETQR.String():
		providerDataContent = qr.genFieldData(constants.VietQRConsumerFieldIDBankBin.String(), qr.Consumer.BankBin) +
			qr.genFieldData(constants.VietQRConsumerFieldIDBankNumber.String(), qr.Consumer.BankNumber)
	case constants.QRProviderGUIDVNPAY.String():
		providerDataContent = qr.Merchant.Id
	default:
		providerDataContent = qr.Provider.Data
	}

	provider := qr.genFieldData(constants.ProviderFieldIDData.String(), providerDataContent)
	service := qr.genFieldData(constants.ProviderFieldIDService.String(), qr.Provider.Service)

	// Default field id from GUID.
	fieldId := qr.Provider.FieldId
	if fieldId == "" {
		switch qr.Provider.GUID {
		case constants.QRProviderGUIDVIETQR.String():
			fieldId = constants.FieldIDVIETQR.String()
		case constants.QRProviderGUIDVNPAY.String():
			fieldId = constants.FieldIDVNPAYQR.String()
		}
	}
	providerData := qr.genFieldData(fieldId, guid+provider+service)

	category := qr.genFieldData(constants.FieldIDCategory.String(), qr.Category)
	currency := qr.genFieldData(constants.FieldIDCurrency.String(), defaultStr(qr.Currency, "704"))
	amountStr := qr.genFieldData(constants.FieldIDAmount.String(), qr.Amount)
	tipAndFeeType := qr.genFieldData(constants.FieldIDTipAndFeeType.String(), qr.TipAndFeeType)
	tipAndFeeAmount := qr.genFieldData(constants.FieldIDTipAndFeeAmount.String(), qr.TipAndFeeAmount)
	tipAndFeePercent := qr.genFieldData(constants.FieldIDTipAndFeePercent.String(), qr.TipAndFeePercent)
	nation := qr.genFieldData(constants.FieldIDNation.String(), defaultStr(qr.Nation, "VN"))
	merchantName := qr.genFieldData(constants.FieldIDMerchantName.String(), qr.Merchant.Name)
	city := qr.genFieldData(constants.FieldIDCity.String(), qr.City)
	zipCode := qr.genFieldData(constants.FieldIDZipCode.String(), qr.ZipCode)

	buildNumber := qr.genFieldData(constants.AdditionalDataIDBillNumber.String(), qr.AdditionalData.BillNumber)
	mobileNumber := qr.genFieldData(constants.AdditionalDataIDMobileNumber.String(), qr.AdditionalData.MobileNumber)
	storeLabel := qr.genFieldData(constants.AdditionalDataIDStoreLabel.String(), qr.AdditionalData.Store)
	loyaltyNumber := qr.genFieldData(constants.AdditionalDataIDLoyaltyNumber.String(), qr.AdditionalData.LoyaltyNumber)
	reference := qr.genFieldData(constants.AdditionalDataIDReferenceLabel.String(), qr.AdditionalData.Reference)
	customerLabel := qr.genFieldData(constants.AdditionalDataIDCustomerLabel.String(), qr.AdditionalData.CustomerLabel)
	terminal := qr.genFieldData(constants.AdditionalDataIDTerminalLabel.String(), qr.AdditionalData.Terminal)
	purpose := qr.genFieldData(constants.AdditionalDataIDPurposeOfTransaction.String(), qr.AdditionalData.Purpose)
	dataRequest := qr.genFieldData(constants.AdditionalDataIDAdditionalConsumerDataRequest.String(), qr.AdditionalData.DataRequest)

	additionalDataContent := buildNumber + mobileNumber + storeLabel + loyaltyNumber + reference + customerLabel + terminal + purpose + dataRequest
	additionalData := qr.genFieldData(constants.FieldIDAdditionalData.String(), additionalDataContent)

	var evmCoContent, unreservedContent string
	for _, k := range sortedKeys(qr.EVMCo) {
		evmCoContent += qr.genFieldData(k, qr.EVMCo[k])
	}
	for _, k := range sortedKeys(qr.Unreserved) {
		unreservedContent += qr.genFieldData(k, qr.Unreserved[k])
	}

	content := version + initMethod + providerData + category + currency + amountStr + tipAndFeeType + tipAndFeeAmount + tipAndFeePercent + nation + merchantName + city + zipCode + additionalData + evmCoContent + unreservedContent + constants.FieldIDCRC.String() + "04"
	crc := qr.GenCRCCode(content)
	return content + crc
}

func defaultStr(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

func (qr *Pay) InitVietQR() {

}

//func (qr *Pay) InitVNPayQR() QRPay {
//	return qr
//}

func (qr *Pay) SetEVMCoField(id, value string) {
	if qr.Unreserved == nil {
		qr.Unreserved = make(map[string]string)
	}
	qr.Unreserved[id] = value
}

func (qr *Pay) SetUnreservedField(id, value string) {
	if qr.Unreserved == nil {
		qr.Unreserved = make(map[string]string)
	}
	qr.Unreserved[id] = value
}

func (qr *Pay) parseRootContent(content string) {
	id, length, value, nextValue := qr.sliceContent(content)
	if len(value) != length {
		return
	}
	switch id {
	case constants.FieldIDVersion.String():
		qr.Version = value
	case constants.FieldIDInitMethod.String():
		qr.InitMethod = value
	case constants.FieldIDVIETQR.String(), constants.FieldIDVNPAYQR.String():
		qr.Provider.FieldId = id
		qr.parseProviderInfo(value)
	case constants.FieldIDCategory.String():
		qr.Category = value
	case constants.FieldIDCurrency.String():
		qr.Currency = value
	case constants.FieldIDAmount.String():
		qr.Amount = value
	case constants.FieldIDTipAndFeeType.String():
		qr.TipAndFeeType = value
	case constants.FieldIDTipAndFeeAmount.String():
		qr.TipAndFeeAmount = value
	case constants.FieldIDTipAndFeePercent.String():
		qr.TipAndFeePercent = value
	case constants.FieldIDNation.String():
		qr.Nation = value
	case constants.FieldIDMerchantName.String():
		qr.Merchant.Name = value
	case constants.FieldIDCity.String():
		qr.City = value
	case constants.FieldIDZipCode.String():
		qr.ZipCode = value
	case constants.FieldIDAdditionalData.String():
		qr.parseAdditionalData(value)
	case constants.FieldIDCRC.String():
		qr.CRC = value
	default:
		idNum, _ := strconv.Atoi(id)
		if idNum >= 65 && idNum <= 79 {
			if qr.EVMCo == nil {
				qr.EVMCo = make(map[string]string)
			}
			qr.EVMCo[id] = value
		} else if idNum >= 80 && idNum <= 99 {
			if qr.Unreserved == nil {
				qr.Unreserved = make(map[string]string)
			}
			qr.Unreserved[id] = value
		}
	}
	// cause of tlv always need 2 chars for tag and 2 chars for length
	if len(nextValue) >= 4 {
		qr.parseRootContent(nextValue)
	}
}

func (qr *Pay) parseProviderInfo(content string) {
	id, _, value, nextValue := qr.sliceContent(content)
	switch id {
	case constants.ProviderFieldIDGUID.String():
		qr.Provider.GUID = value
		// Resolve Name from AID.
		switch value {
		case constants.QRProviderGUIDVIETQR.String():
			qr.Provider.Name = constants.QRProviderVIETQR
		case constants.QRProviderGUIDVNPAY.String():
			qr.Provider.Name = constants.QRProviderVNPAY
		}
	case constants.ProviderFieldIDData.String():
		switch qr.Provider.Name {
		case constants.QRProviderVNPAY:
			qr.Merchant.Id = value
		case constants.QRProviderVIETQR:
			qr.parseVietQRConsumer(value)
		default:
			qr.Provider.Data = value
		}
	case constants.ProviderFieldIDService.String():
		qr.Provider.Service = value
	}
	if len(nextValue) >= 4 {
		qr.parseProviderInfo(nextValue)
	}
}

func (qr *Pay) parseVietQRConsumer(content string) {
	id, _, value, nextValue := qr.sliceContent(content)
	switch id {
	case constants.VietQRConsumerFieldIDBankBin.String():
		qr.Consumer.BankBin = value
	case constants.VietQRConsumerFieldIDBankNumber.String():
		qr.Consumer.BankNumber = value
	}
	if len(nextValue) > 4 {
		qr.parseVietQRConsumer(nextValue)
	}
}

func (qr *Pay) parseAdditionalData(content string) {
	id, _, value, nextValue := qr.sliceContent(content)
	switch id {
	case constants.AdditionalDataIDBillNumber.String():
		qr.AdditionalData.BillNumber = value
	case constants.AdditionalDataIDMobileNumber.String():
		qr.AdditionalData.MobileNumber = value
	case constants.AdditionalDataIDStoreLabel.String():
		qr.AdditionalData.Store = value
	case constants.AdditionalDataIDLoyaltyNumber.String():
		qr.AdditionalData.LoyaltyNumber = value
	case constants.AdditionalDataIDReferenceLabel.String():
		qr.AdditionalData.Reference = value
	case constants.AdditionalDataIDCustomerLabel.String():
		qr.AdditionalData.CustomerLabel = value
	case constants.AdditionalDataIDTerminalLabel.String():
		qr.AdditionalData.Terminal = value
	case constants.AdditionalDataIDPurposeOfTransaction.String():
		qr.AdditionalData.Purpose = value
	case constants.AdditionalDataIDAdditionalConsumerDataRequest.String():
		qr.AdditionalData.DataRequest = value
	}
	if len(nextValue) > 4 {
		qr.parseAdditionalData(nextValue)
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (qr *Pay) verifyCRC(content string) bool {
	if len(content) < 4 {
		return false
	}
	checkContent := content[:len(content)-4]
	crcCode := strings.ToUpper(content[len(content)-4:])
	genCRC := qr.GenCRCCode(checkContent)
	return crcCode == genCRC
}

func (qr *Pay) GenCRCCode(content string) string {
	crc := CRC16CCITT(content)
	return fmt.Sprintf("%04X", crc)
}

// sliceContent splits a TLV: tag, length, value, rest.
func (qr *Pay) sliceContent(content string) (id string, length int, value, nextValue string) {
	id = content[0:2]
	lengthString := content[2:4]
	length, _ = strconv.Atoi(lengthString)
	value = content[4 : 4+length]
	nextValue = content[4+length:]
	return id, length, value, nextValue
}

func (qr *Pay) invalid() {
	qr.IsValid = false
}

func (qr *Pay) genFieldData(id, value string) string {
	if len(id) != 2 || len(value) == 0 {
		return ""
	}
	length := fmt.Sprintf("%02d", len(value))
	return id + length + value
}
