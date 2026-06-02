# Vietnam Payment QR — Specification Summary

This document summarises the QR code formats used for retail payments in Vietnam,
cross-referenced with the sources from Napas, EMVCo, partner banks and e-wallets.
The goal: let a developer understand **what every byte in a QR string means**,
and reconcile the library's behaviour with the spec.

## 1. EMVCo Merchant-Presented structure

Every Vietnamese bank/e-wallet QR is an ASCII string that follows
**EMV® QR Code Specification — Merchant-Presented Mode v1.1**. Every data
object is a `TT LL VV…` triple:

| Field | Description | Length |
|-------|-------------|--------|
| Tag (ID) | 2 digits `00`–`99` | 2 |
| Length   | 2 digits `01`–`99` | 2 |
| Value    | content | equal to Length |

Top-level IDs used in Vietnam:

| ID | Name | Notes |
|----|------|-------|
| `00` | Payload Format Indicator | always `01` |
| `01` | Point of Initiation Method | `11` = static (reusable), `12` = dynamic (one-time) |
| `26` | Merchant Account Information — **VNPAYQR** | sub-template, GUID `A000000775` |
| `38` | Merchant Account Information — **VietQR / Napas** | sub-template, GUID `A000000727` |
| `52` | Merchant Category Code (MCC) | 4 digits; often left empty in VN |
| `53` | Transaction Currency | ISO 4217 numeric — VND = `704` |
| `54` | Transaction Amount | amount in VND, no decimals |
| `55` | Tip or Convenience Indicator | rarely used |
| `56` | Convenience Fee — Fixed | |
| `57` | Convenience Fee — Percentage | |
| `58` | Country Code | `VN` |
| `59` | Merchant Name | UTF-8 |
| `60` | Merchant City | |
| `61` | Postal Code | |
| `62` | Additional Data Field Template | sub-template with Purpose, Bill, Store… |
| `63` | CRC | 4 uppercase hex chars |
| `64` | Merchant Information Language Template | multi-language |
| `65`–`79` | RFU for EMVCo | reserved |
| `80`–`99` | Unreserved Templates | wallets use tag 80 for private metadata |

## 2. Field 38 — VietQR (Napas 247)

The most common form, used by every Napas-member bank.

```
38 LL
   00 10 A000000727                       (Napas / VietQR GUID)
   01 LL
        00 06 <BANK_BIN>                  (Napas bank code)
        01 LL <BANK_NUMBER>               (account or card number)
   02 08 <QRIBFTTA | QRIBFTTC>            (service code)
```

Service codes:

| Code | Meaning |
|------|---------|
| `QRIBFTTA` | Inter-Bank Fund Transfer **To Account** |
| `QRIBFTTC` | Inter-Bank Fund Transfer **To Card** |

`BANK_BIN` is the 6-digit Napas-assigned identifier — see `constants/banks.go`
for the full list (Vietcombank `970436`, Techcombank `970407`, MB Bank `970422`,
ACB `970416`, BVBank `970454`, BIDV `970418`, …).

### Real-world examples

ACB static QR (account `257678859`):

```
00020101021138530010A0000007270123000697041601092576788590208QRIBFTTA53037045802VN6304AE9F
```

Dynamic QR for the same account, 10,000 VND, purpose `Chuyen tien`:

```
00020101021238530010A0000007270123000697041601092576788590208QRIBFTTA53037045405100005802VN62150811Chuyen tien630453E6
```

The only differences: `01 02 12` instead of `11`, field `54` (amount), field
`62` (additional data), and a different CRC.

## 3. Field 26 — VNPAYQR

Same template structure but with a different GUID and no bank sub-template.

```
26 LL
   00 10 A000000775                       (VNPAY GUID)
   01 LL <MERCHANT_ID>                    (merchant ID assigned by VNPAY)
```

Other fields (merchant name, store, terminal, mobile…) live in `59` and `62`.
Example:

```
00020101021126280010A0000007750110010531314453037045408210900005802VN5910CELLPHONES62600312CPSHN ONLINE0517021908061613127850705ONLHN0810CellphoneS63047685
```

