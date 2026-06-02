# emvqr

Go library for Vietnam payment QR codes — **VietQR** (Napas 247) and **VNPAYQR**.
Encodes and decodes per the EMVCo spec; reference fixtures verified against
`xuannghia/vietnam-qr-pay`.

## Install

```bash
go get github.com/khanhuyy/emvgo
```

Minimum Go version: **1.13**.

## Library API — `github.com/khanhuyy/emvgo/bankqr`

The headline package:

```go
import (
    "github.com/khanhuyy/emvgo/bankqr"
    "github.com/khanhuyy/emvgo/constants"
    "github.com/khanhuyy/emvgo/utils"
)

// Static QR (no amount)
qr := bankqr.New(bankqr.Options{
    BIN:           constants.BIN_VIETCOMBANK,
    AccountNumber: "1234567890",
})
fmt.Println(qr.String())

// Dynamic QR (with amount + remark) — optional fields are pointers
qr = bankqr.New(bankqr.Options{
    BIN:           constants.BIN_ACB,
    AccountNumber: "257678859",
    Amount:        utils.Int64Ptr(10_000),
    Purpose:       utils.StrPtr("Chuyen tien"),
})
fmt.Println(qr.String())

// Parse
parsed, err := bankqr.Parse(content)
if err != nil { /* errors.Is(err, bankqr.ErrInvalid) */ }
fmt.Println(parsed.BIN, parsed.AccountNumber, parsed.Amount, parsed.Purpose)

// Mutate then rebuild — CRC is recomputed automatically
parsed.Amount = 999_999
parsed.Purpose = "Cam on nhe"
updated := parsed.String()

// VNPAY merchant QR
vn := bankqr.NewVNPay(bankqr.VNPayOptions{
    MerchantID:   "0102154778",
    MerchantName: utils.StrPtr("TUGIACOMPANY"),
    Store:        utils.StrPtr("TU GIA COMPUTER"),
    Terminal:     utils.StrPtr("TUGIACO1"),
})

// Look up bank by BIN enum
b, ok := bankqr.Bank(constants.BIN_VIETCOMBANK) // b.ShortName == "Vietcombank"
```

### Options layout

`bankqr.Options` follows a "required = value, optional = pointer" convention:

| Field | Type | Required |
|-------|------|----------|
| `BIN` | `constants.BankBIN` | ✅ |
| `AccountNumber` | `string` | ✅ |
| `Amount` | `*int64` | optional (nil = static QR) |
| `Purpose` | `*string` | optional |
| `Service` | `*Service` | optional (default `ServiceTransferToAccount`) |
| `BillNumber`, `MobileNumber`, `Store`, `LoyaltyNumber`, `Reference`, `CustomerLabel`, `Terminal` | `*string` | optional |

Pointer helpers live in `github.com/khanhuyy/emvgo/utils`: `StrPtr`, `IntPtr`,
`Int64Ptr`, `BoolPtr`.

### Bank enums

`constants/bank_bin.go` exposes ~50 well-known `BankBIN` constants
(`BIN_VIETCOMBANK = "970436"`, `BIN_TECHCOMBANK = "970407"`, `BIN_MBBANK`,
`BIN_ACB`, `BIN_BIDV`, `BIN_BVBANK`, …). `BankCode` (full names) and `BankKey`
(URL slugs) both have `.String()` helpers.

Service codes:

| Constant | Code | Use when |
|----------|------|----------|
| `bankqr.ServiceTransferToAccount` | `QRIBFTTA` | Transfer to account (default) |
| `bankqr.ServiceTransferToCard`    | `QRIBFTTC` | Transfer to card |

## QR image generation

```go
qr := bankqr.New(bankqr.Options{
    BIN:           constants.BIN_VIETCOMBANK,
    AccountNumber: "1234567890",
})

// Raw PNG bytes
png, err := qr.PNG(256) // 256x256 px

// Base64-encoded PNG (for JSON / DB / chat)
b64, err := qr.PNGBase64(256)

// Data URL (drop straight into <img src>)
url, err := qr.PNGDataURL(256)
// "data:image/png;base64,iVBORw0KGgo..."

// Custom error correction level
b64, err := qr.PNGBase64WithLevel(256, bankqr.ErrorCorrectionHigh)
```

Levels: `ErrorCorrectionLow` / `Medium` / `High` / `Highest` (≈7/15/25/30%
redundancy). Medium is the default; standard for payment QRs.

## CLI — `cmd/qrpay`

```bash
make build           # → ./bin/qrpay
make install         # → $GOPATH/bin/qrpay

qrpay build  --bin 970436 --account 1234567890 --amount 50000 --purpose "Cam on"
qrpay parse  "<qr-content>"
qrpay crc    123456789
qrpay banks  --filter techcom
```

Append `--json` to `build` / `parse` / `banks` for pipeline-friendly output.

## Lower-level package — `github.com/khanhuyy/emvgo/serve`

`bankqr` wraps `serve.Pay`. Reach into `qrpay/serve` directly if you need raw
EMVCo control (custom field 65-79, unreserved 80-99, tip & fee, multi-language).

## Spec docs

See [`docs/QR_SPEC.md`](docs/QR_SPEC.md) for the full specification: EMVCo
Merchant-Presented structure, field 26 (VNPAY) / 38 (VietQR), additional data
field 62, CRC-16/CCITT-FALSE, how MoMo and ZaloPay tunnel through BVBank, and
a table of the most popular banks.

## Testing

```bash
make test    # all unit + integration tests
make cover   # coverage profile (serve + constants)
make demo    # 3 end-to-end CLI examples
```
