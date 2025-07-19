package constants

// QRProvider defines supported QR code providers
type QRProvider string

const (
	QRProviderVIETQR QRProvider = "VIETQR"
	QRProviderVNPAY  QRProvider = "VNPAY"
)

// QRProviderGUID maps each provider to its AID (Application Identifier)
type QRProviderGUID string

const (
	GUIDVNPAY  QRProviderGUID = "A000000775"
	GUIDVIETQR QRProviderGUID = "A000000727"
)

// FieldID defines standard EMVCo field IDs
type FieldID string

const (
	FieldVersion         FieldID = "00"
	FieldInitMethod      FieldID = "01"
	FieldVNPAYQR         FieldID = "26"
	FieldVIETQR          FieldID = "38"
	FieldCategory        FieldID = "52"
	FieldCurrency        FieldID = "53"
	FieldAmount          FieldID = "54"
	FieldTipAndFeeType   FieldID = "55"
	FieldTipAndFeeAmount FieldID = "56"
	FieldTipAndFeePct    FieldID = "57"
	FieldNation          FieldID = "58"
	FieldMerchantName    FieldID = "59"
	FieldCity            FieldID = "60"
	FieldZipCode         FieldID = "61"
	FieldAdditionalData  FieldID = "62"
	FieldCRC             FieldID = "63"
)

// EVMCoFieldIDs: 65-79
var EVMCoFieldIDs = []string{
	"65", "66", "67", "68", "69", "70", "71", "72", "73", "74", "75", "76", "77", "78", "79",
}

// UnreservedFieldIDs: 80-99
var UnreservedFieldIDs = []string{
	"80", "81", "82", "83", "84", "85", "86", "87", "88", "89",
	"90", "91", "92", "93", "94", "95", "96", "97", "98", "99",
}

// ProviderFieldID defines sub-fields under provider field (like 26 or 38)
type ProviderFieldID string

const (
	ProviderFieldGUID    ProviderFieldID = "00"
	ProviderFieldData    ProviderFieldID = "01"
	ProviderFieldService ProviderFieldID = "02"
)

// VietQRService defines types of VietQR services
type VietQRService string

const (
	VietQRByAccount VietQRService = "QRIBFTTA"
	VietQRByCard    VietQRService = "QRIBFTTC"
)

// VietQRConsumerFieldID defines the sub-fields in VietQR data
type VietQRConsumerFieldID string

const (
	VietQRBankBin    VietQRConsumerFieldID = "00"
	VietQRBankNumber VietQRConsumerFieldID = "01"
)

// AdditionalDataID maps to optional data fields (field 62)
type AdditionalDataID string

const (
	AdditionalBillNumber           AdditionalDataID = "01" // sổ hóa đơn
	AdditionalMobileNumber         AdditionalDataID = "02" // số điện thoại
	AdditionalStoreLabel           AdditionalDataID = "03" // mã cửa hàng
	AdditionalLoyaltyNumber        AdditionalDataID = "04" // mã khách hàng thân thiết
	AdditionalReferenceLabel       AdditionalDataID = "05" // mã tham chiếu
	AdditionalCustomerLabel        AdditionalDataID = "06" // mã khách hàng
	AdditionalTerminalLabel        AdditionalDataID = "07" // mã số điểm bán
	AdditionalPurposeOfTransaction AdditionalDataID = "08" // mục đích giao dịch
	AdditionalDataRequest          AdditionalDataID = "09" // yêu cầ uduwx liệu khách hàng bổ sung
)

type Provide struct {
	fieldId *string
	name    *QRProvider
	guid    *string
	service *string
	data    *string
}

type AdditionalData struct {
	billNumber    *string
	mobileNumber  *string
	store         *string
	loyaltyNumber *string
	reference     *string
	customerLabel *string
	terminal      *string
	purpose       *string
	dataRequest   *string
}

type Consumer struct {
	bankBin    *string
	bankNumber *string
}

type Merchant struct {
	id   *string
	name *string
}
