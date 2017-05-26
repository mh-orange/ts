package ts

import (
	"math"
	"testing"
	"time"
)

func getPacket(pusi bool, rai bool, pts uint64) Packet {
	pkt := NewPacket()
	pkt.SetPID(1)
	pkt.SetPUSI(pusi)
	pkt.SetHasPayload(true)

	if rai {
		pkt.SetHasAdaptationField(true)
		field, _ := pkt.AdaptationField()
		field.SetIsRandomAccess(true)
	}

	if pusi {
		payload, _ := pkt.Payload()
		header := FillPESHeader(payload)
		header.SetHasPTS(true)
		header.SetPTS(newTimestamp(pts, 90000))
	}
	return pkt
}

func TestSegmenter(t *testing.T) {
	clock := uint64(0)
	numSegments := 10
	interval := 10
	numPackets := numSegments * interval
	inCh := make(chan Packet, numPackets)

	for i := 0; i < numSegments; i++ {
		// emit 1 pps for "interval" seconds
		for j := 0; j < interval; j++ {
			clock += PTSFrequency
			inCh <- getPacket(true, j == 0, clock)
		}
	}
	close(inCh)

	rxSegments := 0
	expected := 10 * time.Second
	for segment := range SegmentStream(inCh) {
		received := time.Second * time.Duration(math.Floor(segment.Duration.Seconds()+0.5))
		if received != expected {
			t.Errorf("Segment %d expected to be duration %s but got %s", rxSegments, expected, received)
		}
		rxSegments++
	}

	if rxSegments != numSegments {
		t.Errorf("Expected %d segments but got %d", numSegments, rxSegments)
	}
}
