package test

import (
	"fmt"
	"testing"

	"github.com/khanhuyy/emvqr/serve"
)

// CRC-16/CCITT-FALSE vectors (incl. real QR strings minus the 4-hex CRC).
func TestCRC16CCITT_KnownVectors(t *testing.T) {
	cases := []struct {
		in   string
		want uint16
	}{
		// Reference vector from the CRC catalogue.
		{"123456789", 0x29B1},

		// The CRC input for the static ACB VietQR is everything up to and
		// including the "6304" tag+length prefix of field 63.
		{vietQRStaticACB[:len(vietQRStaticACB)-4], 0xAE9F},
		{vietQRDynamicACB[:len(vietQRDynamicACB)-4], 0x53E6},
		{momoQR[:len(momoQR)-4], 0xEBC8},
		{zaloPayQR[:len(zaloPayQR)-4], 0x073C},
	}
	for _, c := range cases {
		if got := serve.CRC16CCITT(c.in); got != c.want {
			t.Errorf("CRC16CCITT(%q): want %04X, got %04X", c.in, c.want, got)
		}
	}
}

// Every Build() output must self-verify CRC.
func TestCRC_SelfConsistent_ForAllInitVietQRCombinations(t *testing.T) {
	bins := []string{"970436", "970422", "970416", "970407", "970454"}
	for _, bin := range bins {
		for _, amount := range []string{"", "1000", "9999999"} {
			label := fmt.Sprintf("bin=%s amount=%q", bin, amount)
			t.Run(label, func(t *testing.T) {
				out := serve.InitVietQR(serve.InitVietQROptions{
					BankBin:    bin,
					BankNumber: "1234567890",
					Amount:     amount,
				}).Build()
				if !serve.NewQRPay(out).IsValid {
					t.Errorf("self-built QR fails verifyCRC: %s", out)
				}
			})
		}
	}
}
