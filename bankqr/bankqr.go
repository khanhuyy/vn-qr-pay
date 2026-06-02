// Package bankqr builds and parses Vietnam bank QR codes (VietQR + VNPAYQR).
package bankqr

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/khanhuyy/emvgo/constants"
	"github.com/khanhuyy/emvgo/serve"
)

// Service is the VietQR transfer service code (field 38.02).
type Service string

const (
	ServiceTransferToAccount Service = "QRIBFTTA" // transfer to account
	ServiceTransferToCard    Service = "QRIBFTTC" // transfer to card
)

// Options for building a VietQR.
//
// Required fields are plain values; optional fields are pointers so that
// "unset" and "zero value" are distinguishable. Use utils.StrPtr / Int64Ptr /
// the package's WithXxx helpers to construct pointer values inline.
type Options struct {
	// Required.
	BIN           constants.BankBIN // e.g. constants.BIN_VIETCOMBANK
	AccountNumber string

	// Optional.
	Amount  *int64   // VND; nil = static QR
	Purpose *string  // remark
	Service *Service // default ServiceTransferToAccount

	// Field 62 sub-fields — nil omits.
	BillNumber    *string
	MobileNumber  *string
	Store         *string
	LoyaltyNumber *string
	Reference     *string
	CustomerLabel *string
	Terminal      *string
}

// QR is a decoded/encoded Vietnam QR payment.
type QR struct {
	Provider string
	Dynamic  bool

	BIN           constants.BankBIN
	AccountNumber string
	Service       Service

	MerchantID   string
	MerchantName string

	Amount int64

	Purpose       string
	BillNumber    string
	MobileNumber  string
	Store         string
	LoyaltyNumber string
	Reference     string
	CustomerLabel string
	Terminal      string

	UnreservedFields map[string]string

	inner *serve.Pay
}

// New builds a VietQR from opts. BIN and AccountNumber are required.
func New(opts Options) *QR {
	service := ServiceTransferToAccount
	if opts.Service != nil {
		service = *opts.Service
	}

	amountVal := derefInt64(opts.Amount)
	amountStr := ""
	if amountVal > 0 {
		amountStr = strconv.FormatInt(amountVal, 10)
	}

	purpose := derefStr(opts.Purpose)

	inner := serve.InitVietQR(serve.InitVietQROptions{
		BankBin:    opts.BIN.String(),
		BankNumber: opts.AccountNumber,
		Amount:     amountStr,
		Purpose:    purpose,
		Service:    string(service),
	})

	inner.AdditionalData.BillNumber = derefStr(opts.BillNumber)
	inner.AdditionalData.MobileNumber = derefStr(opts.MobileNumber)
	inner.AdditionalData.Store = derefStr(opts.Store)
	inner.AdditionalData.LoyaltyNumber = derefStr(opts.LoyaltyNumber)
	inner.AdditionalData.Reference = derefStr(opts.Reference)
	inner.AdditionalData.CustomerLabel = derefStr(opts.CustomerLabel)
	inner.AdditionalData.Terminal = derefStr(opts.Terminal)

	return &QR{
		Provider:      string(constants.QRProviderVIETQR),
		Dynamic:       amountVal > 0,
		BIN:           opts.BIN,
		AccountNumber: opts.AccountNumber,
		Service:       service,
		Amount:        amountVal,
		Purpose:       purpose,
		BillNumber:    derefStr(opts.BillNumber),
		MobileNumber:  derefStr(opts.MobileNumber),
		Store:         derefStr(opts.Store),
		LoyaltyNumber: derefStr(opts.LoyaltyNumber),
		Reference:     derefStr(opts.Reference),
		CustomerLabel: derefStr(opts.CustomerLabel),
		Terminal:      derefStr(opts.Terminal),
		inner:         inner,
	}
}

// VNPayOptions for building a VNPAY merchant QR. MerchantID required.
type VNPayOptions struct {
	// Required.
	MerchantID string

	// Optional.
	MerchantName  *string
	Store         *string
	Terminal      *string
	Amount        *int64
	Purpose       *string
	BillNumber    *string
	MobileNumber  *string
	LoyaltyNumber *string
	Reference     *string
	CustomerLabel *string
}