Decoded:

- merchant id = `0105313144`
- amount = `21090000`
- store = `CPSHN ONLINE`
- reference = `02190806161312785`
- terminal = `ONLHN`
- purpose = `CellphoneS`

## 4. Field 62 — Additional Data

Sub-template shared between VietQR and VNPAYQR.

| Sub-ID | Name |
|--------|------|
| `01` | Bill Number |
| `02` | Mobile Number |
| `03` | Store Label |
| `04` | Loyalty Number |
| `05` | Reference Label |
| `06` | Customer Label |
| `07` | Terminal Label |
| `08` | Purpose of Transaction |
| `09` | Additional Consumer Data Request |

Length notes: most banks accept up to **25 ASCII characters** for purpose;
keep Vietnamese text un-accented because some POS apps fail to render Unicode
in this field.

## 5. CRC-16/CCITT-FALSE — field 63

Napas requires:

- Polynomial `0x1021`
- Init value `0xFFFF`
- No input/output reflection, no xorout
- Input = the entire string from field 00 up to and including the `6304`
  prefix of field 63
- Output formatted as `%04X` (uppercase hex)

Reference vector: `CRC16-CCITT-FALSE("123456789") = 0x29B1`.

`emvqr/serve` implements this exactly — see `test/crc_test.go` for fixtures
extracted from five real QR strings.

## 6. E-wallets — MoMo and ZaloPay (via BVBank)

MoMo and ZaloPay don't have their own Napas AID; each wallet maps to a virtual
account at **BVBank (BanViet — BIN `970454`)**. When a payer scans the QR,
funds travel through Napas 247 into that BVBank account, and BVBank forwards
them to the wallet.

### MoMo

```
00020101021138620010A00000072701320006970454011899MM24011M348750800208QRIBFTTA53037045802VN62190515MOMOW2W3487508080030466304EBC8
```

Notes:

- BankBin = `970454` (BVBank), BankNumber = `99MM24011M34875080`
- Field 62: reference label = `MOMOW2W` + the last 8 chars of the bank number
- Unreserved field `80`: the last 3 digits of the receiver's **phone number**
  (`046` in this fixture)

### ZaloPay

```
00020101021138620010A00000072701320006970454011899ZP24009M072482670208QRIBFTTA53037045802VN6304073C
```

- BankBin = `970454`, BankNumber = `99ZP24009M07248267`
- Account numbers start with `99ZP…` (ZaloPay prefix)
- Additional metadata is sometimes embedded under field 26, but its semantics
  are not publicly documented yet.

## 7. Popular banks

Excerpted from `constants/banks.go` — 60+ banks in total.

| Bank | BIN | Short | SWIFT | VietQR transfer |
|------|-----|-------|-------|-----------------|
| Vietcombank | 970436 | VCB | BFTVVNVX | ✅ |
| VietinBank | 970415 | CTG | ICBVVNVX | ✅ |
| BIDV | 970418 | BIDV | BIDVVNVX | ✅ |
| Agribank | 970405 | AGB | VBAAVNVX | ✅ |
| Techcombank | 970407 | TCB | VTCBVNVX | ✅ |
| MB Bank | 970422 | MB | MSCBVNVX | ✅ |
| ACB | 970416 | ACB | ASCBVNVX | ✅ |
| VPBank | 970432 | VPB | VPBKVNVX | ✅ |
| TPBank | 970423 | TPB | TPBVVNVX | ✅ |
| Sacombank | 970403 | STB | SGTTVNVX | ✅ |
| HDBank | 970437 | HDB | HDBCVNVX | ✅ |
| OCB | 970448 | OCB | ORCOVNVX | ✅ |
| SHB | 970443 | SHB | SHBAVNVX | ✅ |
| Eximbank | 970431 | EIB | EBVIVNVX | ✅ |
| MSB | 970426 | MSB | MCOBVNVX | ✅ |
| SeABank | 970440 | SEA | SEAVVNVX | ✅ |
| VIB | 970441 | VIB | VNIBVNVX | ✅ |
| BVBank (BanViet) | 970454 | BVB | — | ✅ (MoMo, ZaloPay) |
| LPBank | 970449 | LPB | LVBKVNVX | ✅ |
| Nam A Bank | 970428 | NAB | NAMAVNVX | ✅ |

