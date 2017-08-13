package psi

import (
	"github.com/mh-orange/ts"
)

type PATEntry []byte

func (pe PATEntry) Program() uint16 {
	return ts.Uimsbf16(pe[0:2], 16)
}

func (pe PATEntry) PID() uint16 {
	return ts.Uimsbf16(pe[2:4], 13) & 0x1fff
}

type PAT struct {
	Table
}

func (pat *PAT) Entries() []PATEntry {
	entries := make([]PATEntry, len(pat.Data())/4)
	for i := 0; i < len(pat.Data()); i += 4 {
		entries[i/4] = pat.Data()[i : i+4]
	}
	return entries
}
