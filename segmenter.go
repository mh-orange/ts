package ts

import (
	"bytes"
	"time"
)

type Segment struct {
	Duration time.Duration
	Buffer   []byte
}

func hasPts(pkt Packet) (header PESHeader, found bool) {
	if header, err := pkt.PESHeader(); err == nil {
		found = header.HasPTS()
	}
	return
}

func getPts(pkt Packet) TimeStamp {
	pts := NewPTS()
	if _, ok := hasPts(pkt); ok {
		pay, _ := pkt.Payload()
		pts = CalculatePTS(pay[9:14])
	}
	return pts
}

func hasRAI(pkt Packet) bool {
	if field, err := pkt.AdaptationField(); err == nil {
		return field.IsRandomAccess()
	}
	return false
}

func emit(buffer *bytes.Buffer, duration time.Duration, outCh chan Segment) {
	if buffer.Len() > 0 {
		output := make([]byte, buffer.Len())
		copy(output, buffer.Bytes())
		outCh <- Segment{
			Duration: duration,
			Buffer:   output,
		}
		buffer.Reset()
	}
}

func segmentStream(inCh <-chan Packet, outCh chan Segment) {
	buffer := bytes.NewBuffer([]byte{})
	startPts := NewPTS()
	endPts := NewPTS()
	i := 0
	for pkt := range inCh {
		if _, ok := hasPts(pkt); ok {
			if startPts.Ticks() == 0 && endPts.Ticks() == 0 {
				startPts = getPts(pkt)
			}
		}

		// Segments must start with a keyframe.  The RAI bit in the TS packet
		// indicates whether a video stream can be decoded without error
		// starting at that packet, so we start each segment based on the RAI bit
		//
		// If segments need to be produced for specific time intervals, then an upstream
		// element in the pipeline should be transcoding and inserting keyframes at
		// specific time intervals
		if hasRAI(pkt) {
			emit(buffer, startPts.Delta(endPts).Duration(), outCh)
			startPts = endPts
		}

		buffer.Write(pkt)
		if _, ok := hasPts(pkt); ok {
			endPts = getPts(pkt)
		}
		i++
	}
	emit(buffer, startPts.Delta(endPts).Duration(), outCh)
	close(outCh)
}

func SegmentStream(inCh <-chan Packet) <-chan Segment {
	outCh := make(chan Segment)
	go segmentStream(inCh, outCh)
	return outCh
}
