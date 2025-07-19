package test

import (
	"qrpay/qr"
	"testing"
)

func TestBuild(t *testing.T) {
	qr.NewQRPay()
	result := qr.Build(2, 3)
	expected := 5
	if result != expected {
		t.Errorf("Add(2, 3) = %d; want %d", result, expected)
	}
}
