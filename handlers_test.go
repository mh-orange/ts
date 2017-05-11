package ts

import (
	"bytes"
	"encoding/hex"
	"math/rand"
	"testing"

	"github.com/Comcast/gots"
	"github.com/Comcast/gots/packet"
)

func createPayload(length int) []byte {
	// add crc
	length += 4

	payload := make([]byte, length)

	// random payload
	rand.Read(payload)

	// section length
	payload[1] = uint8(uint16(length-3) >> 8 & 0x03)
	payload[2] = uint8(0xff&length - 3)

	// compute CRC
	computedCrc := gots.ComputeCRC(payload[0 : len(payload)-4])
	copy(payload[len(payload)-4:], computedCrc)

	return payload
}

var cc uint8

func createPackets(offset int, discontinuity bool, expectedDiscontinuity bool, payload []byte) []packet.Packet {
	p := make([]byte, len(payload)+offset+1)
	copy(p[offset+1:], payload)

	// offset byte
	p[0] = uint8(offset)

	// filler
	for i := 1; i < offset+1; i++ {
		p[i] = 0xff
	}

	packets := make([]packet.Packet, 0)

	offset = 0
	for offset < len(payload) {
		if discontinuity {
			cc += 10
		}

		pusi := true
		if offset != 0 {
			pusi = false
		}

		var pkt packet.Packet
		if expectedDiscontinuity {
			pkt = packet.CreateDCPacket(1, cc)
			packet.WithHasPayloadFlag(&pkt)
			if pusi {
				packet.WithPUSI(&pkt)
			}
			pkt[4] = 0x01
		} else {
			pkt = packet.CreateTestPacket(1, cc, pusi, true)
		}

		offset += packet.SetPayload(&pkt, p[offset:])
		packets = append(packets, pkt)

		cc += 1
	}

	if offset != len(p) {
		panic("bad offset length")
	}
	return packets
}

func TestTableHandler(t *testing.T) {
	tests := []*struct {
		offset                int
		length                int
		payload               []byte
		discontinuity         bool
		expectedDiscontinuity bool
	}{
		{
			offset:                0,
			length:                10,
			discontinuity:         false,
			expectedDiscontinuity: false,
		},
		{
			offset:                0,
			length:                200,
			discontinuity:         false,
			expectedDiscontinuity: false,
		},
		{
			offset:                4,
			length:                200,
			discontinuity:         false,
			expectedDiscontinuity: false,
		},
		{
			offset:                4,
			length:                300,
			discontinuity:         true,
			expectedDiscontinuity: false,
		},
		{
			offset:                4,
			length:                300,
			discontinuity:         true,
			expectedDiscontinuity: true,
		},
	}

	expectedPayloads := 0
	packets := make([]packet.Packet, 0)

	for _, test := range tests {
		test.payload = createPayload(test.length)
		if !test.discontinuity || test.expectedDiscontinuity {
			expectedPayloads += 1
		}
		pkts := createPackets(test.offset, test.discontinuity, test.expectedDiscontinuity, test.payload)
		for _, pkt := range pkts {
			packets = append(packets, pkt)
		}
	}

	inCh := make(chan packet.Packet, len(packets))
	for _, pkt := range packets {
		inCh <- pkt
	}

	outCh := HandleTableStreams(inCh)
	close(inCh)

	foundPayloads := 0
	payloads := make([][]byte, 0)
	for payload := range outCh {
		foundPayloads++
		payloads = append(payloads, payload)
	}

	if expectedPayloads != foundPayloads {
		t.Errorf("Expected %d payloads but got %d", expectedPayloads, foundPayloads)
	}

	i := 0
	for _, test := range tests {
		if !test.discontinuity || test.expectedDiscontinuity {
			if !bytes.Equal(test.payload, payloads[i]) {
				t.Errorf("Payload %d expected\n%s\ngot\n%s\n", i, hex.Dump(test.payload), hex.Dump(payloads[i]))
			}
		} else {
			i--
		}
		i++
	}
}

func TestCrcHandler(t *testing.T) {
	tests := []struct {
		b    []byte
		good bool
	}{
		{
			b:    []byte{0xff, 0xfe, 0xfd, 0xfc, 0xac, 0x69, 0x14, 0x51},
			good: true,
		},
		{
			b:    []byte{0xfe, 0xfe, 0xfd, 0xfc, 0xac, 0x69, 0x14, 0x51},
			good: false,
		},
	}

	inCh := make(chan []byte, len(tests))
	expectedGood := 0
	for _, test := range tests {
		inCh <- test.b
		if test.good {
			expectedGood += 1
		}
	}
	outCh := HandleCrcStreams(inCh)
	close(inCh)

	foundGood := 0
	for range outCh {
		foundGood += 1
	}

	if foundGood != expectedGood {
		t.Errorf("Expected %d good packets but got %d", expectedGood, foundGood)
	}
}