Foreign-flagged banks operating in Vietnam (Shinhan `970424`, Woori `970457`,
UOB `970458`, Standard Chartered `970410`, HSBC `458761`, Indovina `970434`,
Public Bank `970439`, KasikornBank `668888`, …) are also in `BanksMap`.

## 8. Library coverage

After the bug fixes, the following flows match the upstream
`xuannghia/vietnam-qr-pay` fixtures byte-for-byte:

| Flow | Test | Status |
|------|------|--------|
| Build VietQR static | `TestBuild_VietQR_Static_ACB` | ✅ |
| Build VietQR dynamic | `TestBuild_VietQR_Dynamic_ACB` | ✅ |
| Build MoMo (QRIBFTTA + reference + tag 80) | `TestBuild_MoMo_via_BVBank` | ✅ |
| Build ZaloPay | `TestBuild_ZaloPay_via_BVBank` | ✅ |
| Build transfer-to-card (QRIBFTTC) | `TestBuild_VietQR_TransferToCard` | ✅ |
| Build shorthand `instance.Build` | `TestBuild_Instance_Shorthand` | ✅ |
| Parse static / dynamic VietQR | `TestParse_VietQR_*` | ✅ |
| Parse MoMo (unreserved, reference) | `TestParse_MoMo_*` | ✅ |
| Round-trip parse → build → parse | `TestRoundTrip_*` | ✅ |
| Mutate amount/purpose then rebuild | `TestMutate_ChangeAmountAndPurpose` | ✅ |
| CRC-16/CCITT-FALSE standard vector | `TestCRC16CCITT_KnownVectors` | ✅ |
| PNG / Base64 / Data URL image | `TestPNG_*` | ✅ |

## 9. Known gaps

The items below are **not yet** implemented — reasonable TODOs:

1. **VNPAYQR end-to-end fixture.** `InitVNPayQR` exists but no real-merchant
   reference string is asserted yet.
2. **Field 64 — Merchant Information Language Template** (`64`): currently
   read/written in `instance` but not surfaced on `qrPay`. Some banks
   (Vietcombank app) emit the Vietnamese merchant name through this field.
3. **Fields 55–57 (Tip & Fee).** Struct slots exist but no tests; EMVCo
   requires Convenience Fee Fixed when Tip Indicator = 02, not enforced.
4. **Per-field length / charset validation** per EMVCo (amount 1–13, country 2
   chars, currency 3 chars, CRC exactly 4 hex). The library only checks CRC.
5. **Unify `instance` and `qrPay`.** Two parallel APIs cause confusion; merge
   them, or clarify each role in the README.
6. **Vietnamese-diacritic encoding helper** for field 59/62: bytes accept
   UTF-8 but many POS scanners corrupt it; an optional helper to strip
   diacritics would be useful.

## 10. Sources

- EMVCo, *EMV® QR Code Specification for Payment Systems — Merchant-Presented
  Mode v1.1*, https://www.emvco.com/specifications/
- Napas, *VietQR Format Regulations for the NAPAS247 Service* (Vietnamese
  original "Quy dinh Dinh Dang QR VietQR trong Dich vu NAPAS247", available
  in both Vietnamese and English on Studocu / Scribd)
- `xuannghia/vietnam-qr-pay` (TypeScript reference; source of the MoMo/
  ZaloPay/ACB fixtures): https://github.com/xuannghia/vietnam-qr-pay
- VNPAY Sandbox Portal: https://sandbox.vnpayment.vn/apis/
- NCB VNPAY-QR API Portal: https://develop.ncb-bank.vn/apis/vnpay-qr/overview
- `thanhtinhpas1/vietqr-parser` (Go reference):
  https://github.com/thanhtinhpas1/vietqr-parser
- `subiz/vietqr` (Go reference): https://pkg.go.dev/github.com/subiz/vietqr
