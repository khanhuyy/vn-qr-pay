package bankqr_test

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"strings"
	"testing"

	"github.com/khanhuyy/emvqr/bankqr"
)

// PNG magic bytes.
var pngMagic = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

func vcbQR() *bankqr.QR {
	return bankqr.New(bankqr.Options{
		BIN:           "970436",
		AccountNumber: "1234567890",
		Amount:        50_000,
		Purpose:       "Cam on",
	})
}

func TestPNG_ReturnsValidPNGBytes(t *testing.T) {
	raw, err := vcbQR().PNG(256)
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) < 8 || !bytes.Equal(raw[:8], pngMagic) {
		t.Fatalf("output is not a PNG (first 8 bytes: %x)", raw[:8])
	}

	img, err := png.Decode(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("png.Decode: %v", err)
	}
	b := img.Bounds()
	if b.Dx() != 256 || b.Dy() != 256 {
		t.Errorf("size: want 256x256, got %dx%d", b.Dx(), b.Dy())
	}
}

func TestPNG_DifferentSizesYieldDifferentBytes(t *testing.T) {
	small, _ := vcbQR().PNG(128)
	large, _ := vcbQR().PNG(512)
	if bytes.Equal(small, large) {
		t.Error("128px and 512px outputs should differ")
	}
}

func TestPNGBase64_DecodesToPNG(t *testing.T) {
	b64, err := vcbQR().PNGBase64(256)
	if err != nil {
		t.Fatal(err)
	}
	if b64 == "" {
		t.Fatal("empty base64 output")
	}

	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		t.Fatalf("base64 decode: %v", err)
	}
	if !bytes.HasPrefix(raw, pngMagic) {
		t.Errorf("decoded bytes are not a PNG: %x", raw[:8])
	}
}

func TestPNGDataURL_HasCorrectPrefixAndPayload(t *testing.T) {
	url, err := vcbQR().PNGDataURL(256)
	if err != nil {
		t.Fatal(err)
	}
	const prefix = "data:image/png;base64,"
	if !strings.HasPrefix(url, prefix) {
		t.Fatalf("missing prefix %q in: %s", prefix, url[:min(60, len(url))])
	}
	payload := url[len(prefix):]
	raw, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		t.Fatalf("payload not valid base64: %v", err)
	}
	if !bytes.HasPrefix(raw, pngMagic) {
		t.Error("payload does not decode to PNG")
	}
}

func TestPNGWithLevel_HigherCorrectionStillValid(t *testing.T) {
	for _, lvl := range []bankqr.ErrorCorrection{
		bankqr.ErrorCorrectionLow,
		bankqr.ErrorCorrectionMedium,
		bankqr.ErrorCorrectionHigh,
		bankqr.ErrorCorrectionHighest,
	} {
		raw, err := vcbQR().PNGWithLevel(256, lvl)
		if err != nil {
			t.Errorf("level %v: %v", lvl, err)
			continue
		}
		if _, err := png.Decode(bytes.NewReader(raw)); err != nil {
			t.Errorf("level %v produced invalid PNG: %v", lvl, err)
		}
	}
}

// Indirect round-trip via determinism.
func TestPNG_DeterministicForSameInput(t *testing.T) {
	a, _ := vcbQR().PNG(256)
	b, _ := vcbQR().PNG(256)
	if !bytes.Equal(a, b) {
		t.Error("same QR + same size must produce byte-identical PNG")
	}
}

// min for Go < 1.21.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
