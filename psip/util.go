package psip

import (
	"unicode/utf16"
)

var bitMask = []byte{
	0x80,
	0x40,
	0x20,
	0x10,
	0x08,
	0x04,
	0x02,
	0x01,
}

var mask = []byte{
	0x00,
	0x01,
	0x03,
	0x07,
	0x0f,
	0x1f,
	0x3f,
	0x7f,
	0xff,
}

func shortBuffer(b []byte, offset int, length int) bool {
	return false
}

func utf16ToString(b []byte, offset int, length int) string {
	if shortBuffer(b, offset, length) {
		return ""
	}

	utfStr := make([]uint16, 0)
	for i := 0; i < length; i += 2 {
		if b[i] == 0 && b[i+1] == 0 {
			break
		}
		utfStr = append(utfStr, uint16(b[i]<<8)|uint16(b[i+1]))
	}
	return string(utf16.Decode(utfStr))
}

func getBool(b []byte, offset int, bit int) bool {
	if shortBuffer(b, offset, 1) {
		return false
	}

	return b[offset]&bitMask[bit] > 0
}

func getUimsbf16(b []byte, offset int, width int) uint16 {
	if shortBuffer(b, offset, width) {
		return 0
	}

	return uint16(b[offset]&mask[width-8])<<8 | uint16(b[offset+1])
}

func getUimsbf8(b []byte, offset int) uint8 {
	if shortBuffer(b, offset, 1) {
		return 0
	}

	return uint8(b[offset])
}
