package ts

import (
	"time"
)

const (
	// Presentation Time Stamp (PTS) is based on a 90khz clock
	PTSFrequency = uint64(90000)
)

type TimeStamp interface {
	Ticks() uint64
	Increment(ticks uint64)
	Duration() time.Duration
	Delta(endTs TimeStamp) TimeStamp
}

func CalculatePTS(b []byte) TimeStamp {
	ticks := uint64(b[0]>>1&0x07) << 30
	ticks |= uint64(b[1]) << 22
	ticks |= uint64(b[2]>>1&0x7f) << 15
	ticks |= uint64(b[3]) << 7
	ticks |= uint64(b[4] >> 1 & 0x7f)

	pts := NewPTS()
	pts.Increment(ticks)
	return pts
}

func NewPTS() TimeStamp {
	return newTimestamp(0, PTSFrequency)
}

func newTimestamp(ticks uint64, frequency uint64) *timestamp {
	return &timestamp{ticks, 1000000000 / frequency}
}

type timestamp struct {
	ticks    uint64
	interval uint64 // number of nanoseconds in one period
}

func (t *timestamp) Increment(ticks uint64) {
	t.ticks += ticks
}

func (t *timestamp) Ticks() uint64 {
	return t.ticks
}

func (t *timestamp) Duration() time.Duration {
	return time.Duration(t.ticks * t.interval)
}

func (t *timestamp) Delta(endTs TimeStamp) TimeStamp {
	return &timestamp{0x3fffffffff & (endTs.Ticks() - t.ticks), t.interval}
}
