package psip

type cvct struct {
	*tvct
}

func newCVCT(payload []byte) VCT {
	return &cvct{&tvct{table(payload)}}
}
