package utils

import "unicode/utf8"

func ToPointer[T any](t T) *T {
	return &t
}

// Utility functions
func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

func StringToUint8Array(content string) []byte {
	var octets []byte
	for i := 0; i < len(content); {
		r, size := utf8.DecodeRuneInString(content[i:])
		codePoint := int(r)
		switch {
		case codePoint <= 0x7F:
			octets = append(octets, byte(codePoint))
		case codePoint <= 0x7FF:
			octets = append(octets,
				0xC0|byte(codePoint>>6),
				0x80|byte(codePoint&0x3F))
		case codePoint <= 0xFFFF:
			octets = append(octets,
				0xE0|byte(codePoint>>12),
				0x80|byte((codePoint>>6)&0x3F),
				0x80|byte(codePoint&0x3F))
		case codePoint <= 0x1FFFFF:
			octets = append(octets,
				0xF0|byte(codePoint>>18),
				0x80|byte((codePoint>>12)&0x3F),
				0x80|byte((codePoint>>6)&0x3F),
				0x80|byte(codePoint&0x3F))
		}
		i += size
	}
	return octets
}
