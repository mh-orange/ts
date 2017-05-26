package ts

import (
	"bytes"
)

type tableStreamHandler struct {
	buffer []byte
	inCh   <-chan Packet
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
		nextCc := pkt.ContinuityCounter()
		if !first && cc+1 != nextCc {
			// check for expected discontinuity
			if field, err := pkt.AdaptationField(); err == nil {
				discontinuity = !field.IsDiscontinuous()
			}
		}
		first = false
		cc = nextCc

		if pkt.HasPayload() {
			payload, _ = pkt.Payload()
		}

		// check payload unit start indicator
		if pkt.PUSI() {
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

func HandleTableStreams(inCh <-chan Packet) <-chan []byte {
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
			crc := ComputeCRC(buffer[0 : len(buffer)-4])
			if bytes.Equal(crc, buffer[len(buffer)-4:]) {
				outCh <- buffer
			}
		}
		close(outCh)
	}(inCh, outCh)
	return outCh
}
