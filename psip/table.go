package psip

import (
	"github.com/mh-orange/ts"
)

type Table interface {
	ID() uint8
	SectionSyntaxIndicator() bool
	PrivateIndicator() bool
	SectionLength() int
	IDExtension() uint16
	VersionNumber() uint8
	IsCurrent() bool
	IsNext() bool
	SectionNumber() uint8
	LastSectionNumber() uint8
	ProtocolVersion() uint8
	CRC() []byte
}

type table []byte

func (t table) ID() uint8 {
	return ts.Uimsbf8(t[0])
}

func (t table) SectionSyntaxIndicator() bool {
	return ts.Bool(t[1], 0)
}

func (t table) PrivateIndicator() bool {
	return ts.Bool(t[1], 1)
}

func (t table) SectionLength() int {
	return int(ts.Uimsbf16(t[1:3], 12))
}

func (t table) IDExtension() uint16 {
	return ts.Uimsbf16(t[3:5], 16)
}

func (t table) VersionNumber() uint8 {
	return uint8((0x3f & t[5]) >> 1)
}

func (t table) IsCurrent() bool {
	return ts.Bool(t[5], 7)
}

func (t table) IsNext() bool {
	return !t.IsCurrent()
}

func (t table) SectionNumber() uint8 {
	return ts.Uimsbf8(t[6])
}

func (t table) LastSectionNumber() uint8 {
	return ts.Uimsbf8(t[7])
}

func (t table) ProtocolVersion() uint8 {
	return ts.Uimsbf8(t[8])
}

func (t table) Data() []byte {
	// 6 bytes precede the data field and 4 bytes (crc)
	// follow the data field (10 bytes) the offset is
	// starting at byte 9, so the ending index is the
	// section length
	return t[9:t.SectionLength()]
}

func (t table) CRC() []byte {
	return t[len(t)-4:]
}
