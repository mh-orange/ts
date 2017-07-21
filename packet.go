package ts

import (
	"fmt"
)

const (
	PacketSize = 188
)

var (
	ErrNoPCR                  = fmt.Errorf("AdaptationField has no Program Clock Reference")
	ErrNoOPCR                 = fmt.Errorf("AdaptationField has no Original Program Clock Reference")
	ErrNoAdaptationField      = fmt.Errorf("Packet has no Adaptation Field")
	ErrNoPayload              = fmt.Errorf("Packet has no payload")
	ErrNoPESHeader            = fmt.Errorf("Packet does not contain a PES Header")
	ErrNoSplicingPoint        = fmt.Errorf("Adaptation Field has no splicing point")
	ErrNoTransportPrivateData = fmt.Errorf("Adaptation Field has no transport private data")
	ErrNoAdaptationExtension  = fmt.Errorf("Adaptation Field does not have an extension field")
	ErrNoLTW                  = fmt.Errorf("Adaptation Extension Field has no Legal Time Window")
	ErrNoPiecewiseRate        = fmt.Errorf("Adaptation Extension Field has no Piecewise Rate")
	ErrNoSeamlessSplice       = fmt.Errorf("Adaptation Extension Field has no Seamless Splice")
)

type AdaptationExtension []byte

func (ae AdaptationExtension) Length() int {
	return int(ae[0])
}

func (ae AdaptationExtension) HasLTW() bool {
	return Bool(ae[1], 7)
}

func (ae AdaptationExtension) HasPiecewiseRate() bool {
	return Bool(ae[1], 6)
}

func (ae AdaptationExtension) HasSeamlessSplice() bool {
	return Bool(ae[1], 5)
}

func (ae AdaptationExtension) LTWIsValid() (bool, error) {
	if ae.HasLTW() {
		return Bool(ae[2], 7), nil
	}
	return false, ErrNoLTW
}

func (ae AdaptationExtension) LTWOffset() (uint16, error) {
	if ae.HasLTW() {
		return Uimsbf16(ae[2:4], 15), nil
	}
	return 0, ErrNoLTW
}

func (ae AdaptationExtension) ltwOffset() int {
	return 2
}

func (ae AdaptationExtension) piecewiseRateOffset() int {
	if ae.HasLTW() {
		return ae.ltwOffset() + 2
	}
	return ae.ltwOffset()
}

func (ae AdaptationExtension) PiecewiseRate() (uint32, error) {
	if ae.HasPiecewiseRate() {
		offset := ae.piecewiseRateOffset()
		return Uimsbf32(ae[offset:offset+3], 22), nil
	}
	return 0, ErrNoPiecewiseRate
}

func (ae AdaptationExtension) seamlessSpliceOffset() int {
	if ae.HasPiecewiseRate() {
		return ae.piecewiseRateOffset() + 3
	}
	return ae.piecewiseRateOffset()
}

func (ae AdaptationExtension) SpliceType() (byte, error) {
	if ae.HasSeamlessSplice() {
		return ae[ae.seamlessSpliceOffset()] >> 4, nil
	}
	return 0x00, ErrNoSeamlessSplice
}

func (ae AdaptationExtension) DTSNextAccessUnit() (uint64, error) {
	if ae.HasSeamlessSplice() {
		return uint64(0), nil
	}
	return uint64(0), ErrNoSeamlessSplice
}

type AdaptationField []byte

func (a AdaptationField) Length() int {
	return int(a[0])
}

func (a AdaptationField) IsDiscontinuous() bool {
	return Bool(a[1], 7)
}

func (a AdaptationField) SetIsDiscontinuous(isDiscontinuous bool) {
	SetBool(&a[1], 7, isDiscontinuous)
}

func (a AdaptationField) IsRandomAccess() bool {
	return Bool(a[1], 6)
}

func (a AdaptationField) SetIsRandomAccess(isRandomAccess bool) {
	SetBool(&a[1], 6, isRandomAccess)
}

func (a AdaptationField) HasElementaryStreamPriorty() bool {
	return Bool(a[1], 5)
}

func (a AdaptationField) HasPCR() bool {
	return Bool(a[1], 4)
}

func (a AdaptationField) HasOPCR() bool {
	return Bool(a[1], 3)
}

func (a AdaptationField) HasSplicingPoint() bool {
	return Bool(a[1], 2)
}

func (a AdaptationField) HasTransportPrivateData() bool {
	return Bool(a[1], 1)
}

func (a AdaptationField) HasAdaptationExtension() bool {
	return Bool(a[1], 0)
}

func (a AdaptationField) PCR() (TimeStamp, error) {
	if a.HasPCR() {
		return CalculatePCR(a[2:]), nil
	}
	return nil, ErrNoPCR
}

func (a AdaptationField) pcrOffset() int {
	return 2
}

func (a AdaptationField) opcrOffset() int {
	if a.HasPCR() {
		return a.pcrOffset() + 6
	}
	return a.pcrOffset()
}

