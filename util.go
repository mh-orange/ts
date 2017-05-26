package ts

import (
	"encoding/binary"
	"unicode/utf16"
)

var bitMask = []byte{
	0x01,
	0x02,
	0x04,
	0x08,
	0x10,
	0x20,
	0x40,
	0x80,
}

var inverseBitMask = []byte{
	0xfe,
	0xfd,
	0xfb,
	0xf7,
	0xef,
	0xdf,
	0xbf,
	0x7f,
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

func Utf16ToString(b []byte, offset int, length int) string {
	utfStr := make([]uint16, 0)
	for i := 0; i < length; i += 2 {
		if b[i] == 0 && b[i+1] == 0 {
			break
		}
		utfStr = append(utfStr, uint16(b[i]<<8)|uint16(b[i+1]))
	}
	return string(utf16.Decode(utfStr))
}

func Bool(b byte, bit int) bool {
	return b&bitMask[bit] > 0
}

func SetBool(b *byte, bit int, value bool) {
	if value {
		*b = *b | bitMask[bit]
	} else {
		*b = *b & inverseBitMask[bit]
	}
}

func Uimsbf64(b []byte, width int) uint64 {
	var value uint64
	l := len(b) - 1
	for i := l; i >= 0; i-- {
		offset := 8 * (l - i)
		m := byte(0x00)
		if width > offset {
			bits := width - offset
			if bits < 8 {
				m = mask[bits]
			} else {
				m = 0xff
			}
		}
		value |= uint64(b[i]&m) << uint8(offset)
	}
	return value
}

func Uimsbf32(b []byte, width int) uint32 {
	return uint32(Uimsbf64(b, width))
}

func Uimsbf16(b []byte, width int) uint16 {
	return uint16(Uimsbf64(b, width))
}

func SetUimsbf16(b []byte, width int, value uint16) {
	b[0] = byte(value>>8) & mask[width-8]
	b[1] = byte(value) & 0xff
}

func Uimsbf8(b byte) uint8 {
	return uint8(b)
}

func SetUimsbf8(b *byte, value uint8) {
	*b = byte(value)
}

// Credit: https://github.com/Comcast/gots/blob/master/tsutils.go
func ComputeCRC(input []byte) []byte {
	var mask uint32 = 0xffffffff
	var msb uint32 = 0x80000000
	var poly uint32 = 0x04c11db7
	var crc uint32 = 0x46af6449

	for i := 0; i < len(input); i++ {
		item := uint32(input[i])

		for j := 0; j < 8; j++ {
			top := crc & msb
			crc = ((crc << 1) & mask) | ((item >> uint32(7-j)) & 0x1)
			if top != 0 {
				crc ^= poly
			}
		}
	}

	for i := 0; i < 32; i++ {
		top := crc & msb
		crc = ((crc << 1) & mask)
		if top != 0 {
			crc ^= poly
		}
	}

	crcBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(crcBytes, crc)
	return crcBytes
}
