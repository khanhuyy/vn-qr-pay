// qrpay CLI: build, parse, crc, banks.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/khanhuyy/emvgo/constants"
	"github.com/khanhuyy/emvgo/serve"
)

const usage = `qrpay — Vietnam QR Pay (VietQR / VNPAYQR) command-line tool

Usage:
  qrpay build  --bin BIN --account NUMBER [--amount N] [--purpose TEXT] [--service QRIBFTTA|QRIBFTTC] [--json]
  qrpay parse  QR_CONTENT [--json]
  qrpay crc    CONTENT
  qrpay banks  [--filter KEYWORD] [--json]
  qrpay help

Run 'qrpay <command> --help' for command-specific flags.`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	switch os.Args[1] {
	case "build":
		runBuild(os.Args[2:])
	case "parse":
		runParse(os.Args[2:])
	case "crc":
		runCRC(os.Args[2:])
	case "banks":
		runBanks(os.Args[2:])
	case "help", "-h", "--help":
		fmt.Println(usage)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n%s\n", os.Args[1], usage)
		os.Exit(2)
	}
}

// build

func runBuild(args []string) {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	bin := fs.String("bin", "", "bank BIN, e.g. 970436 (Vietcombank). Required for VietQR.")
	account := fs.String("account", "", "bank account or card number. Required for VietQR.")
	amount := fs.String("amount", "", "transaction amount in VND (omit for static QR).")
	purpose := fs.String("purpose", "", "transfer purpose / remark (max 25 ASCII chars recommended).")
	service := fs.String("service", "QRIBFTTA", "service code: QRIBFTTA (account) or QRIBFTTC (card).")
	asJSON := fs.Bool("json", false, "emit structured JSON instead of just the QR string.")
	_ = fs.Parse(args)

	if *bin == "" || *account == "" {
		fmt.Fprintln(os.Stderr, "build: --bin and --account are required")
		fs.Usage()
		os.Exit(2)
	}

	qr := serve.InitVietQR(serve.InitVietQROptions{
		BankBin:    *bin,
		BankNumber: *account,
		Amount:     *amount,
		Purpose:    *purpose,
		Service:    *service,
	})
	content := qr.Build()

	if *asJSON {
		out := buildSummary{
			Content:    content,
			InitMethod: qr.InitMethod,
			Provider:   string(qr.Provider.Name),
			GUID:       qr.Provider.GUID,
			Service:    qr.Provider.Service,
			BankBin:    qr.Consumer.BankBin,
			BankNumber: qr.Consumer.BankNumber,
			Amount:     qr.Amount,
			Purpose:    qr.AdditionalData.Purpose,
		}
		_ = json.NewEncoder(os.Stdout).Encode(out)
		return
	}

	fmt.Println(content)
}

type buildSummary struct {
	Content    string `json:"content"`
	InitMethod string `json:"initMethod"`
	Provider   string `json:"provider"`
	GUID       string `json:"guid"`
	Service    string `json:"service,omitempty"`
	BankBin    string `json:"bankBin,omitempty"`
	BankNumber string `json:"bankNumber,omitempty"`
	Amount     string `json:"amount,omitempty"`
	Purpose    string `json:"purpose,omitempty"`
}

// parse

