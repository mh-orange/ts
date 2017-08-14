package psip

import (
	"github.com/mh-orange/ts"
	"github.com/mh-orange/ts/psi"
)

const (
	MODULATION_8VSB = 0x04
)

type Channel interface {
	ShortName() string
	MajorNumber() uint16
	MinorNumber() uint16
	ModulationMode() uint8
	CarrierFrequency() uint32
	TSID() uint16
	Program() uint16
	ETMLocation() uint8
	AccessControlled() bool
	Hidden() bool
	HideGuide() bool
	ServiceType() uint8
	SourceID() uint16
	DescriptorsLength() int
	Descriptors() []psi.Descriptor
	Length() int
	CRC() []byte
}

type channel []byte

const (
	c_minimum_length = 32
	c_major_offset   = 14
)

func (c channel) ShortName() string {
	return ts.Utf16ToString(c, 0, 14)
}

func (c channel) MajorNumber() uint16 {
	return ts.Uimsbf16(c[14:16], 16) >> 2 & 0x3ff
}

func (c channel) MinorNumber() uint16 {
	return ts.Uimsbf16(c[15:17], 10)
}

func (c channel) ModulationMode() uint8 {
	return ts.Uimsbf8(c[17])
}

func (c channel) CarrierFrequency() uint32 {
	return ts.Uimsbf32(c[18:22], 32)
}

func (c channel) TSID() uint16 {
	return ts.Uimsbf16(c[22:24], 16)
}

func (c channel) Program() uint16 {
	return ts.Uimsbf16(c[24:26], 16)
}

func (c channel) ETMLocation() uint8 {
	return ts.Uimsbf8(c[26]) >> 6
}

func (c channel) AccessControlled() bool {
	return ts.Bool(c[26], 5)
}

func (c channel) Hidden() bool {
	return ts.Bool(c[26], 4)
}

func (c channel) HideGuide() bool {
	return ts.Bool(c[26], 1)
}

func (c channel) ServiceType() uint8 {
	return ts.Uimsbf8(c[27])
}

func (c channel) SourceID() uint16 {
	return ts.Uimsbf16(c[28:30], 16)
}

func (c channel) DescriptorsLength() int {
	return int(c[30]&0x03)<<8 | int(c[31])
}

func (c channel) Descriptors() []psi.Descriptor {
	return psi.Descriptors(c[32 : 32+c.DescriptorsLength()])
}

func (c channel) CRC() []byte {
	return c[len(c)-4 : len(c)]
}

func (c channel) Length() int {
	return 36 + c.DescriptorsLength()
}
