package psip

import (
	"github.com/Comcast/gots/packet"
)

const (
	// BasePID is the PID for all ATSC PSIP tables
	BasePID     = uint16(0x1FFB)
	MGTTableID  = uint8(0xC7)
	TVCTTableID = uint8(0xC8)
	CVCTTableID = uint8(0xC9)
	SVCTTableID = uint8(0xDA)
)

type Channel interface {
	MajorNumber() uint16
	MinorNumber() uint16
	DescriptorsLength() int
	ShortName() string
	Length() int
}

// Master Guide Table
type MGT interface {
}

// TVCT represents operations on a Terrestrial Virtual Channel Table
type VCT interface {
	SectionSyntaxIndicator() bool
	PrivateIndicator() bool
	SectionLength() int
	TransportStreamID() uint16
	VersionNumber() uint8
	CurrentNextIndicator() bool
	SectionNumber() uint8
	LastSectionNumber() uint8
	ProtocolVersion() uint8
	NumChannelsInSection() int
	Channels() []Channel
	channelLength() int
	AdditionalDescriptorsLength() uint16
	Crc() []byte
}

// Tables are the collection of all ATSC PSIP tables
type Tables interface {
	Update(packet.Packet) error
	VCT() VCT
}
