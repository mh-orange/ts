package psi

import (
	"bytes"
	"encoding/hex"
	"math/rand"
	"testing"

	"github.com/mh-orange/ts"
)

func createTable(id int, crc bool, length int) Table {
	// random payload
	payload := make([]byte, length)
	rand.Read(payload)

	table := CreateTable(uint8(id), payload)

	if !crc {
		table[len(table)-4] = 0x00
		table[len(table)-3] = 0x00
		table[len(table)-2] = 0x00
		table[len(table)-1] = 0x00
	}

	return table
}

var cc uint8

func createPackets(offset int, discontinuity bool, expectedDiscontinuity bool, table Table) []ts.Packet {
	payload := make([]byte, len(table)+offset+1)
	copy(payload[offset+1:], table)

	// offset byte
	payload[0] = uint8(offset)

	// filler
	for i := 1; i < offset+1; i++ {
		payload[i] = 0xff
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
		pkt.SetHasPayload()

		if expectedDiscontinuity {
			pkt.SetHasAdaptationField()
			field, _ := pkt.AdaptationField()
			field.SetIsDiscontinuous()
		}

		length := pkt.SetPayload(payload[offset:])
		offset += length
		p, _ := pkt.Payload()

		// set the rest of the packet to 0xff for stuffing bits
		for i := length; i < len(p); i++ {
			p[i] = 0xff
		}
		packets = append(packets, pkt)

		cc += 1
	}

	if offset != len(payload) {
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
	receivedTables := make([]Table, 0)
	packets := make([]ts.Packet, 0)

	demuxer := NewTableDemux()
	foundTables := 0
	i := 0

	for j, test := range tests {
		test.table = createTable(j, test.crc, test.length)
		if test.crc && (!test.discontinuity || test.expectedDiscontinuity) {
			expectedTables += 1
		}

		pkts := createPackets(test.offset, test.discontinuity, test.expectedDiscontinuity, test.table)
		for _, pkt := range pkts {
			packets = append(packets, pkt)
		}

		demuxer.Select(uint8(j), TableHandlerFunc(func(table Table) {
			foundTables++
			i++
			receivedTables = append(receivedTables, table)
		}))
	}

	for _, pkt := range packets {
		demuxer.Handle(pkt)
	}

	if expectedTables != foundTables {
		t.Errorf("Expected %d tables but got %d", expectedTables, foundTables)
	}

	i = 0
	for j, test := range tests {
		if test.crc && (!test.discontinuity || test.expectedDiscontinuity) {
			if i < len(receivedTables) && !bytes.Equal(test.table, receivedTables[i]) {
				t.Errorf("Test %d: Table %d expected\n%s\ngot\n%s\n", j, i, hex.Dump(test.table), hex.Dump(receivedTables[i]))
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
