package ts

import (
	"bytes"

	"github.com/Comcast/gots"
	"github.com/Comcast/gots/packet"
	"github.com/Comcast/gots/packet/adaptationfield"
)

type tableStreamHandler struct {
	buffer []byte
	inCh   <-chan packet.Packet
	outCh  chan []byte
}

func (tsh *tableStreamHandler) emit() {
	if len(tsh.buffer) > 0 {
		start := 1 + int(tsh.buffer[0])
		sectionLength := int(tsh.buffer[start+1]&0x0f)<<8 | int(tsh.buffer[start+2])
		end := start + sectionLength + 3
		if end < len(tsh.buffer) {
			tsh.outCh <- tsh.buffer[start:end]
		}
	}
}

func (tsh *tableStreamHandler) run() {
	first := true
	var cc uint8
	var payload []byte

	count := 0
	for pkt := range tsh.inCh {
		discontinuity := false
		count++
		// check continuity counter
		nextCc, _ := packet.ContinuityCounter(pkt)
		if !first && cc+1 != nextCc {
			discontinuity = true
			// check for expected discontinuity
			if ok, _ := packet.ContainsAdaptationField(pkt); ok {
				discontinuity = !adaptationfield.IsDiscontinuous(pkt)
			}
		}
		first = false
		cc = nextCc

		if ok, _ := packet.ContainsPayload(pkt); ok {
			payload, _ = packet.Payload(pkt)
		}

		// check payload unit start indicator
		if ok, _ := packet.PayloadUnitStartIndicator(pkt); ok {
			tsh.emit()
			tsh.buffer = make([]byte, len(payload))
			copy(tsh.buffer, payload)
		} else if !discontinuity {
			tsh.buffer = append(tsh.buffer, payload...)
		} else {
			tsh.buffer = make([]byte, 0)
		}
	}

	tsh.emit()
	close(tsh.outCh)
}

func HandleTableStreams(inCh <-chan packet.Packet) <-chan []byte {
	outCh := make(chan []byte)
	handler := &tableStreamHandler{
		inCh:  inCh,
		outCh: outCh,
	}
	go handler.run()
	return HandleCrcStreams(outCh)
}

func HandleCrcStreams(inCh <-chan []byte) <-chan []byte {
	outCh := make(chan []byte)
	go func(inCh <-chan []byte, outCh chan []byte) {
		for buffer := range inCh {
			crc := gots.ComputeCRC(buffer[0 : len(buffer)-4])
			if bytes.Equal(crc, buffer[len(buffer)-4:]) {
				outCh <- buffer
			}
		}
		close(outCh)
	}(inCh, outCh)
	return outCh
}