func runParse(args []string) {
	fs := flag.NewFlagSet("parse", flag.ExitOnError)
	asJSON := fs.Bool("json", false, "emit structured JSON instead of a human-readable table.")
	_ = fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "parse: missing QR content")
		os.Exit(2)
	}

	qr := serve.NewQRPay(fs.Arg(0))
	if !qr.IsValid {
		fmt.Fprintln(os.Stderr, "parse: QR is invalid (CRC mismatch or malformed TLV)")
		os.Exit(1)
	}

	if *asJSON {
		_ = json.NewEncoder(os.Stdout).Encode(qr)
		return
	}

	w := os.Stdout
	fmt.Fprintln(w, "Valid           :", qr.IsValid)
	fmt.Fprintln(w, "Version         :", qr.Version)
	fmt.Fprintln(w, "Init method     :", qr.InitMethod, initMethodLabel(qr.InitMethod))
	fmt.Fprintln(w, "Provider name   :", qr.Provider.Name)
	fmt.Fprintln(w, "Provider GUID   :", qr.Provider.GUID)
	if qr.Provider.Service != "" {
		fmt.Fprintln(w, "Service code    :", qr.Provider.Service, serviceLabel(qr.Provider.Service))
	}
	if qr.Consumer.BankBin != "" {
		fmt.Fprintln(w, "Bank BIN        :", qr.Consumer.BankBin, "("+bankShortName(qr.Consumer.BankBin)+")")
		fmt.Fprintln(w, "Bank account    :", qr.Consumer.BankNumber)
	}
	if qr.Merchant.Id != "" {
		fmt.Fprintln(w, "Merchant ID     :", qr.Merchant.Id)
	}
	if qr.Merchant.Name != "" {
		fmt.Fprintln(w, "Merchant name   :", qr.Merchant.Name)
	}
	if qr.Amount != "" {
		fmt.Fprintln(w, "Amount          :", qr.Amount, "VND")
	}
	fmt.Fprintln(w, "Currency        :", qr.Currency)
	fmt.Fprintln(w, "Nation          :", qr.Nation)
	if qr.AdditionalData.Purpose != "" {
		fmt.Fprintln(w, "Purpose         :", qr.AdditionalData.Purpose)
	}
	if qr.AdditionalData.Reference != "" {
		fmt.Fprintln(w, "Reference       :", qr.AdditionalData.Reference)
	}
	if qr.AdditionalData.Store != "" {
		fmt.Fprintln(w, "Store           :", qr.AdditionalData.Store)
	}
	if qr.AdditionalData.Terminal != "" {
		fmt.Fprintln(w, "Terminal        :", qr.AdditionalData.Terminal)
	}
	if len(qr.Unreserved) > 0 {
		fmt.Fprintln(w, "Unreserved (80-99):")
		for k, v := range qr.Unreserved {
			fmt.Fprintf(w, "  %s = %s\n", k, v)
		}
	}
	fmt.Fprintln(w, "CRC             :", qr.CRC)
}

func initMethodLabel(m string) string {
	switch m {
	case "11":
		return "(static / reusable)"
	case "12":
		return "(dynamic / one-time)"
	}
	return ""
}

func serviceLabel(s string) string {
	switch s {
	case "QRIBFTTA":
		return "(transfer to account)"
	case "QRIBFTTC":
		return "(transfer to card)"
	}
	return ""
}

// crc

func runCRC(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "crc: missing content")
		os.Exit(2)
	}
	fmt.Printf("%04X\n", serve.CRC16CCITT(args[0]))
}

// banks

type bankRow struct {
	Key          string `json:"key"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	ShortName    string `json:"shortName"`
	BIN          string `json:"bin"`
	VietQRStatus int    `json:"vietQRStatus"`
}

func runBanks(args []string) {
	fs := flag.NewFlagSet("banks", flag.ExitOnError)
	filter := fs.String("filter", "", "case-insensitive substring filter on key / code / name / BIN.")
	asJSON := fs.Bool("json", false, "emit JSON array instead of a table.")
	_ = fs.Parse(args)

	needle := strings.ToLower(*filter)

	var rows []bankRow
	for _, b := range constants.BanksMap {
		if needle != "" {
			hay := strings.ToLower(string(b.Key) + " " + string(b.Code) + " " + b.Name + " " + b.ShortName + " " + b.BIN)
			if !strings.Contains(hay, needle) {
				continue
			}
		}
		rows = append(rows, bankRow{
			Key:          string(b.Key),
			Code:         string(b.Code),
			Name:         b.Name,
			ShortName:    b.ShortName,
			BIN:          b.BIN,
			VietQRStatus: int(b.VietQRStatus),
		})
	}

	if *asJSON {
		_ = json.NewEncoder(os.Stdout).Encode(rows)
		return
	}

	fmt.Printf("%-22s %-6s %-22s %s\n", "KEY", "BIN", "SHORT NAME", "NAME")
	for _, r := range rows {
		fmt.Printf("%-22s %-6s %-22s %s\n", r.Key, r.BIN, r.ShortName, r.Name)
	}
}

func bankShortName(bin string) string {
	for _, b := range constants.BanksMap {
		if b.BIN == bin {
			return b.ShortName
		}
	}
	return "unknown"
}