func (a AdaptationField) OPCR() (TimeStamp, error) {
	if a.HasOPCR() {
		return CalculatePCR(a[a.opcrOffset():]), nil
	}
	return nil, ErrNoOPCR
}

func (a AdaptationField) spliceCountdownOffset() int {
	if a.HasOPCR() {
		return a.opcrOffset() + 6
	}
	return a.opcrOffset()
}

func (a AdaptationField) SpliceCountdown() (int, error) {
	if a.HasSplicingPoint() {
		return int(a[a.spliceCountdownOffset()]), nil
	}
	return 0, ErrNoSplicingPoint
}

func (a AdaptationField) transportPrivateDataOffset() int {
	if a.HasSplicingPoint() {
		return a.spliceCountdownOffset() + 1
	}
	return a.spliceCountdownOffset()
}

func (a AdaptationField) TransportPrivateDataLength() (int, error) {
	if a.HasTransportPrivateData() {
		return int(a[a.transportPrivateDataOffset()]), nil
	}
	return 0, ErrNoTransportPrivateData
}

func (a AdaptationField) TransportPrivateData() ([]byte, error) {
	if a.HasTransportPrivateData() {
		offset := a.transportPrivateDataOffset() + 1
		length, _ := a.TransportPrivateDataLength()
		return a[offset : offset+length], nil
	}
	return nil, ErrNoTransportPrivateData
}

func (a AdaptationField) adaptationExtensionOffset() int {
	if a.HasTransportPrivateData() {
		return a.transportPrivateDataOffset() + int(a[a.transportPrivateDataOffset()]) + 1
	}
	return a.transportPrivateDataOffset()
}

func (a AdaptationField) AdaptationExtension() (AdaptationExtension, error) {
	if a.HasAdaptationExtension() {
		return AdaptationExtension(a[a.adaptationExtensionOffset():]), nil
	}
	return nil, ErrNoAdaptationExtension
}

type Packet []byte

func NewPacket() Packet {
	p := make(Packet, PacketSize)
	p[0] = 0x47
	return p
}

func (p Packet) SyncByte() byte {
	return p[0]
}

func (p Packet) TEI() bool {
	return Bool(p[1], 0)
}

func (p Packet) PUSI() bool {
	return Bool(p[1], 1)
}

func (p Packet) SetPUSI(pusi bool) {
	SetBool(&p[1], 1, pusi)
}

func (p Packet) TransportPriority() bool {
	return Bool(p[1], 2)
}

func (p Packet) PID() uint16 {
	return Uimsbf16(p[1:3], 13)
}

func (p Packet) SetPID(pid uint16) {
	SetUimsbf16(p[1:3], 13, pid)
}

func (p Packet) TSC() byte {
	return p[3] >> 6
}

func (p Packet) Scrambled() bool {
	return p.TSC() == 0x00
}

func (p Packet) ScrambledEven() bool {
	return p.TSC() == 0x02
}

func (p Packet) ScrambledOdd() bool {
	return p.TSC() == 0x03
}

func (p Packet) HasAdaptationField() bool {
	return Bool(p[3], 2)
}

func (p Packet) SetHasAdaptationField(hasAdaptationField bool) {
	SetBool(&p[3], 2, hasAdaptationField)
	SetUimsbf8(&p[4], 1)
}

func (p Packet) HasPayload() bool {
	return Bool(p[3], 3)
}

func (p Packet) SetHasPayload(hasPayload bool) {
	SetBool(&p[3], 3, hasPayload)
}

func (p Packet) ContinuityCounter() uint8 {
	return uint8(p[3] >> 4)
}

func (p Packet) SetContinuityCounter(cc uint8) {
	p[3] = (cc << 4) | p[3]&0x0f
}

func (p Packet) AdaptationField() (AdaptationField, error) {
	if p.HasAdaptationField() {
		return AdaptationField(p[4:]), nil
	}
	return nil, ErrNoAdaptationField
}

func (p Packet) payload() []byte {
	if p.HasAdaptationField() {
		return p[5+int(p[4]):]
	}
	return p[4:]
}

func (p Packet) Payload() ([]byte, error) {
	if p.HasPayload() {
		return p.payload(), nil
	}
	return nil, ErrNoPayload
}

func (p Packet) SetPayload(payload []byte) int {
	p.SetHasPayload(len(payload) > 0)
	return copy(p.payload(), payload)
}

func (p Packet) IsNull() bool {
	return p.PID() == 0x1fff
}

func (p Packet) HasPESHeader() bool {
	if p.PUSI() {
		if payload, err := p.Payload(); err == nil {
			return payload[0] == 0x00 && payload[1] == 0x00 && payload[2] == 0x01
		}
	}
	return false
}

func (p Packet) PESHeader() (PESHeader, error) {
	if p.HasPESHeader() {
		return PESHeader(p.payload()), nil
	}
	return nil, ErrNoPESHeader
}