// NewVNPay builds a VNPAY merchant QR.
func NewVNPay(opts VNPayOptions) *QR {
	req := serve.VNPaymentRequest{
		MerchantID:   opts.MerchantID,
		MerchantName: derefStr(opts.MerchantName),
		Store:        derefStr(opts.Store),
		Terminal:     derefStr(opts.Terminal),
	}
	if opts.Amount != nil && *opts.Amount > 0 {
		s := strconv.FormatInt(*opts.Amount, 10)
		req.Amount = &s
	}
	req.Purpose = opts.Purpose
	req.BillNumber = opts.BillNumber
	req.MobileNumber = opts.MobileNumber
	req.LoyaltyNumber = opts.LoyaltyNumber
	req.Reference = opts.Reference
	req.CustomerLabel = opts.CustomerLabel

	inner := serve.InitVNPayQR(req)
	return &QR{
		Provider:      string(constants.QRProviderVNPAY),
		Dynamic:       opts.Amount != nil && *opts.Amount > 0,
		MerchantID:    opts.MerchantID,
		MerchantName:  derefStr(opts.MerchantName),
		Amount:        derefInt64(opts.Amount),
		Purpose:       derefStr(opts.Purpose),
		BillNumber:    derefStr(opts.BillNumber),
		MobileNumber:  derefStr(opts.MobileNumber),
		Store:         derefStr(opts.Store),
		Terminal:      derefStr(opts.Terminal),
		Reference:     derefStr(opts.Reference),
		LoyaltyNumber: derefStr(opts.LoyaltyNumber),
		CustomerLabel: derefStr(opts.CustomerLabel),
		inner:         inner,
	}
}

// ErrInvalid is returned when CRC fails or content is malformed.
var ErrInvalid = errors.New("bankqr: invalid QR content")

// Parse decodes an EMVCo QR string.
func Parse(content string) (*QR, error) {
	inner := serve.NewQRPay(content)
	if !inner.IsValid {
		return nil, ErrInvalid
	}

	q := &QR{
		Provider:         string(inner.Provider.Name),
		Dynamic:          inner.InitMethod == "12",
		BIN:              constants.BankBIN(inner.Consumer.BankBin),
		AccountNumber:    inner.Consumer.BankNumber,
		Service:          Service(inner.Provider.Service),
		MerchantID:       inner.Merchant.Id,
		MerchantName:     inner.Merchant.Name,
		Purpose:          inner.AdditionalData.Purpose,
		BillNumber:       inner.AdditionalData.BillNumber,
		MobileNumber:     inner.AdditionalData.MobileNumber,
		Store:            inner.AdditionalData.Store,
		LoyaltyNumber:    inner.AdditionalData.LoyaltyNumber,
		Reference:        inner.AdditionalData.Reference,
		CustomerLabel:    inner.AdditionalData.CustomerLabel,
		Terminal:         inner.AdditionalData.Terminal,
		UnreservedFields: inner.Unreserved,
		inner:            inner,
	}
	if inner.Amount != "" {
		if n, err := strconv.ParseInt(inner.Amount, 10, 64); err == nil {
			q.Amount = n
		}
	}
	return q, nil
}

// String renders the QR as an EMVCo payload (CRC included).
func (q *QR) String() string {
	if q.inner == nil {
		q.inner = serve.InitVietQR(serve.InitVietQROptions{
			BankBin:    q.BIN.String(),
			BankNumber: q.AccountNumber,
			Service:    string(q.Service),
		})
	}

	switch q.Provider {
	case string(constants.QRProviderVNPAY):
		q.inner.Merchant.Id = q.MerchantID
		q.inner.Merchant.Name = q.MerchantName
	default:
		q.inner.Consumer.BankBin = q.BIN.String()
		q.inner.Consumer.BankNumber = q.AccountNumber
		if q.Service != "" {
			q.inner.Provider.Service = string(q.Service)
		}
	}
	if q.Amount > 0 {
		q.inner.Amount = strconv.FormatInt(q.Amount, 10)
		q.inner.InitMethod = "12"
	} else {
		q.inner.Amount = ""
		q.inner.InitMethod = "11"
	}
	q.inner.AdditionalData.Purpose = q.Purpose
	q.inner.AdditionalData.BillNumber = q.BillNumber
	q.inner.AdditionalData.MobileNumber = q.MobileNumber
	q.inner.AdditionalData.Store = q.Store
	q.inner.AdditionalData.LoyaltyNumber = q.LoyaltyNumber
	q.inner.AdditionalData.Reference = q.Reference
	q.inner.AdditionalData.CustomerLabel = q.CustomerLabel
	q.inner.AdditionalData.Terminal = q.Terminal
	if q.UnreservedFields != nil {
		q.inner.Unreserved = q.UnreservedFields
	}

	return q.inner.Build()
}

// Bank looks up a bank by Napas BIN.
func Bank(bin constants.BankBIN) (constants.Bank, bool) {
	s := bin.String()
	for _, b := range constants.BanksMap {
		if b.BIN == s {
			return b, true
		}
	}
	return constants.Bank{}, false
}

// MustNew is like New but panics if the result fails to round-trip.
func MustNew(opts Options) *QR {
	qr := New(opts)
	if _, err := Parse(qr.String()); err != nil {
		panic(fmt.Sprintf("bankqr.MustNew: %v", err))
	}
	return qr
}

// Tiny pointer-deref helpers — internal only.

func derefStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func derefInt64(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}
