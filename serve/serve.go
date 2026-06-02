package serve

import "github.com/khanhuyy/emvgo/constants"

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

// InitVNPayQR builds a VNPAY merchant QR (field 26).
func InitVNPayQR(options VNPaymentRequest) *Pay {
	qr := NewQRPay("")

	qr.Provider.FieldId = constants.FieldIDVNPAYQR.String()
	qr.Provider.GUID = constants.QRProviderGUIDVNPAY.String()
	qr.Provider.Name = constants.QRProviderVNPAY

	qr.Merchant.Id = options.MerchantID
	qr.Merchant.Name = options.MerchantName

	if options.Amount != nil {
		qr.Amount = *options.Amount
		qr.InitMethod = "12"
	} else {
		qr.InitMethod = "11"
	}
	if options.Purpose != nil {
		qr.AdditionalData.Purpose = *options.Purpose
	}
	if options.BillNumber != nil {
		qr.AdditionalData.BillNumber = *options.BillNumber
	}
	if options.MobileNumber != nil {
		qr.AdditionalData.MobileNumber = *options.MobileNumber
	}
	if options.Store != "" {
		qr.AdditionalData.Store = options.Store
	}
	if options.Terminal != "" {
		qr.AdditionalData.Terminal = options.Terminal
	}
	if options.LoyaltyNumber != nil {
		qr.AdditionalData.LoyaltyNumber = *options.LoyaltyNumber
	}
	if options.Reference != nil {
		qr.AdditionalData.Reference = *options.Reference
	}
	if options.CustomerLabel != nil {
		qr.AdditionalData.CustomerLabel = *options.CustomerLabel
	}
	return qr
}

type InitVietQROptions struct {
	BankBin    string
	BankNumber string
	Amount     string
	Purpose    string
	Service    string
}

func InitVietQR(options InitVietQROptions) *Pay {
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
