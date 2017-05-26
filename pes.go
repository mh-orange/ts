package ts

import (
	"fmt"
)

var (
	ErrNoPESPTS                          = fmt.Errorf("PES Header has no Presentation Time Stamp")
	ErrNoPESDTS                          = fmt.Errorf("PES Header has no Decoding Time Stamp")
	ErrNoPESESCR                         = fmt.Errorf("PES Header has no Elementary Stream Clock Reference")
	ErrNoPESESRate                       = fmt.Errorf("PES Header has no Elementary Stream Rate")
	ErrNoPESAdditionalCopyInfo           = fmt.Errorf("PES Header has no additional copy information")
	ErrNoPESCRC                          = fmt.Errorf("PES Header has no CRC")
	ErrNoPESExtension                    = fmt.Errorf("PES Header has no header extension")
	ErrNoPESPrivateData                  = fmt.Errorf("PES Header has no private data")
	ErrNoPESPackHeaderField              = fmt.Errorf("PES Header has no pack header field")
	ErrNoPESProgramPacketSequenceCounter = fmt.Errorf("PES Header has no Program Packet Sequence Counter")
	ErrNoPStdBufferInfo                  = fmt.Errorf("PES Header has no P-STD Buffer information")
	ErrNoPESSecondExtension              = fmt.Errorf("PES Header has no additional PES extension")
)

type PESHeader []byte

func FillPESHeader(payload []byte) PESHeader {
	payload[0] = 0x00
	payload[1] = 0x00
	payload[2] = 0x01
	return PESHeader(payload)
}

func (p PESHeader) StreamID() uint8 {
	return Uimsbf8(p[3])
}

func (p PESHeader) Length() uint16 {
	return Uimsbf16(p[4:6], 16)
}

func (p PESHeader) updateLength() {
	length := p.extensionOffset()
	if extension, err := p.Extension(); err == nil {
		length += extension.Length()
	}

	SetUimsbf16(p[4:6], 16, uint16(length))
}

func (p PESHeader) ScramblingControl() byte {
	return (p[6] >> 4) & 0x03
}

func (p PESHeader) HasPriority() bool {
	return Bool(p[6], 3)
}

func (p PESHeader) IsAligned() bool {
	return Bool(p[6], 2)
}

func (p PESHeader) HasCopyright() bool {
	return Bool(p[6], 1)
}

func (p PESHeader) IsOriginal() bool {
	return Bool(p[6], 0)
}

func (p PESHeader) HasPTS() bool {
	return Bool(p[7], 7)
}

func (p PESHeader) SetHasPTS(hasPts bool) {
	SetBool(&p[7], 7, hasPts)
	p.updateLength()
}

func (p PESHeader) HasDTS() bool {
	return Bool(p[7], 6)
}

func (p PESHeader) SetHasDTS(hasDts bool) {
	SetBool(&p[7], 6, hasDts)
}

func (p PESHeader) HasESCR() bool {
	return Bool(p[7], 5)
}

func (p PESHeader) HasESRate() bool {
	return Bool(p[7], 4)
}

func (p PESHeader) HasDSMTrickMode() bool {
	return Bool(p[7], 3)
}

func (p PESHeader) HasAdditionalCopyInfo() bool {
	return Bool(p[7], 2)
}

func (p PESHeader) HasCRC() bool {
	return Bool(p[7], 1)
}

func (p PESHeader) HasExtension() bool {
	return Bool(p[7], 0)
}

func (p PESHeader) AdditionalLength() uint8 {
	return Uimsbf8(p[8])
}

func (p PESHeader) ptsOffset() int {
	return 9
}

func (p PESHeader) PTS() (TimeStamp, error) {
	if p.HasPTS() {
		offset := p.ptsOffset()
		return CalculatePTS(p[offset : offset+4]), nil
	}
	return nil, ErrNoPESPTS
}

func (p PESHeader) SetPTS(timestamp TimeStamp) {
	offset := p.ptsOffset()
	copy(p[offset:offset+4], DumpPTS(timestamp))
}

func (p PESHeader) dtsOffset() int {
	if p.HasPTS() {
		return p.ptsOffset() + 4
	}
	return p.ptsOffset()
}

func (p PESHeader) DTS() (TimeStamp, error) {
	if p.HasDTS() {
		offset := p.dtsOffset()
		return CalculatePTS(p[offset : offset+4]), nil
	}
	return nil, ErrNoPESDTS
}

func (p PESHeader) escrOffset() int {
	if p.HasDTS() {
		return p.dtsOffset() + 4
	}
	return p.dtsOffset()
}

func (p PESHeader) ESCR() (TimeStamp, error) {
	if p.HasESCR() {
		offset := p.escrOffset()
		return CalculateESCR(p[offset : offset+4]), nil
	}
	return nil, ErrNoPESESCR
}

func (p PESHeader) esRateOffset() int {
	if p.HasESCR() {
		return p.escrOffset() + 4
	}
	return p.escrOffset()
}

