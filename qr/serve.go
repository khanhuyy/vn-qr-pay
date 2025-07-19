package qr

import (
	"fmt"
	"github.com/howeyc/crc16"
	"sort"
	"strconv"
	"strings"
)

const (
	PF_GUID    = "00"
	PF_DATA    = "01"
	PF_SERVICE = "02"

	GUID_VNPAY  = "A000000775"
	GUID_VIETQR = "A000000727"

	QRProviderVNPAY  = "VNPAY"
	QRProviderVIETQR = "VIETQR"

	VQF_BANK_BIN    = "00"
	VQF_BANK_NUMBER = "01"
)

// Field IDs
const (
	FieldVersion       = "00"
	FieldInitMethod    = "01"
	FieldVNPAYQR       = "26"
	FieldVIETQR        = "38"
	FieldCategory      = "52"
	FieldCurrency      = "53"
	FieldAmount        = "54"
	FieldTipAndFeeType = "55"
	FieldTipAndFeeAmt  = "56"
	FieldTipAndFeePct  = "57"
	FieldNation        = "58"
	FieldMerchantName  = "59"
	FieldCity          = "60"
	FieldZipCode       = "61"
	FieldAdditional    = "62"
	FieldCRC           = "63"
)

// ProviderFieldID
const (
	ProvFieldGUID    = "00"
	ProvFieldData    = "01"
	ProvFieldService = "02"
)

// QRProvider GUIDs
const (
	GUIDVNPAY  = "A000000775"
	GUIDVIETQR = "A000000727"
)

// VietQR Consumer FieldID
const (
	ConsFieldBankBin    = "00"
	ConsFieldBankNumber = "01"
)

// AdditionalData IDs
const (
	ADBillNumber      = "01"
	ADMobileNumber    = "02"
	ADStoreLabel      = "03"
	ADLoyaltyNumber   = "04"
	ADReferenceLabel  = "05"
	ADCustomerLabel   = "06"
	ADTerminalLabel   = "07"
	ADPurpose         = "08"
	ADConsumerDataReq = "09"
)

// Structs
type Provider struct {
	FieldID string
	Name    string
	GUID    string
	Service string
	Data    string
}

type Merchant struct {
	ID   string
	Name string
}

type Consumer struct {
	BankBin    string
	BankNumber string
}

type AdditionalData struct {
	BillNumber   string
	MobileNumber string
	Store        string
	Loyalty      string
	Reference    string
	Customer     string
	Terminal     string
	Purpose      string
	DataRequest  string
}

type QRPay struct {
	IsValid       bool
	Version       string
	InitMethod    string
	Provider      Provider
	Merchant      Merchant
	Consumer      Consumer
	Category      string
	Currency      string
	Amount        string
	TipFeeType    string
	TipFeeAmount  string
	TipFeePercent string
	Nation        string
	City          string
	ZipCode       string
	Additional    AdditionalData
	EVMCo         map[string]string
	Unreserved    map[string]string
	CRC           string
}

func sliceContent(content string) (id, value, next string, err error) {
	if len(content) < 4 {
		return "", "", "", fmt.Errorf("invalid length")
	}
	id = content[:2]
	length, err := strconv.Atoi(content[2:4])
	if err != nil || len(content) < 4+length {
		return "", "", "", fmt.Errorf("invalid length")
	}
	value = content[4 : 4+length]
	next = content[4+length:]
	return id, value, next, nil
}

func genFieldData(id, value string) string {
	if len(id) != 2 || len(value) == 0 {
		return ""
	}
	return fmt.Sprintf("%s%02d%s", id, len(value), value)
}

func genCRC(content string) string {
	crc := crc16.Checksum([]byte(content), crc16.CCITTFalseTable)
	return fmt.Sprintf("%04X", crc)
}

func NewQRPay(payload string) *QRPay {
	q := &QRPay{
		IsValid:    true,
		EVMCo:      make(map[string]string),
		Unreserved: make(map[string]string),
	}
	if !q.verifyCRC(payload) {
		q.IsValid = false
		return q
	}
	q.parseRootContent(payload)
	return q
}

