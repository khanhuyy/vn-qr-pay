package test

import (
	"testing"

	"github.com/khanhuyy/emvgo/constants"
	"github.com/khanhuyy/emvgo/serve"
)

// Parse static VietQR.
func TestParse_VietQR_Static_ACB(t *testing.T) {
	qr := serve.NewQRPay(vietQRStaticACB)
	if !qr.IsValid {
		t.Fatal("expected valid QR")
	}
	if qr.Version != "01" {
		t.Errorf("Version: want 01, got %q", qr.Version)
	}
	if qr.InitMethod != "11" {
		t.Errorf("InitMethod: want 11, got %q", qr.InitMethod)
	}
	if qr.Provider.Name != constants.QRProviderVIETQR {
		t.Errorf("Provider.Name: want VIETQR, got %q", qr.Provider.Name)
	}
	if qr.Provider.GUID != "A000000727" {
		t.Errorf("Provider.GUID: want A000000727, got %q", qr.Provider.GUID)
	}
	if qr.Provider.Service != string(constants.VietQRServiceByAccount) {
		t.Errorf("Provider.Service: want QRIBFTTA, got %q", qr.Provider.Service)
	}
	if qr.Consumer.BankBin != "970416" {
		t.Errorf("BankBin: want 970416, got %q", qr.Consumer.BankBin)
	}
	if qr.Consumer.BankNumber != "257678859" {
		t.Errorf("BankNumber: want 257678859, got %q", qr.Consumer.BankNumber)
	}
	if qr.Amount != "" {
		t.Errorf("static QR should have empty Amount, got %q", qr.Amount)
	}
}

// Parse dynamic VietQR.
func TestParse_VietQR_Dynamic_ACB(t *testing.T) {
	qr := serve.NewQRPay(vietQRDynamicACB)
	if !qr.IsValid {
		t.Fatal("expected valid QR")
	}
	if qr.InitMethod != "12" {
		t.Errorf("InitMethod: want 12, got %q", qr.InitMethod)
	}
	if qr.Amount != "10000" {
		t.Errorf("Amount: want 10000, got %q", qr.Amount)
	}
	if qr.AdditionalData.Purpose != "Chuyen tien" {
		t.Errorf("Purpose: want 'Chuyen tien', got %q", qr.AdditionalData.Purpose)
	}
}

// Parse MoMo (unreserved tag 80 + MOMOW2W ref).
func TestParse_MoMo_ExtractsUnreservedAndReference(t *testing.T) {
	qr := serve.NewQRPay(momoQR)
	if !qr.IsValid {
		t.Fatal("expected valid MoMo QR")
	}
	if qr.Consumer.BankBin != "970454" {
		t.Errorf("BankBin: want 970454 (BVBank), got %q", qr.Consumer.BankBin)
	}
	if qr.Consumer.BankNumber != "99MM24011M34875080" {
		t.Errorf("BankNumber: got %q", qr.Consumer.BankNumber)
	}
	if got := qr.AdditionalData.Reference; got != "MOMOW2W34875080" {
		t.Errorf("Reference: want MOMOW2W34875080, got %q", got)
	}
	if got, ok := qr.Unreserved["80"]; !ok || got != "046" {
		t.Errorf("Unreserved[80]: want 046, got %q (ok=%v)", got, ok)
	}
}

// Parse ZaloPay.
func TestParse_ZaloPay(t *testing.T) {
	qr := serve.NewQRPay(zaloPayQR)
	if !qr.IsValid {
		t.Fatal("expected valid ZaloPay QR")
	}
	if qr.Consumer.BankNumber != "99ZP24009M07248267" {
		t.Errorf("BankNumber: got %q", qr.Consumer.BankNumber)
	}
}

// Bad CRC must invalidate.
func TestParse_BadCRC_InvalidatesQR(t *testing.T) {
	// flip the last hex char so the CRC no longer matches
	bad := vietQRStaticACB[:len(vietQRStaticACB)-1] + "0"
	qr := serve.NewQRPay(bad)
	if qr.IsValid {
		t.Error("QR with mutated CRC should be invalid")
	}
}

func TestParse_TruncatedContent_InvalidatesQR(t *testing.T) {
	qr := serve.NewQRPay("00")
	if qr.IsValid {
		t.Error("truncated content should be invalid")
	}
}

// Round-trip stability.
func TestRoundTrip_VietQRDynamic(t *testing.T) {
	first := serve.NewQRPay(vietQRDynamicACB)
	if !first.IsValid {
		t.Fatal("seed QR invalid")
	}

	rebuilt := first.Build()
	if rebuilt != vietQRDynamicACB {
		t.Errorf("rebuild not byte-identical\nwant: %s\n got: %s", vietQRDynamicACB, rebuilt)
	}

	second := serve.NewQRPay(rebuilt)
	if !second.IsValid {
		t.Fatal("rebuilt QR failed to re-parse")
	}
	if second.Consumer.BankNumber != first.Consumer.BankNumber {
		t.Errorf("BankNumber drifted: %q vs %q", first.Consumer.BankNumber, second.Consumer.BankNumber)
	}
	if second.Amount != first.Amount {
		t.Errorf("Amount drifted: %q vs %q", first.Amount, second.Amount)
	}
}

func TestRoundTrip_MoMo_PreservesUnreserved(t *testing.T) {
	first := serve.NewQRPay(momoQR)
	if !first.IsValid {
		t.Fatal("seed MoMo QR invalid")
	}
	rebuilt := first.Build()
	second := serve.NewQRPay(rebuilt)
	if !second.IsValid {
		t.Fatal("rebuilt MoMo QR failed to re-parse")
	}
	if got := second.Unreserved["80"]; got != "046" {
		t.Errorf("Unreserved[80] not preserved on round trip: got %q", got)
	}
	if second.AdditionalData.Reference != "MOMOW2W34875080" {
		t.Errorf("Reference not preserved: got %q", second.AdditionalData.Reference)
	}
}

// Mutate amount/purpose then rebuild.
func TestMutate_ChangeAmountAndPurpose(t *testing.T) {
	qr := serve.NewQRPay(vietQRDynamicACB)
	qr.Amount = "999999"
	qr.AdditionalData.Purpose = "Cam on nhe"

	const want = "00020101021238530010A0000007270123000697041601092576788590208QRIBFTTA530370454069999995802VN62140810Cam on nhe6304E786"
	if got := qr.Build(); got != want {
		t.Errorf("mutated rebuild mismatch\nwant: %s\n got: %s", want, got)
	}
}
