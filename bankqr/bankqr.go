// Package bankqr builds and parses Vietnam bank QR codes (VietQR + VNPAYQR).
package bankqr

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/khanhuyy/emvqr/constants"
	"github.com/khanhuyy/emvqr/serve"
)

// Service is the VietQR transfer service code (field 38.02).
type Service string

const (
	ServiceTransferToAccount Service = "QRIBFTTA" // transfer to account
	ServiceTransferToCard    Service = "QRIBFTTC" // transfer to card
)

// Options for building a VietQR.
type Options struct {
	BIN           string // bank BIN, e.g. 970436
	AccountNumber string
	Amount        int64 // 0 = static QR
	Purpose       string
	Service       Service // default QRIBFTTA

	// Field 62 sub-fields.
	BillNumber    string
	MobileNumber  string
	Store         string
	LoyaltyNumber string
	Reference     string
	CustomerLabel string
	Terminal      string
}

// QR is a decoded/encoded Vietnam QR payment.
type QR struct {
	Provider string
	Dynamic  bool

	BIN           string
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

// New builds a VietQR from opts.
func New(opts Options) *QR {
	service := opts.Service
	if service == "" {
		service = ServiceTransferToAccount
	}

	amountStr := ""
	if opts.Amount > 0 {
		amountStr = strconv.FormatInt(opts.Amount, 10)
	}

	inner := serve.InitVietQR(serve.InitVietQROptions{
		BankBin:    opts.BIN,
		BankNumber: opts.AccountNumber,
		Amount:     amountStr,
		Purpose:    opts.Purpose,
		Service:    string(service),
	})

	inner.AdditionalData.BillNumber = opts.BillNumber
	inner.AdditionalData.MobileNumber = opts.MobileNumber
	inner.AdditionalData.Store = opts.Store
	inner.AdditionalData.LoyaltyNumber = opts.LoyaltyNumber
	inner.AdditionalData.Reference = opts.Reference
	inner.AdditionalData.CustomerLabel = opts.CustomerLabel
	inner.AdditionalData.Terminal = opts.Terminal

	return &QR{
		Provider:      string(constants.QRProviderVIETQR),
		Dynamic:       opts.Amount > 0,
		BIN:           opts.BIN,
		AccountNumber: opts.AccountNumber,
		Service:       service,
		Amount:        opts.Amount,
		Purpose:       opts.Purpose,
		BillNumber:    opts.BillNumber,
		MobileNumber:  opts.MobileNumber,
		Store:         opts.Store,
		LoyaltyNumber: opts.LoyaltyNumber,
		Reference:     opts.Reference,
		CustomerLabel: opts.CustomerLabel,
		Terminal:      opts.Terminal,
		inner:         inner,
	}
}

// VNPayOptions for building a VNPAY merchant QR.
type VNPayOptions struct {
	MerchantID   string
	MerchantName string
	Store        string
	Terminal     string

	Amount        int64
	Purpose       string
	BillNumber    string
	MobileNumber  string
	LoyaltyNumber string
	Reference     string
	CustomerLabel string
}

// NewVNPay builds a VNPAY merchant QR.
func NewVNPay(opts VNPayOptions) *QR {
	req := serve.VNPaymentRequest{
		MerchantID:   opts.MerchantID,
		MerchantName: opts.MerchantName,
		Store:        opts.Store,
		Terminal:     opts.Terminal,
	}
	if opts.Amount > 0 {
		s := strconv.FormatInt(opts.Amount, 10)
		req.Amount = &s
	}
	if opts.Purpose != "" {
		req.Purpose = &opts.Purpose
	}
	if opts.BillNumber != "" {
		req.BillNumber = &opts.BillNumber
	}
	if opts.MobileNumber != "" {
		req.MobileNumber = &opts.MobileNumber
	}
	if opts.LoyaltyNumber != "" {
		req.LoyaltyNumber = &opts.LoyaltyNumber
	}
	if opts.Reference != "" {
		req.Reference = &opts.Reference
	}
	if opts.CustomerLabel != "" {
		req.CustomerLabel = &opts.CustomerLabel
	}

	inner := serve.InitVNPayQR(req)
	return &QR{
		Provider:      string(constants.QRProviderVNPAY),
		Dynamic:       opts.Amount > 0,
		MerchantID:    opts.MerchantID,
		MerchantName:  opts.MerchantName,
		Amount:        opts.Amount,
		Purpose:       opts.Purpose,
		BillNumber:    opts.BillNumber,
		MobileNumber:  opts.MobileNumber,
		Store:         opts.Store,
		Terminal:      opts.Terminal,
		Reference:     opts.Reference,
		LoyaltyNumber: opts.LoyaltyNumber,
		CustomerLabel: opts.CustomerLabel,
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
		BIN:              inner.Consumer.BankBin,
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
			BankBin:    q.BIN,
			BankNumber: q.AccountNumber,
			Service:    string(q.Service),
		})
	}

	switch q.Provider {
	case string(constants.QRProviderVNPAY):
		q.inner.Merchant.Id = q.MerchantID
		q.inner.Merchant.Name = q.MerchantName
	default:
		q.inner.Consumer.BankBin = q.BIN
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

// Bank looks up a bank by 6-digit Napas BIN.
func Bank(bin string) (constants.Bank, bool) {
	for _, b := range constants.BanksMap {
		if b.BIN == bin {
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
