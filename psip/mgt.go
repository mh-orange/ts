package psip

type mgt []byte

func newMGT(b []byte) MGT {
	m := make(mgt, 0)
	return &m
}

func (m mgt) Crc() []byte {
	if len(m) < 4 {
		return nil
	}
	return m[len(m)-4:]
}
