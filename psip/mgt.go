package psip

import (
	"github.com/mh-orange/ts"
)

type TableInfo []byte

func (t TableInfo) Type() uint16 {
	return ts.Uimsbf16(t[0:2], 16)
}

func (t TableInfo) PID() uint16 {
	return ts.Uimsbf16(t[2:4], 16) & 0x1fff
}

func (t TableInfo) Version() uint8 {
	return ts.Uimsbf8(t[4]) & 0x1f
}

func (t TableInfo) TableLength() int {
	return int(ts.Uimsbf32(t[5:9], 32))
}

func (t TableInfo) DescriptorsLength() int {
	return int(ts.Uimsbf16(t[9:11], 16) & 0x0fff)
}

func (t TableInfo) Descriptors() []ts.Descriptor {
	return ts.Descriptors(t[11 : 11+t.DescriptorsLength()])
}

func (t TableInfo) Length() int {
	return 11 + t.DescriptorsLength()
}

// Master Guide Table
type MGT interface {
	NumTables() int
	Tables() []TableInfo
	Descriptors() []ts.Descriptor
}

type mgt struct {
	table
}

func newMGT(payload []byte) MGT {
	return &mgt{table(payload)}
}

func (m *mgt) NumTables() int {
	return int(ts.Uimsbf16(m.Data()[0:2], 2))
}

func (m *mgt) Tables() []TableInfo {
	tables := make([]TableInfo, m.NumTables())
	start := 2
	for i := 0; i < len(tables); i++ {
		tables[i] = TableInfo(m.Data()[start:])
		start += tables[i].Length()
	}
	return tables
}

func (m *mgt) descriptorLengthOffset() int {
	offset := 2
	for _, tableInfo := range m.Tables() {
		offset += tableInfo.Length()
	}
	return offset
}

func (m *mgt) DescriptorsLength() int {
	offset := m.descriptorLengthOffset()
	return int(ts.Uimsbf16(m.Data()[offset:offset+2], 16) & 0x0fff)
}

func (m *mgt) Descriptors() []ts.Descriptor {
	start := m.descriptorLengthOffset() + 2
	end := len(m.Data()) - 4
	return ts.Descriptors(m.Data()[start:end])
}
