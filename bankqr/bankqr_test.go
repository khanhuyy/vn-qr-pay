package bankqr_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/khanhuyy/emvgo/bankqr"
	"github.com/khanhuyy/emvgo/constants"
	"github.com/khanhuyy/emvgo/utils"
)

// Reference fixtures.
const (
	vietQRStaticACB  = "00020101021138530010A0000007270123000697041601092576788590208QRIBFTTA53037045802VN6304AE9F"
	vietQRDynamicACB = "00020101021238530010A0000007270123000697041601092576788590208QRIBFTTA53037045405100005802VN62150811Chuyen tien630453E6"
	momoQR           = "00020101021138620010A00000072701320006970454011899MM24011M348750800208QRIBFTTA53037045802VN62190515MOMOW2W3487508080030466304EBC8"
)

func TestNew_StaticVietQR_MatchesReference(t *testing.T) {
	qr := bankqr.New(bankqr.Options{
		BIN:           constants.BIN_ACB,
		AccountNumber: "257678859",
	})
	if got := qr.String(); got != vietQRStaticACB {
		t.Fatalf("static VietQR mismatch\nwant: %s\n got: %s", vietQRStaticACB, got)
	}
	if qr.Dynamic {
		t.Error("static QR should have Dynamic=false")
	}
}

func TestNew_DynamicVietQR_MatchesReference(t *testing.T) {
	qr := bankqr.New(bankqr.Options{
		BIN:           constants.BIN_ACB,
		AccountNumber: "257678859",
		Amount:        utils.Int64Ptr(10_000),
		Purpose:       utils.StrPtr("Chuyen tien"),
	})
	if got := qr.String(); got != vietQRDynamicACB {
		t.Fatalf("dynamic VietQR mismatch\nwant: %s\n got: %s", vietQRDynamicACB, got)
	}
	if !qr.Dynamic {
		t.Error("amount-bearing QR should have Dynamic=true")
	}
}

func TestNew_DefaultServiceIsTransferToAccount(t *testing.T) {
	qr := bankqr.New(bankqr.Options{BIN: constants.BIN_VIETCOMBANK, AccountNumber: "1"})
	if qr.Service != bankqr.ServiceTransferToAccount {
		t.Errorf("default service: want QRIBFTTA, got %q", qr.Service)
	}
}

func TestNew_TransferToCardServiceFlows(t *testing.T) {
	svc := bankqr.ServiceTransferToCard
	qr := bankqr.New(bankqr.Options{
		BIN:           constants.BIN_MBBANK,
		AccountNumber: "9704229876543210",
		Service:       &svc,
	})
	out := qr.String()
	if want := "0208QRIBFTTC"; !contains(out, want) {
		t.Errorf("missing %q in %s", want, out)
	}
}

func TestNew_BinEnumProducesCorrectBytes(t *testing.T) {
	// Verify the constant resolves to the same 6 digits we documented.
	if constants.BIN_VIETCOMBANK.String() != "970436" {
		t.Errorf("BIN_VIETCOMBANK: want 970436, got %q", constants.BIN_VIETCOMBANK.String())
	}
	if constants.BIN_TECHCOMBANK.String() != "970407" {
		t.Errorf("BIN_TECHCOMBANK: want 970407, got %q", constants.BIN_TECHCOMBANK.String())
	}

	out := bankqr.New(bankqr.Options{
		BIN:           constants.BIN_VIETCOMBANK,
		AccountNumber: "1234567890",
	}).String()
	if !contains(out, "0006970436") {
		t.Errorf("BIN 970436 not present in output: %s", out)
	}
}

func TestParse_VietQRDynamic(t *testing.T) {
	qr, err := bankqr.Parse(vietQRDynamicACB)
	if err != nil {
		t.Fatal(err)
	}
	if qr.BIN != constants.BIN_ACB {
		t.Errorf("BIN: want BIN_ACB, got %q", qr.BIN)
	}
	if qr.AccountNumber != "257678859" {
		t.Errorf("AccountNumber: got %q", qr.AccountNumber)
	}
	if qr.Amount != 10_000 {
		t.Errorf("Amount: want 10000, got %d", qr.Amount)
	}
	if qr.Purpose != "Chuyen tien" {
		t.Errorf("Purpose: got %q", qr.Purpose)
	}
	if !qr.Dynamic {
		t.Error("expected Dynamic=true")
	}
}

