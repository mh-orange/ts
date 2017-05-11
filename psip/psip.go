package psip

import (
	"github.com/abates/orange-ts"
)

type PSIPHandler struct {
	vctCh     chan VCT
	mgtCh     chan MGT
	payloadCh <-chan []byte
}

func HandlePSIPStreams(demux ts.Demux) *PSIPHandler {
	handler := &PSIPHandler{
		payloadCh: ts.HandleTableStreams(demux.Select(BasePID)),
	}

	go handler.run()

	return handler
}

func (ph *PSIPHandler) run() {
	for payload := range ph.payloadCh {
		tableID := uint8(payload[0])
		switch tableID {
		case MGTTableID:
			if ph.mgtCh != nil {
				ph.mgtCh <- newMGT(payload)
			}
		case TVCTTableID:
			if ph.vctCh != nil {
				ph.vctCh <- newTVCT(payload)
			}
		case CVCTTableID:
			if ph.vctCh != nil {
				ph.vctCh <- newCVCT(payload)
			}
		case SVCTTableID:
			if ph.vctCh != nil {
				ph.vctCh <- newSVCT(payload)
			}
		}
	}

	switch {
	case ph.mgtCh != nil:
		close(ph.mgtCh)
	case ph.vctCh != nil:
		close(ph.vctCh)
	}
}

func (ph *PSIPHandler) SelectVCT() <-chan VCT {
	ph.vctCh = make(chan VCT)

	return ph.vctCh
}

func (ph *PSIPHandler) SelectMGT() <-chan MGT {
	ph.mgtCh = make(chan MGT)

	return ph.mgtCh
}
