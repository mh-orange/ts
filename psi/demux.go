package psi

import (
	"bytes"
	"sync"

	"github.com/mh-orange/ts"
)

type TableHandler interface {
	Handle(table Table)
}

type TableHandlerFunc func(table Table)

func (thf TableHandlerFunc) Handle(table Table) {
	thf(table)
}

type TableDemux interface {
	Handle(pkt ts.Packet)
	Clear(id uint8)
	Select(id uint8, handler TableHandler)
}

type TableBuffer interface {
	Append([]byte)
	Flush()
	HasNext() bool
	Next() Table
}

type defaultTableBuffer struct {
	buffer []byte
	tables []Table
}

func sectionLength(data []byte) (length int) {
	if len(data) > 3 && data[0] != 0xff {
		length = 3 + int(ts.Uimsbf16(data[1:3], 16))
	}
	return length
}

func (tb *defaultTableBuffer) Next() (table Table) {
	if len(tb.tables) > 0 {
		table = tb.tables[0]
		tb.tables = tb.tables[1:]
	}
	return
}

func (tb *defaultTableBuffer) HasNext() bool {
	return len(tb.tables) > 0
}

func (tb *defaultTableBuffer) Flush() {
	tb.buffer = nil
}

func (tb *defaultTableBuffer) Append(data []byte) {
	tb.buffer = append(tb.buffer, data...)
	// stopping condition is end of the buffer (do nothing) or
	// table id of 0xff (flush remaining bytes)
	for len(tb.buffer) > 0 && tb.buffer[0] != 0xff {
		length := sectionLength(tb.buffer)
		if length > 0 && length < len(tb.buffer) {
			table := make(Table, length)
			copy(table, tb.buffer[0:length])
			tb.tables = append(tb.tables, table)
			tb.buffer = tb.buffer[length:]
		} else {
			break
		}
	}

	if len(tb.buffer) > 0 && tb.buffer[0] == 0xff {
		tb.buffer = nil
	}
}

type tableDemux struct {
	buffer   TableBuffer
	mu       sync.Mutex
	handlers map[uint8]TableHandler
	first    bool
	cc       uint8
}

func NewTableDemux() TableDemux {
	return &tableDemux{
		buffer:   &defaultTableBuffer{},
		handlers: make(map[uint8]TableHandler),
		first:    true,
	}
}

func (td *tableDemux) Select(tableID uint8, handler TableHandler) {
	td.mu.Lock()
	td.handlers[tableID] = handler
	td.mu.Unlock()
}

func (td *tableDemux) Clear(tableID uint8) {
	td.mu.Lock()
	delete(td.handlers, tableID)
	td.mu.Unlock()
}

func (td *tableDemux) Handle(pkt ts.Packet) {
	var payload []byte
	discontinuous := false
	// check continuity counter
	nextCc := pkt.ContinuityCounter()
	if !td.first && td.cc+1 != nextCc {
		// check for expected discontinuity
		if field, err := pkt.AdaptationField(); err == nil {
			discontinuous = !field.IsDiscontinuous()
		}
	}
	td.first = false
	td.cc = nextCc

	if pkt.HasPayload() {
		payload, _ = pkt.Payload()
	}

	// If the Payload Unit Start Indicator (PUSI) is set, then a new PSI table
	// begins in the payload section of the packet.  If the PUSI is not set,
	// then the payload is a continuation of a previous packet's payload.  This
	// can simply be appended to the current buffer as long as there is no
	// discontinuity
	//
	// When PUSI is set it doesn't matter if there is a discontinuity since we're
	// starting with a new table
	if pkt.PUSI() {
		td.buffer.Append(payload[int(payload[0])+1:])
	} else if !discontinuous {
		td.buffer.Append(payload)
	} else {
		td.buffer.Flush()
	}

	for td.buffer.HasNext() {
		table := td.buffer.Next()
		crc := ts.ComputeCRC(table[0 : len(table)-4])
		if !bytes.Equal(crc, table[len(table)-4:]) {
			continue
		}
		td.mu.Lock()
		if handler, ok := td.handlers[table.ID()]; ok {
			handler.Handle(table)
		}
		td.mu.Unlock()
	}
}
