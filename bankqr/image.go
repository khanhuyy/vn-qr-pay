package bankqr

import (
	"encoding/base64"
	"fmt"

	qrcode "github.com/skip2/go-qrcode"
)

// ErrorCorrection level for QR redundancy (~7/15/25/30%).
type ErrorCorrection int

const (
	ErrorCorrectionLow ErrorCorrection = iota
	ErrorCorrectionMedium
	ErrorCorrectionHigh
	ErrorCorrectionHighest
)

func (e ErrorCorrection) toQR() qrcode.RecoveryLevel {
	switch e {
	case ErrorCorrectionLow:
		return qrcode.Low
	case ErrorCorrectionHigh:
		return qrcode.High
	case ErrorCorrectionHighest:
		return qrcode.Highest
	default:
		return qrcode.Medium
	}
}

// PNG renders the QR as PNG bytes. size = pixel side length.
func (q *QR) PNG(size int) ([]byte, error) {
	return q.PNGWithLevel(size, ErrorCorrectionMedium)
}

// PNGWithLevel renders with an explicit correction level.
func (q *QR) PNGWithLevel(size int, level ErrorCorrection) ([]byte, error) {
	if q == nil {
		return nil, fmt.Errorf("bankqr: PNG called on nil *QR")
	}
	content := q.String()
	if content == "" {
		return nil, fmt.Errorf("bankqr: cannot render PNG of empty QR")
	}
	return qrcode.Encode(content, level.toQR(), size)
}

// PNGBase64 returns base64-encoded PNG (no data URL prefix).
func (q *QR) PNGBase64(size int) (string, error) {
	return q.PNGBase64WithLevel(size, ErrorCorrectionMedium)
}

func (q *QR) PNGBase64WithLevel(size int, level ErrorCorrection) (string, error) {
	raw, err := q.PNGWithLevel(size, level)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(raw), nil
}

// PNGDataURL returns a data:image/png;base64,... URL.
func (q *QR) PNGDataURL(size int) (string, error) {
	return q.PNGDataURLWithLevel(size, ErrorCorrectionMedium)
}

func (q *QR) PNGDataURLWithLevel(size int, level ErrorCorrection) (string, error) {
	b64, err := q.PNGBase64WithLevel(size, level)
	if err != nil {
		return "", err
	}
	return "data:image/png;base64," + b64, nil
}
