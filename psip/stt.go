package psip

import (
	"time"

	"github.com/mh-orange/ts"
)

type DSControl []byte

func (dsc DSControl) Status() bool {
	return ts.Bool(dsc[0], 0)
}

func (dsc DSControl) DayOfMonth() int {
	return int(0x1f & dsc[0])
}

func (dsc DSControl) Hour() int {
	return int(dsc[1])
}

type STT struct {
	*Table
}

var (
	baseTime = time.Date(1980, time.January, 6, 0, 0, 0, 0, time.UTC)
)

func newSTT(payload []byte) *STT {
	return &STT{&Table{payload}}
}

func (s *STT) SystemTime() time.Time {
	return baseTime.Add(time.Duration(ts.Uimsbf32(s.Data()[0:4], 32)) * time.Second)
}

func (s *STT) GPSOffset() uint8 {
	return ts.Uimsbf8(s.Data()[4])
}

func (s *STT) DaylightSaving() DSControl {
	return DSControl(s.Data()[5:7])
}
