package ts

import (
	"bufio"
	"bytes"
	"sync"
	"testing"

	"github.com/Comcast/gots/packet"
)

func TestSelect(t *testing.T) {
	packets := []packet.Packet{
		packet.CreateTestPacket(0, 0, true, false),
		packet.CreateTestPacket(1, 0, true, false),
		packet.CreateTestPacket(0, 1, true, false),
		packet.CreateTestPacket(1, 1, true, false),
		packet.CreateTestPacket(1, 2, true, false),
		packet.CreateTestPacket(0, 2, true, false),
		packet.CreateTestPacket(0, 3, true, false),
		packet.CreateTestPacket(1, 3, true, false),
		packet.CreateTestPacket(0, 4, true, false),
		packet.CreateTestPacket(0x1fff, 0, true, false),
		packet.CreateTestPacket(1, 4, true, false),
		packet.CreateTestPacket(1, 5, true, false),
		packet.CreateTestPacket(1, 6, true, false),
		packet.CreateTestPacket(0, 5, true, false),
		packet.CreateTestPacket(0, 6, true, false),
		packet.CreateTestPacket(0, 7, true, false),
		packet.CreateTestPacket(1, 7, true, false),
	}

	b := []byte{}

	pidCounts := make(map[uint16]int)
	pids := make(map[uint16][]packet.Packet, 0)
	for _, pkt := range packets {
		b = append(b, []byte(pkt)...)
		pid, _ := packet.Pid(pkt)
		if _, found := pids[pid]; !found {
			pids[pid] = make([]packet.Packet, 0)
			pidCounts[pid] = 0
		}

		cc, _ := packet.ContinuityCounter(pkt)
		if ok, _ := packet.IsNull(pkt); cc < 4 && !ok {
			pidCounts[pid]++
		}
	}

	buffer := bytes.NewReader(b)

	demux := NewDemux(Reader(bufio.NewReader(buffer)))
	var wg sync.WaitGroup
	var mutex sync.Mutex
	for i := 0; i < len(pids); i++ {
		wg.Add(1)
		ch := demux.Select(uint16(i))
		go func(pid uint16, ch <-chan packet.Packet) {
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
			if len(pids[pid]) > 4 {
				mutex.Lock()
				pids[pid] = pids[pid][0 : len(pids[pid])-1]
				mutex.Unlock()
			}
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
