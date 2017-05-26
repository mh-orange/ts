package ts

import (
	"bytes"
	"sync"
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

		cc := pkt.ContinuityCounter()
		if cc < 4 && !pkt.IsNull() {
			pidCounts[pid]++
		}
	}

	buffer := bytes.NewReader(b)

	demux := NewDemux(Reader(buffer))
	var wg sync.WaitGroup
	var mutex sync.Mutex
	for i := 0; i < len(pids); i++ {
		wg.Add(1)
		ch := demux.Select(uint16(i))
		go func(pid uint16, ch <-chan Packet) {
			count := 0
			for packet := range ch {
				mutex.Lock()
				pids[pid] = append(pids[pid], packet)
				mutex.Unlock()
				count++
				if count >= 4 {
					demux.Clear(pid)
				}
			}

			// sometimes an extra packet is written if the demuxer is blocked
			// at this pid's channel write when we clear the selection, this
			// should be at most one extra packet written to the channel
			mutex.Lock()
			if len(pids[pid]) > 4 {
				pids[pid] = pids[pid][0 : len(pids[pid])-1]
			}
			mutex.Unlock()
			wg.Done()
		}(uint16(i), ch)
	}

	demux.Run()
	wg.Wait()

	for pid, pidCount := range pidCounts {
		if pidCount != len(pids[pid]) {
			t.Errorf("Pid %d expected %d packets but got %d", pid, pidCount, len(pids[pid]))
		}
	}
}
