package constants

// QRProvider defines supported QR code providers
type QRProvider string

const (
	QRProviderVIETQR QRProvider = "VIETQR"
	QRProviderVNPAY  QRProvider = "VNPAY"
)

func (q QRProvider) String() string {
	return string(q)
}

type QRProviderGUID string

const (
	GUIDVNPAY  QRProviderGUID = "A000000775"
	GUIDVIETQR QRProviderGUID = "A000000727"
)

// FieldID defines standard EMVCo field IDs
type FieldID string

const (
	FieldIDVersion          FieldID = "00"
	FieldIDInitMethod       FieldID = "01"
	FieldIDVNPAYQR          FieldID = "26"
	FieldIDVIETQR           FieldID = "38"
	FieldIDCategory         FieldID = "52"
	FieldIDCurrency         FieldID = "53"
	FieldIDAmount           FieldID = "54"
	FieldIDTipAndFeeType    FieldID = "55"
	FieldIDTipAndFeeAmount  FieldID = "56"
	FieldIDTipAndFeePercent FieldID = "57"
	FieldIDNation           FieldID = "58"
	FieldIDMerchantName     FieldID = "59"
	FieldIDCity             FieldID = "60"
	FieldIDZipCode          FieldID = "61"
	FieldIDAdditionalData   FieldID = "62"
	FieldIDCRC              FieldID = "63"
)

func (f FieldID) String() string {
	return string(f)
}

var EVMCoFieldIDs = []string{
	"65", "66", "67", "68", "69", "70", "71", "72", "73", "74", "75", "76", "77", "78", "79",
}

func GetEVMCoFieldIDs() []string {
	return EVMCoFieldIDs
}

var UnreservedFieldIDs = []string{
	"80", "81", "82", "83", "84", "85", "86", "87", "88", "89",
	"90", "91", "92", "93", "94", "95", "96", "97", "98", "99",
}

func GetUnreservedFieldIDs() []string {
	return UnreservedFieldIDs
}

type ProviderFieldID string

const (
	ProviderFieldIDGUID    ProviderFieldID = "00"
	ProviderFieldIDData    ProviderFieldID = "01"
	ProviderFieldIDService ProviderFieldID = "02"
)

func (p ProviderFieldID) String() string {
	return string(p)
}

type VietQRService string

const (
	VietQRByAccount VietQRService = "QRIBFTTA"
	VietQRByCard    VietQRService = "QRIBFTTC"
)

// VietQRConsumerFieldID defines the sub-fields in VietQR data
type VietQRConsumerFieldID string

const (
	VietQRConsumerFieldIDBankBin    VietQRConsumerFieldID = "00"
	VietQRConsumerFieldIDBankNumber VietQRConsumerFieldID = "01"
)

func (v VietQRConsumerFieldID) String() string {
	return string(v)
}

// AdditionalDataID maps to optional data fields (field 62)
type AdditionalDataID string

const (
	AdditionalDataIDBillNumber                    AdditionalDataID = "01" // sổ hóa đơn
	AdditionalDataIDMobileNumber                  AdditionalDataID = "02" // số điện thoại
	AdditionalDataIDStoreLabel                    AdditionalDataID = "03" // mã cửa hàng
	AdditionalDataIDLoyaltyNumber                 AdditionalDataID = "04" // mã khách hàng thân thiết
	AdditionalDataIDReferenceLabel                AdditionalDataID = "05" // mã tham chiếu
	AdditionalDataIDCustomerLabel                 AdditionalDataID = "06" // mã khách hàng
	AdditionalDataIDTerminalLabel                 AdditionalDataID = "07" // mã số điểm bán
	AdditionalDataIDPurposeOfTransaction          AdditionalDataID = "08" // mục đích giao dịch
	AdditionalDataIDAdditionalConsumerDataRequest AdditionalDataID = "09" // yêu cầ uduwx liệu khách hàng bổ sung
)

func (a AdditionalDataID) String() string {
	return string(a)
}

type Provider struct {
	FieldId string
	Name    QRProvider
	GUID    string
	Service string
	Data    string
}

type AdditionalData struct {
	BillNumber    string
	MobileNumber  string
	Store         string
	LoyaltyNumber string
	Reference     string
	CustomerLabel string
	Terminal      string
	Purpose       string
	DataRequest   string
}

type Consumer struct {
	BankBin    string
	BankNumber string
}

type Merchant struct {
	Id   string
	Name string
}
