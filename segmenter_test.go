package ts

import (
	"testing"
	"time"

	"github.com/Comcast/gots"
	"github.com/Comcast/gots/packet"
)

func getPacket(pusi bool, rai bool, pts uint64) packet.Packet {
	pkt := packet.Create(1)
	if pusi {
		packet.WithPUSI(&pkt)
	}

	packet.WithHasPayloadFlag(&pkt)
	if rai {
		packet.WithHasAdaptationFieldFlag(&pkt)
		pkt[4] = 0x01
		pkt[5] |= 0x40
	}

	if pusi {
		pay, _ := packet.Payload(pkt)
		pay[0] = 0x00
		pay[1] = 0x00
		pay[2] = 0x01
		pay[3] = 0xb8
		pay[4] = 0x00
		pay[5] = 0x08
		pay[6] = 0x20
		pay[7] = 0x80
		pay[8] = 0x05

		gots.InsertPTS(pay[9:14], pts)
	}
	return pkt
}

func TestSegmenter(t *testing.T) {
	clock := uint64(0)
	numSegments := 10
	interval := 10
	numPackets := numSegments * interval
	inCh := make(chan packet.Packet, numPackets)

	for i := 0; i < numSegments; i++ {
		// emit 1 pps for "interval" seconds
		for j := 0; j < interval; j++ {
			clock += PTSFrequency
			inCh <- getPacket(true, j == 0, clock)
		}
	}
	close(inCh)

	rxSegments := 0
	expected := time.Duration(uint64(interval) * PTSFrequency * (1000000000 / PTSFrequency))
	for segment := range SegmentStream(inCh) {
		if segment.Duration != expected {
			t.Errorf("Segment %d expected to be duration %s but got %s", rxSegments, expected, segment.Duration)
		}
		rxSegments++
	}

	if rxSegments != numSegments {
		t.Errorf("Expected %d segments but got %d", numSegments, rxSegments)
	}
}
