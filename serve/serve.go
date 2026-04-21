package serve

import "qrpay/constants"

type VNPaymentRequest struct {
	MerchantID   string `json:"merchantId"`
	MerchantName string `json:"merchantName"`
	Store        string `json:"store"`
	Terminal     string `json:"terminal"`

	// additional data
	Amount        *string `json:"amount,omitempty"`
	Purpose       *string `json:"purpose,omitempty"`
	BillNumber    *string `json:"billNumber,omitempty"`
	MobileNumber  *string `json:"mobileNumber,omitempty"`
	LoyaltyNumber *string `json:"loyaltyNumber,omitempty"`
	Reference     *string `json:"reference,omitempty"`
	CustomerLabel *string `json:"customerLabel,omitempty"`
}

//func InitVNPayQR(options VNPaymentRequest) QRPay {
//	var qr = NewQRPay("hello")
//	qr.Merchant.id = options.MerchantID
//	qr.merchant.name = options.MerchantName
//	qr.provider.fieldId = constants.FieldIDVNPAYQR
//	qr.provider.guid = QRProviderGUID.VNPAY
//	qr.provider.name = QRProvider.VNPAY
//	qr.amount = options.amount
//	qr.additionalData.purpose = options.purpose
//	qr.additionalData.billNumber = options.billNumber
//	qr.additionalData.mobileNumber = options.mobileNumber
//	qr.additionalData.store = options.store
//	qr.additionalData.terminal = options.terminal
//	qr.additionalData.loyaltyNumber = options.loyaltyNumber
//	qr.additionalData.reference = options.reference
//	qr.additionalData.customerLabel = options.customerLabel
//	return qr
//}

type InitVietQROptions struct {
	BankBin    string
	BankNumber string
	Amount     string
	Purpose    string
	Service    string
}

func InitVietQR(options InitVietQROptions) *qrPay {
	qr := NewQRPay("")

	if options.Amount != "" {
		qr.InitMethod = "12"
	} else {
		qr.InitMethod = "11"
	}

	qr.Provider.FieldId = constants.FieldIDVIETQR.String()
	qr.Provider.GUID = constants.QRProviderGUIDVIETQR.String()
	qr.Provider.Name = constants.QRProviderVIETQR

	if options.Service != "" {
		qr.Provider.Service = options.Service
	} else {
		qr.Provider.Service = string(constants.VietQRServiceByAccount)
	}

	qr.Consumer.BankBin = options.BankBin
	qr.Consumer.BankNumber = options.BankNumber
	qr.Amount = options.Amount
	qr.AdditionalData.Purpose = options.Purpose

	return qr
}
