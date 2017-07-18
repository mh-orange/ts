package psip

type cvct struct {
	*tvct
}

func newCVCT(payload []byte) VCT {
	vct := tvct(payload)
	return &cvct{&vct}
}