func (p PESHeader) ESRate() (uint16, error) {
	if p.HasESRate() {
		offset := p.esRateOffset()
		return Uimsbf16(p[offset:offset+2], 15) >> 1, nil
	}
	return 0, ErrNoPESESRate
}

func (p PESHeader) additionalCopyInfoOffset() int {
	if p.HasESRate() {
		return p.esRateOffset() + 2
	}
	return p.esRateOffset()
}

func (p PESHeader) AdditionalCopyInfo() (uint8, error) {
	if p.HasAdditionalCopyInfo() {
		return Uimsbf8(p[p.additionalCopyInfoOffset()]) & 0x7f, nil
	}
	return 0, ErrNoPESAdditionalCopyInfo
}

func (p PESHeader) crcOffset() int {
	if p.HasAdditionalCopyInfo() {
		return p.additionalCopyInfoOffset() + 1
	}
	return p.additionalCopyInfoOffset()
}

func (p PESHeader) CRC() ([]byte, error) {
	if p.HasCRC() {
		offset := p.crcOffset()
		return p[offset : offset+4], nil
	}
	return nil, ErrNoPESCRC
}

func (p PESHeader) extensionOffset() int {
	if p.HasCRC() {
		return p.crcOffset() + 4
	}
	return p.crcOffset()
}

type PESExtension []byte

func (p PESHeader) Extension() (PESExtension, error) {
	if p.HasExtension() {
		return PESExtension(p[p.extensionOffset():]), nil
	}
	return nil, ErrNoPESExtension
}

func (pe PESExtension) HasPrivateData() bool {
	return Bool(pe[0], 7)
}

func (pe PESExtension) HasPackHeaderField() bool {
	return Bool(pe[0], 6)
}

func (pe PESExtension) HasProgramPacketSequenceCounter() bool {
	return Bool(pe[0], 5)
}

func (pe PESExtension) HasPStdBufferInfo() bool {
	return Bool(pe[0], 4)
}

func (pe PESExtension) HasExtension() bool {
	return Bool(pe[0], 0)
}

func (pe PESExtension) privateDataOffset() int {
	return 1
}

func (pe PESExtension) PrivateData() ([]byte, error) {
	if pe.HasPrivateData() {
		offset := pe.privateDataOffset()
		return pe[offset : offset+16], nil
	}
	return nil, ErrNoPESPrivateData
}

func (pe PESExtension) packHeaderFieldOffset() int {
	if pe.HasPrivateData() {
		return pe.privateDataOffset() + 16
	}
	return pe.privateDataOffset()
}

func (pe PESExtension) PackHeaderField() (byte, error) {
	if pe.HasPrivateData() {
		return pe[pe.packHeaderFieldOffset()], nil
	}
	return 0x00, ErrNoPESPackHeaderField
}

func (pe PESExtension) programPacketSequenceCounterOffset() int {
	if pe.HasPackHeaderField() {
		return pe.packHeaderFieldOffset() + 1
	}
	return pe.packHeaderFieldOffset()
}

type ProgramPacketSequenceCounter uint16

func (pe PESExtension) ProgramPacketSequenceCounter() (ProgramPacketSequenceCounter, error) {
	if pe.HasProgramPacketSequenceCounter() {
		offset := pe.programPacketSequenceCounterOffset()
		return ProgramPacketSequenceCounter(Uimsbf16(pe[offset:offset+2], 16)), nil
	}
	return ProgramPacketSequenceCounter(0), ErrNoPESProgramPacketSequenceCounter
}

type PStdBufferInfo uint16

func (pe PESExtension) pStdBufferInfoOffset() int {
	if pe.HasProgramPacketSequenceCounter() {
		return pe.programPacketSequenceCounterOffset() + 2
	}
	return pe.programPacketSequenceCounterOffset()
}

func (pe PESExtension) PStdBufferInfo() (PStdBufferInfo, error) {
	if pe.HasPStdBufferInfo() {
		offset := pe.pStdBufferInfoOffset()
		return PStdBufferInfo(Uimsbf16(pe[offset:offset+2], 16)), nil
	}
	return 0, ErrNoPStdBufferInfo
}

func (pe PESExtension) extensionOffset() int {
	if pe.HasPStdBufferInfo() {
		return pe.pStdBufferInfoOffset() + 2
	}
	return pe.pStdBufferInfoOffset()
}

type PESSecondExtension uint16

func (pe PESExtension) Extension() (PESSecondExtension, error) {
	if pe.HasExtension() {
		offset := pe.extensionOffset()
		return PESSecondExtension(Uimsbf16(pe[offset:offset+2], 16)), nil
	}
	return 0, ErrNoPESSecondExtension
}

func (pe PESExtension) Length() int {
	if pe.HasExtension() {
		return pe.extensionOffset() + 2
	}
	return pe.extensionOffset()
}
