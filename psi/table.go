package psi

import (
	"github.com/mh-orange/ts"
)

type Table []byte

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

func (t Table) ProtocolVersion() uint8 {
	return ts.Uimsbf8(t[8])
}

func (t Table) Data() []byte {
	// 6 bytes precede the data field and 4 bytes (crc)
	// follow the data field (10 bytes) the offset is
	// starting at byte 9, so the ending index is the
	// section length
	return t[9:t.SectionLength()]
}

func (t Table) CRC() []byte {
	return t[len(t)-4:]
}
