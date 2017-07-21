package ts

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestUtf16(t *testing.T) {
	tests := []struct {
		input  []byte
		output string
	}{
		{
			input:  []byte{0x00, 0x66, 0x00, 0x6f, 0x00, 0x6f, 0x00, 0x00},
			output: "foo",
		},
		{
			input:  []byte{0x00, 0x62, 0x00, 0x61, 0x00, 0x72, 0x00, 0x00},
			output: "bar",
		},
	}

	for i, test := range tests {
		str := Utf16ToString(test.input, 0, len(test.input))
		if str != test.output {
			t.Errorf("Test %d decoding expected \"%s\" but got \"%s\"", i, test.output, str)
		}

		output := StringToUtf16(test.output)
		if !bytes.Equal(test.input, output) {
			t.Errorf("Test %d encoding expected %v but got %v", i, hex.Dump(test.input), hex.Dump(output))
		}
	}
}

func TestBool(t *testing.T) {
	tests := []struct {
		input    byte
		bit      int
		expected bool
	}{
		{0x00, 0, false},
		{0x01, 7, true},
		{0x00, 6, false},
		{0x02, 6, true},
		{0x00, 5, false},
		{0x04, 5, true},
		{0x00, 4, false},
		{0x08, 4, true},
		{0x00, 3, false},
		{0x10, 3, true},
		{0x00, 2, false},
		{0x20, 2, true},
		{0x00, 1, false},
		{0x40, 1, true},
		{0x00, 0, false},
		{0x80, 0, true},
	}

	for i, test := range tests {
		received := Bool(test.input, test.bit)
		if received != test.expected {
			t.Errorf("Test %d expected %v but got %v", i, test.expected, received)
		}
	}
}

func TestSetBool(t *testing.T) {
	tests := []struct {
		input    bool
		original byte
		expected byte
		bit      int
	}{
		{true, 0x00, 0x01, 7},
		{false, 0x01, 0x00, 7},
		{true, 0x00, 0x02, 6},
		{false, 0x02, 0x00, 6},
		{true, 0xfe, 0xff, 7},
		{false, 0xff, 0xfe, 7},
		{true, 0xfd, 0xff, 6},
		{false, 0xff, 0xfd, 6},
	}

	for i, test := range tests {
		SetBool(&test.original, test.bit, test.input)
		if test.expected != test.original {
			t.Errorf("Test %d failed. Expected 0x%02x but got 0x%02x", i, test.expected, test.original)
		}
	}
}

func TestUimsbf64(t *testing.T) {
	tests := []struct {
		input []byte
		width int
		exp64 uint64
		exp32 uint32
		exp16 uint16
	}{
		{[]byte{0x01}, 64, uint64(0x01), uint32(0x01), uint16(0x01)},
		{[]byte{0x01, 0x00}, 64, uint64(0x0100), uint32(0x0100), uint16(0x0100)},
		{[]byte{0x01, 0x00, 0x00}, 64, uint64(0x010000), uint32(0x010000), uint16(0x0000)},
		{[]byte{0x01, 0x00, 0x00, 0x00}, 64, uint64(0x01000000), uint32(0x01000000), uint16(0x0000)},
		{[]byte{0x01, 0x00, 0x00, 0x00, 0x00}, 64, uint64(0x0100000000), uint32(0x00000000), uint16(0x0000)},
		{[]byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00}, 64, uint64(0x010000000000), uint32(0x00000000), uint16(0x0000)},
		{[]byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 64, uint64(0x01000000000000), uint32(0x00000000), uint16(0x0000)},
		{[]byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 64, uint64(0x0100000000000000), uint32(0x00000000), uint16(0x0000)},
		{[]byte{0x01}, 8, uint64(0x01), uint32(0x01), uint16(0x01)},
		{[]byte{0xff, 0x01}, 8, uint64(0x01), uint32(0x01), uint16(0x01)},
		{[]byte{0xff, 0xff, 0x01}, 16, uint64(0xff01), uint32(0xff01), uint16(0xff01)},
		{[]byte{0xff, 0xff, 0x01}, 22, uint64(0x3fff01), uint32(0x3fff01), uint16(0xff01)},
		{[]byte{0xff, 0xff}, 13, uint64(0x1fff), uint32(0x1fff), uint16(0x1fff)},
	}

	for i, test := range tests {
		r64 := Uimsbf64(test.input, test.width)
		if test.exp64 != r64 {
			t.Errorf("Test %d expected %v but got %v", i, test.exp64, r64)
		}

		r32 := Uimsbf32(test.input, test.width)
		if test.exp32 != r32 {
			t.Errorf("Test %d expected %v but got %v", i, test.exp32, r32)
		}

		r16 := Uimsbf16(test.input, test.width)
		if test.exp16 != r16 {
			t.Errorf("Test %d expected %v but got %v", i, test.exp16, r16)
		}
	}
}

func TestUimsbf8(t *testing.T) {
	tests := []struct {
		input    byte
		expected uint8
	}{
		{0x01, uint8(1)},
		{0xff, uint8(0xff)},
	}

	for i, test := range tests {
		received := Uimsbf8(test.input)
		if test.expected != received {
			t.Errorf("Test %d expected %v but got %v", i, test.expected, received)
		}
	}
}

func TestComputeCRC(t *testing.T) {
	tests := []struct {
		input    []byte
		expected []byte
	}{
		{[]byte{0x00, 0xb0, 0x0d, 0x59, 0x81, 0xeb, 0x00, 0x00, 0x00, 0x01, 0xe0, 0x42}, []byte{0x5e, 0x44, 0x05, 0x9a}},
		{[]byte{0xff, 0xfe, 0xfd, 0xfc}, []byte{0xac, 0x69, 0x14, 0x51}},
	}

	for i, test := range tests {
		received := ComputeCRC(test.input)
		if !bytes.Equal(test.expected, received) {
			t.Errorf("Test %d expected\n%s\nbut got\n%s", i, hex.Dump(test.expected), hex.Dump(received))
		}
	}
}
