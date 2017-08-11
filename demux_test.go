package ts

import (
	"bytes"
	"testing"
)

func createTestPacket(pid uint16, cc uint8, pusi bool) Packet {
	p := NewPacket()
	p.SetPID(pid)
	p.SetContinuityCounter(cc)
	p.SetPUSI(pusi)
	p.SetHasPayload(true)
	return p
}

func TestSelect(t *testing.T) {
	packets := []Packet{
		createTestPacket(0, 0, true),
		createTestPacket(1, 0, true),
		createTestPacket(0, 1, true),
		createTestPacket(1, 1, true),
		createTestPacket(1, 2, true),
		createTestPacket(0, 2, true),
		createTestPacket(0, 3, true),
		createTestPacket(1, 3, true),
		createTestPacket(0, 4, true),
		createTestPacket(0x1fff, 0, true),
		createTestPacket(1, 4, true),
		createTestPacket(1, 5, true),
		createTestPacket(1, 6, true),
		createTestPacket(0, 5, true),
		createTestPacket(0, 6, true),
		createTestPacket(0, 7, true),
		createTestPacket(1, 7, true),
	}

	b := []byte{}

	pidCounts := make(map[uint16]int)
	pids := make(map[uint16][]Packet, 0)
	for _, pkt := range packets {
		b = append(b, []byte(pkt)...)
		pid := pkt.PID()
		if _, found := pids[pid]; !found {
			pids[pid] = make([]Packet, 0)
			pidCounts[pid] = 0
		}

		if !pkt.IsNull() {
			pidCounts[pid]++
		}
	}

	buffer := bytes.NewReader(b)

	demux := NewDemux(NewReader(buffer))
	for i := 0; i < len(pids); i++ {
		pid := uint16(i)
		demux.Select(pid, PacketHandlerFunc(func(pkt Packet) {
			pids[pid] = append(pids[pid], pkt)
		}))
	}

	demux.Run()

	for pid, pidCount := range pidCounts {
		if pidCount != len(pids[pid]) {
			t.Errorf("Pid %d expected %d packets but got %d", pid, pidCount, len(pids[pid]))
		}
	}
}

func TestClear(t *testing.T) {
	buffer := bytes.NewReader(nil)
	d := NewDemux(NewReader(buffer)).(*demux)
	d.Select(42, PacketHandlerFunc(func(pkt Packet) {}))
	if _, ok := d.handlers[42]; !ok {
		t.Errorf("Select should have added a channel to the channels map")
	}

	d.Clear(42)
	if _, ok := d.handlers[42]; ok {
		t.Errorf("Clear should have removed a channel to the channels map")
	}
}