func (q *QRPay) verifyCRC(content string) bool {
	if len(content) < 4 {
		return false
	}
	data := content[:len(content)-4]
	crcStr := strings.ToUpper(content[len(content)-4:])
	return crcStr == genCRC(data)
}

func (q *QRPay) Build() string {
	var builder strings.Builder

	// 00: Version
	builder.WriteString(genFieldData(FieldVersion, defaultIfEmpty(q.Version, "01")))

	// 01: Init Method
	builder.WriteString(genFieldData(FieldInitMethod, defaultIfEmpty(q.InitMethod, "11")))

	// 26 / 38: Provider block (VNPay or VietQR)
	providerContent := ""
	providerContent += genFieldData(ProvFieldGUID, q.Provider.GUID)

	switch q.Provider.GUID {
	case GUIDVIETQR:
		providerContent += genFieldData(ProvFieldData,
			genFieldData(ConsFieldBankBin, q.Consumer.BankBin)+
				genFieldData(ConsFieldBankNumber, q.Consumer.BankNumber))
		q.Provider.Name = "VietQR"
	case GUIDVNPAY:
		providerContent += genFieldData(ProvFieldData, q.Merchant.ID)
		q.Provider.Name = "VNPay"
	default:
		providerContent += genFieldData(ProvFieldData, q.Provider.Data)
	}
	providerContent += genFieldData(ProvFieldService, q.Provider.Service)

	builder.WriteString(genFieldData(q.Provider.FieldID, providerContent))

	// 52: Category
	builder.WriteString(genFieldData(FieldCategory, q.Category))

	// 53: Currency (default: 704)
	builder.WriteString(genFieldData(FieldCurrency, defaultIfEmpty(q.Currency, "704")))

	// 54: Amount
	builder.WriteString(genFieldData(FieldAmount, q.Amount))

	// 55 - 57: Tip & fee
	builder.WriteString(genFieldData(FieldTipAndFeeType, q.TipFeeType))
	builder.WriteString(genFieldData(FieldTipAndFeeAmt, q.TipFeeAmount))
	builder.WriteString(genFieldData(FieldTipAndFeePct, q.TipFeePercent))

	// 58: Nation
	builder.WriteString(genFieldData(FieldNation, defaultIfEmpty(q.Nation, "VN")))

	// 59: Merchant Name
	builder.WriteString(genFieldData(FieldMerchantName, q.Merchant.Name))

	// 60 - 61: City, Zip
	builder.WriteString(genFieldData(FieldCity, q.City))
	builder.WriteString(genFieldData(FieldZipCode, q.ZipCode))

	// 62: Additional Data
	additional := ""
	additional += genFieldData(ADBillNumber, q.Additional.BillNumber)
	additional += genFieldData(ADMobileNumber, q.Additional.MobileNumber)
	additional += genFieldData(ADStoreLabel, q.Additional.Store)
	additional += genFieldData(ADLoyaltyNumber, q.Additional.Loyalty)
	additional += genFieldData(ADReferenceLabel, q.Additional.Reference)
	additional += genFieldData(ADCustomerLabel, q.Additional.Customer)
	additional += genFieldData(ADTerminalLabel, q.Additional.Terminal)
	additional += genFieldData(ADPurpose, q.Additional.Purpose)
	additional += genFieldData(ADConsumerDataReq, q.Additional.DataRequest)

	if len(additional) > 0 {
		builder.WriteString(genFieldData(FieldAdditional, additional))
	}

	// 65-79: EVMCo fields
	if len(q.EVMCo) > 0 {
		keys := sortedKeys(q.EVMCo)
		for _, k := range keys {
			builder.WriteString(genFieldData(k, q.EVMCo[k]))
		}
	}

	// 80-99: Unreserved fields
	if len(q.Unreserved) > 0 {
		keys := sortedKeys(q.Unreserved)
		for _, k := range keys {
			builder.WriteString(genFieldData(k, q.Unreserved[k]))
		}
	}

	// 63: CRC placeholder
	crcBase := builder.String() + FieldCRC + "04"
	crc := genCRC(crcBase)
	return crcBase + crc
}

