package test

import (
	"qrpay/serve"
	"testing"
)

func TestBuild(t *testing.T) {
	serve.NewQRPay()
	result := serve.Build(2, 3)
	expected := 5
	if result != expected {
		t.Errorf("Add(2, 3) = %d; want %d", result, expected)
	}
}
