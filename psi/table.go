package psi

import (
	"github.com/mh-orange/ts"
)

type Table []byte

func CreateTable(id uint8, data []byte) Table {
	// table is data length + 8 (for header bits) + 4 (for CRC)
	length := len(data) + 12

	table := make(Table, length)
	table[0] = id
	table[1] |= 0x30
	table[1] = uint8(uint16(length-3) >> 8 & 0x03)
	table[2] = uint8(0xff & (length - 3))

	copy(table[8:], data)
	copy(table[length-4:], ts.ComputeCRC(table[0:length-4]))

	return table
}

func (t Table) ID() uint8 {
	return ts.Uimsbf8(t[0])
}

func (t Table) SectionSyntaxIndicator() bool {
	return ts.Bool(t[1], 0)
}

func (t Table) PrivateIndicator() bool {
	return ts.Bool(t[1], 1)
}

func (t Table) SectionLength() int {
	return int(ts.Uimsbf16(t[1:3], 12))
}

func (t Table) IDExtension() uint16 {
	return ts.Uimsbf16(t[3:5], 16)
}

func (t Table) VersionNumber() uint8 {
	return uint8((0x3f & t[5]) >> 1)
}

func (t Table) IsCurrent() bool {
	return ts.Bool(t[5], 7)
}

func (t Table) IsNext() bool {
	return !t.IsCurrent()
}

func (t Table) SectionNumber() uint8 {
	return ts.Uimsbf8(t[6])
}

func (t Table) LastSectionNumber() uint8 {
	return ts.Uimsbf8(t[7])
}

func (t Table) Data() []byte {
	return t[8 : len(t)-4]
}

func (t Table) CRC() []byte {
	return t[len(t)-4:]
}