func (q *QRPay) parseRootContent(content string) {
	for len(content) > 4 {
		id, val, next, err := sliceContent(content)
		if err != nil {
			q.IsValid = false
			return
		}
		switch id {
		case FieldVersion:
			q.Version = val
		case FieldInitMethod:
			q.InitMethod = val
		case FieldVNPAYQR, FieldVIETQR:
			q.Provider.FieldID = id
			q.parseProviderInfo(val)
		case FieldCategory:
			q.Category = val
		case FieldCurrency:
			q.Currency = val
		case FieldAmount:
			q.Amount = val
		case FieldTipAndFeeType:
			q.TipFeeType = val
		case FieldTipAndFeeAmt:
			q.TipFeeAmount = val
		case FieldTipAndFeePct:
			q.TipFeePercent = val
		case FieldNation:
			q.Nation = val
		case FieldMerchantName:
			q.Merchant.Name = val
		case FieldCity:
			q.City = val
		case FieldZipCode:
			q.ZipCode = val
		case FieldAdditional:
			q.parseAdditionalData(val)
		case FieldCRC:
			q.CRC = val
		default:

			n := parseInt(id)
			if 65 <= n && n <= 79 {
				q.EVMCo[id] = val
			} else if 80 <= n && n <= 99 {
				q.Unreserved[id] = val
			}
		}
		content = next
	}
}

func defaultIfEmpty(val, def string) string {
	if val == "" {
		return def
	}
	return val
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (q *QRPay) parseAdditionalData(content string) {
	for len(content) >= 4 {
		id := content[0:2]
		length, err := strconv.Atoi(content[2:4])
		if err != nil || len(content) < 4+length {
			q.IsValid = false
			return
		}
		value := content[4 : 4+length]
		content = content[4+length:]

		switch id {
		case ADBillNumber:
			q.Additional.BillNumber = value
		case ADMobileNumber:
			q.Additional.MobileNumber = value
		case ADStoreLabel:
			q.Additional.Store = value
		case ADLoyaltyNumber:
			q.Additional.Loyalty = value
		case ADReferenceLabel:
			q.Additional.Reference = value
		case ADCustomerLabel:
			q.Additional.Customer = value
		case ADTerminalLabel:
			q.Additional.Terminal = value
		case ADPurpose:
			q.Additional.Purpose = value
		case ADConsumerDataReq:
			q.Additional.DataRequest = value
		}
	}
}

func (q *QRPay) parseProviderInfo(content string) {
	for len(content) >= 4 {
		id := content[0:2]
		length, err := strconv.Atoi(content[2:4])
		if err != nil || len(content) < 4+length {
			q.IsValid = false
			return
		}
		value := content[4 : 4+length]
		content = content[4+length:]

		switch id {
		case PF_GUID:
			q.Provider.GUID = value

		case PF_DATA:
			q.Provider.Data = value
			switch q.Provider.GUID {
			case GUID_VNPAY:
				q.Provider.Name = QRProviderVNPAY
				q.Merchant.ID = value
			case GUID_VIETQR:
				q.Provider.Name = QRProviderVIETQR
				q.parseVietQRConsumer(value)
			}

		case PF_SERVICE:
			q.Provider.Service = value
		}
	}
}

func (q *QRPay) parseVietQRConsumer(content string) {
	for len(content) >= 4 {
		id := content[0:2]
		length, err := strconv.Atoi(content[2:4])
		if err != nil || len(content) < 4+length {
			q.IsValid = false
			return
		}
		value := content[4 : 4+length]
		content = content[4+length:]

		switch id {
		case VQF_BANK_BIN:
			q.Consumer.BankBin = value
		case VQF_BANK_NUMBER:
			q.Consumer.BankNumber = value
		}
	}
}

func parseInt(s string) int64 {
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return int64(value)
}
