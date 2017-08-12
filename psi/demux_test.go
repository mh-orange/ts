package psi

import (
	"bytes"
	"encoding/hex"
	"math/rand"
	"testing"

	"github.com/mh-orange/ts"
)

func createTable(crc bool, length int) Table {
	// add crc
	length += 4

	payload := make([]byte, length)

	// random payload
	rand.Read(payload)

	// table ID
	payload[0] = uint8(1)

	// section length
	payload[1] = uint8(uint16(length-3) >> 8 & 0x03)
	payload[2] = uint8(0xff&length - 3)

	// compute CRC
	if crc {
		computedCrc := ts.ComputeCRC(payload[0 : len(payload)-4])
		copy(payload[len(payload)-4:], computedCrc)
	}

	return Table(payload)
}

var cc uint8

func createPackets(offset int, discontinuity bool, expectedDiscontinuity bool, table Table) []ts.Packet {
	p := make([]byte, len(table)+offset+1)
	copy(p[offset+1:], table)

	// offset byte
	p[0] = uint8(offset)

	// filler
	for i := 1; i < offset+1; i++ {
		p[i] = 0xff
	}

	packets := make([]ts.Packet, 0)

	offset = 0
	for offset < len(table) {
		if discontinuity {
			cc += 10
		}
		pkt := ts.NewPacket()

		pkt.SetContinuityCounter(cc)
		pkt.SetPUSI(offset == 0)
		pkt.SetHasPayload(true)

		if expectedDiscontinuity {
			pkt.SetHasAdaptationField(true)
			field, _ := pkt.AdaptationField()
			field.SetIsDiscontinuous(true)
		}

		length := pkt.SetPayload(p[offset:])
		offset += length
		payload, _ := pkt.Payload()

		// set the rest of the packet to 0xff for stuffing bits
		for i := length; i < len(payload); i++ {
			payload[i] = 0xff
		}
		packets = append(packets, pkt)

		cc += 1
	}

	if offset != len(p) {
		panic("bad offset length")
	}
	return packets
}

func TestTableDemux(t *testing.T) {
	tests := []*struct {
		offset                int
		length                int
		table                 Table
		discontinuity         bool
		expectedDiscontinuity bool
		crc                   bool
	}{
		{
			offset:                0,
			length:                10,
			discontinuity:         false,
			expectedDiscontinuity: false,
			crc: true,
		},
		{
			offset:                0,
			length:                200,
			discontinuity:         false,
			expectedDiscontinuity: false,
			crc: true,
		},
		{
			offset:                4,
			length:                200,
			discontinuity:         false,
			expectedDiscontinuity: false,
			crc: true,
		},
		{
			offset:                4,
			length:                300,
			discontinuity:         true,
			expectedDiscontinuity: false,
			crc: true,
		},
		{
			offset:                4,
			length:                300,
			discontinuity:         true,
			expectedDiscontinuity: true,
			crc: true,
		},
		{
			offset:                0,
			length:                10,
			discontinuity:         false,
			expectedDiscontinuity: false,
			crc: false,
		},
	}

	expectedTables := 0
	packets := make([]ts.Packet, 0)

	for _, test := range tests {
		test.table = createTable(test.crc, test.length)
		if test.crc && (!test.discontinuity || test.expectedDiscontinuity) {
			expectedTables += 1
		}
		pkts := createPackets(test.offset, test.discontinuity, test.expectedDiscontinuity, test.table)
		for _, pkt := range pkts {
			packets = append(packets, pkt)
		}
	}

	demuxer := NewTableDemux()
	foundTables := 0
	tables := make([]Table, 0)
	demuxer.Select(1, TableHandlerFunc(func(table Table) {
		foundTables++
		tables = append(tables, table)
	}))

	for _, pkt := range packets {
		demuxer.Handle(pkt)
	}

	if expectedTables != foundTables {
		t.Errorf("Expected %d payloads but got %d", expectedTables, foundTables)
	}

	i := 0
	for _, test := range tests {
		if test.crc && (!test.discontinuity || test.expectedDiscontinuity) {
			if !bytes.Equal(test.table, tables[i]) {
				t.Errorf("Payload %d expected\n%s\ngot\n%s\n", i, hex.Dump(test.table), hex.Dump(tables[i]))
			}
		} else {
			i--
		}
		i++
	}
}

func TestClear(t *testing.T) {
	d := NewTableDemux().(*tableDemux)
	d.Select(42, TableHandlerFunc(func(table Table) {}))
	if _, ok := d.handlers[42]; !ok {
		t.Errorf("Select should have added a channel to the channels map")
	}

	d.Clear(42)
	if _, ok := d.handlers[42]; ok {
		t.Errorf("Clear should have removed a channel to the channels map")
	}
}
