package psip

type CVCT struct {
	*TVCT
}

func newCVCT(payload []byte) VCT {
	return &CVCT{&TVCT{&Table{payload}}}
}
