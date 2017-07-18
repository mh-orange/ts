package ts

import (
	"testing"
	"time"
)

type timestampTest struct {
	input    []byte
	ticks    uint64
	duration time.Duration
}

func TestCalculatePCR(t *testing.T) {
	tests := []timestampTest{
		{
			input:    []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			ticks:    0,
			duration: time.Duration(0),
		},
		{
			input:    []byte{0x00, 0x00, 0x00, 0x00, 0x80, 0x00},
			ticks:    300,
			duration: time.Duration(11 * time.Microsecond),
		},
	}

	for i, test := range tests {
		pcr := CalculatePCR(test.input)
		if pcr.Ticks() != test.ticks {
			t.Errorf("Test %d failed.  Expected %d got %d", i, test.ticks, pcr.Ticks())
		}

		if pcr.Duration().Truncate(time.Microsecond) != test.duration {
			t.Errorf("Test %d failed.  Expected %s got %s", i, test.duration, pcr.Duration().Truncate(time.Microsecond))
		}
	}
}

func TestCalculatePTS(t *testing.T) {
	tests := []timestampTest{
		{
			input:    []byte{0xf1, 0x00, 0x01, 0x00, 0x01},
			ticks:    0,
			duration: time.Duration(0),
		},
		{
			input:    []byte{0xf1, 0x00, 0x01, 0x00, 0x03},
			ticks:    1,
			duration: 11111 * time.Nanosecond,
		},
		{
			input:    []byte{0xf1, 0x00, 0x05, 0xbf, 0x21},
			ticks:    90000,
			duration: 90000 * 11111 * time.Nanosecond,
		},
		{
			input:    []byte{0xff, 0xff, 0xff, 0xff, 0xff},
			ticks:    8589934591,
			duration: 8589934591 * 11111 * time.Nanosecond,
		},
	}

	for i, test := range tests {
		pts := CalculatePTS(test.input)
		if pts.Ticks() != test.ticks {
			t.Errorf("Test %d failed.  Expected %d got %d", i, test.ticks, pts.Ticks())
		}

		if pts.Duration() != test.duration {
			t.Errorf("Test %d failed.  Expected %s got %s", i, test.duration, pts.Duration())
		}
	}
}

func TestDelta(t *testing.T) {
	tests := []struct {
		startTicks uint64
		endTicks   uint64
		duration   time.Duration
	}{
		{0, 1, time.Second},
		{0x3fffffffff, 0, time.Second},
	}

	for i, test := range tests {
		startTs := &timestamp{test.startTicks, uint64(time.Second)}
		endTs := &timestamp{test.endTicks, uint64(time.Second)}

		ts := startTs.Delta(endTs)
		if test.duration != ts.Duration() {
			t.Errorf("Test %d expected %s but got %s", i, test.duration, ts.Duration())
		}
	}
}
