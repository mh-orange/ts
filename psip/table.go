package psip

import (
	"github.com/mh-orange/ts"
	"github.com/mh-orange/ts/psi"
)

type Table struct {
	psi.Table
}

func (t *Table) ProtocolVersion() uint8 {
	return ts.Uimsbf8(t.Table.Data()[0])
}

func (t *Table) Data() []byte {
	return t.Table.Data()[1:]
}