func TestParse_InvalidCRCReturnsErrInvalid(t *testing.T) {
	bad := vietQRStaticACB[:len(vietQRStaticACB)-1] + "0"
	_, err := bankqr.Parse(bad)
	if !errors.Is(err, bankqr.ErrInvalid) {
		t.Errorf("want ErrInvalid, got %v", err)
	}
}

func TestRoundTrip_MutateAmountAndRebuild(t *testing.T) {
	qr, err := bankqr.Parse(vietQRDynamicACB)
	if err != nil {
		t.Fatal(err)
	}
	qr.Amount = 999_999
	qr.Purpose = "Cam on nhe"

	const want = "00020101021238530010A0000007270123000697041601092576788590208QRIBFTTA530370454069999995802VN62140810Cam on nhe6304E786"
	if got := qr.String(); got != want {
		t.Errorf("rebuild mismatch\nwant: %s\n got: %s", want, got)
	}
}

func TestRoundTrip_MoMoPreservesUnreserved(t *testing.T) {
	qr, err := bankqr.Parse(momoQR)
	if err != nil {
		t.Fatal(err)
	}
	if v := qr.UnreservedFields["80"]; v != "046" {
		t.Errorf("Unreserved[80]: want 046, got %q", v)
	}
	if qr.Reference != "MOMOW2W34875080" {
		t.Errorf("Reference: got %q", qr.Reference)
	}
	if got := qr.String(); got != momoQR {
		t.Errorf("rebuild not byte-identical:\nwant %s\n got %s", momoQR, got)
	}
}

func TestBank_FoundAndMissing(t *testing.T) {
	b, ok := bankqr.Bank(constants.BIN_VIETCOMBANK)
	if !ok {
		t.Fatal("expected to find Vietcombank by BIN")
	}
	if b.ShortName == "" {
		t.Error("ShortName should be populated")
	}

	if _, ok := bankqr.Bank(constants.BankBIN("000000")); ok {
		t.Error("expected miss for fake BIN")
	}
}

func TestBankCode_StringMethod(t *testing.T) {
	if constants.BankCode_VIETCOMBANK.String() != "VIETCOMBANK" {
		t.Errorf("BankCode.String() mismatch: %q", constants.BankCode_VIETCOMBANK.String())
	}
	if constants.BankKey_VIETCOMBANK.String() != "vietcombank" {
		t.Errorf("BankKey.String() mismatch: %q", constants.BankKey_VIETCOMBANK.String())
	}
}

// Godoc examples.
func ExampleNew_static() {
	qr := bankqr.New(bankqr.Options{
		BIN:           constants.BIN_ACB,
		AccountNumber: "257678859",
	})
	fmt.Println(qr.String())
	// Output: 00020101021138530010A0000007270123000697041601092576788590208QRIBFTTA53037045802VN6304AE9F
}

func ExampleNew_dynamic() {
	qr := bankqr.New(bankqr.Options{
		BIN:           constants.BIN_ACB,
		AccountNumber: "257678859",
		Amount:        utils.Int64Ptr(10_000),
		Purpose:       utils.StrPtr("Chuyen tien"),
	})
	fmt.Println(qr.String())
	// Output: 00020101021238530010A0000007270123000697041601092576788590208QRIBFTTA53037045405100005802VN62150811Chuyen tien630453E6
}

func ExampleParse() {
	qr, err := bankqr.Parse("00020101021238530010A0000007270123000697041601092576788590208QRIBFTTA53037045405100005802VN62150811Chuyen tien630453E6")
	if err != nil {
		panic(err)
	}
	fmt.Println(qr.BIN, qr.AccountNumber, qr.Amount, qr.Purpose)
	// Output: 970416 257678859 10000 Chuyen tien
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
