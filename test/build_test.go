package test

import (
	"strings"
	"testing"

	"github.com/khanhuyy/emvgo/constants"
	"github.com/khanhuyy/emvgo/serve"
)

// Empty constructor must produce a valid, buildable instance.
func TestBuild_EmptyConstructorStillValid(t *testing.T) {
	qr := serve.NewQRPay("")
	qr.Merchant.Name = "Test Merchant"
	qr.City = "Hanoi"

	if !qr.IsValid {
		t.Error("NewQRPay(\"\") should produce a valid (buildable) instance")
	}
	if got := qr.Build(); got == "" {
		t.Error("Build() returned empty string")
	}
}

// ACB static fixture from xuannghia/vietnam-qr-pay.
const vietQRStaticACB = "00020101021138530010A0000007270123000697041601092576788590208QRIBFTTA53037045802VN6304AE9F"

func TestBuild_VietQR_Static_ACB(t *testing.T) {
	qr := serve.InitVietQR(serve.InitVietQROptions{
		BankBin:    "970416",
		BankNumber: "257678859",
	})
	got := qr.Build()
	if got != vietQRStaticACB {
		t.Fatalf("VietQR static mismatch\nwant: %s\n got: %s", vietQRStaticACB, got)
	}
	if qr.InitMethod != "11" {
		t.Errorf("static QR must use initMethod=11, got %q", qr.InitMethod)
	}
}

// ACB dynamic fixture (amount + purpose).
const vietQRDynamicACB = "00020101021238530010A0000007270123000697041601092576788590208QRIBFTTA53037045405100005802VN62150811Chuyen tien630453E6"

func TestBuild_VietQR_Dynamic_ACB(t *testing.T) {
	qr := serve.InitVietQR(serve.InitVietQROptions{
		BankBin:    "970416",
		BankNumber: "257678859",
		Amount:     "10000",
		Purpose:    "Chuyen tien",
	})
	got := qr.Build()
	if got != vietQRDynamicACB {
		t.Fatalf("VietQR dynamic mismatch\nwant: %s\n got: %s", vietQRDynamicACB, got)
	}
	if qr.InitMethod != "12" {
		t.Errorf("dynamic QR must use initMethod=12, got %q", qr.InitMethod)
	}
}

// MoMo via BVBank (tag-80 + MOMOW2W reference).
const momoQR = "00020101021138620010A00000072701320006970454011899MM24011M348750800208QRIBFTTA53037045802VN62190515MOMOW2W3487508080030466304EBC8"

func TestBuild_MoMo_via_BVBank(t *testing.T) {
	const acc = "99MM24011M34875080"

	qr := serve.InitVietQR(serve.InitVietQROptions{
		BankBin:    "970454", // BanViet (BVBank)
		BankNumber: acc,
	})
	qr.AdditionalData.Reference = "MOMOW2W" + acc[10:] // "MOMOW2W34875080"
	qr.SetUnreservedField("80", "046")

	got := qr.Build()
	if got != momoQR {
		t.Fatalf("MoMo QR mismatch\nwant: %s\n got: %s", momoQR, got)
	}
}

// ZaloPay via BVBank.
const zaloPayQR = "00020101021138620010A00000072701320006970454011899ZP24009M072482670208QRIBFTTA53037045802VN6304073C"

func TestBuild_ZaloPay_via_BVBank(t *testing.T) {
	qr := serve.InitVietQR(serve.InitVietQROptions{
		BankBin:    "970454",
		BankNumber: "99ZP24009M07248267",
	})
	got := qr.Build()
	if got != zaloPayQR {
		t.Fatalf("ZaloPay QR mismatch\nwant: %s\n got: %s", zaloPayQR, got)
	}
}

// Transfer-to-card variant.
func TestBuild_VietQR_TransferToCard(t *testing.T) {
	qr := serve.InitVietQR(serve.InitVietQROptions{
		BankBin:    "970422", // MB Bank
		BankNumber: "9704229876543210",
		Service:    string(constants.VietQRServiceByCard),
	})
	out := qr.Build()
	if !strings.Contains(out, "0208QRIBFTTC") {
		t.Errorf("expected QRIBFTTC service tag, got %s", out)
	}
}

// Shorthand instance.Build() must produce a parseable, CRC-valid QR.
func TestBuild_Instance_Shorthand(t *testing.T) {
	out := serve.NewMyQRPay().Build(serve.BuildQROptions{
		BankBin:    "970407", // Techcombank
		BankNumber: "19033868110065",
		Amount:     "100000",
		Remark:     "test",
	})

	parsed := serve.NewQRPay(out)
	if !parsed.IsValid {
		t.Fatalf("shorthand output failed CRC: %s", out)
	}
	if parsed.Consumer.BankBin != "970407" {
		t.Errorf("bankBin: want 970407, got %q", parsed.Consumer.BankBin)
	}
	if parsed.Consumer.BankNumber != "19033868110065" {
		t.Errorf("bankNumber: want 19033868110065, got %q", parsed.Consumer.BankNumber)
	}
	if parsed.Amount != "100000" {
		t.Errorf("amount: want 100000, got %q", parsed.Amount)
	}
	if parsed.AdditionalData.Purpose != "test" {
		t.Errorf("purpose: want test, got %q", parsed.AdditionalData.Purpose)
	}
	if parsed.InitMethod != "12" {
		t.Errorf("dynamic QR should have initMethod=12, got %q", parsed.InitMethod)
	}
}

// CRC must be 4 uppercase hex chars after field "6304".
func TestBuild_CRC_IsAppendedAsUppercaseHex(t *testing.T) {
	qr := serve.InitVietQR(serve.InitVietQROptions{
		BankBin:    "970436", // Vietcombank
		BankNumber: "1234567890",
	})
	out := qr.Build()
	if len(out) < 8 {
		t.Fatalf("output too short: %s", out)
	}
	crc := out[len(out)-4:]
	if crc != strings.ToUpper(crc) {
		t.Errorf("CRC must be uppercase hex, got %s", crc)
	}
	for _, c := range crc {
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F')) {
			t.Errorf("CRC contains non-hex char: %s", crc)
			break
		}
	}
	if !strings.Contains(out, "6304"+crc) {
		t.Errorf("expected output to end with %q field, got %s", "6304"+crc, out)
	}
}

// Default currency = 704 (VND), nation = VN.
func TestBuild_DefaultsCurrencyAndNation(t *testing.T) {
	out := serve.InitVietQR(serve.InitVietQROptions{
		BankBin:    "970436",
		BankNumber: "1234567890",
	}).Build()

	if !strings.Contains(out, "5303704") {
		t.Errorf("missing default currency 704 (VND): %s", out)
	}
	if !strings.Contains(out, "5802VN") {
		t.Errorf("missing default country VN: %s", out)
	}
}
