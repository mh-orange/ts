package ts

import (
	"sync"

	"github.com/Comcast/gots/packet"
)

type Demux interface {
	Select(pid uint16) <-chan packet.Packet
	Clear(pid uint16)
	Run()
}

type demux struct {
	inCh       <-chan packet.Packet
	channels   map[uint16]chan packet.Packet
	chMutex    sync.Mutex
	clearPids  []uint16
	clearMutex sync.Mutex
}

func NewDemux(inCh <-chan packet.Packet) Demux {
	return &demux{
		inCh:     inCh,
		channels: make(map[uint16]chan packet.Packet),
	}
}

func (d *demux) Select(pid uint16) <-chan packet.Packet {
	ch := make(chan packet.Packet)
	d.chMutex.Lock()
	d.channels[pid] = ch
	d.chMutex.Unlock()
	return ch
}

func (d *demux) Clear(pid uint16) {
	d.clearMutex.Lock()
	d.clearPids = append(d.clearPids, pid)
	d.clearMutex.Unlock()
}

func (d *demux) Run() {
	for pkt := range d.inCh {
		if ok, _ := packet.IsNull(pkt); ok {
			continue
		}

		pid, _ := packet.Pid(pkt)
		d.chMutex.Lock()
		if channel, found := d.channels[pid]; found {
			channel <- pkt
		}
		d.chMutex.Unlock()

		d.clearMutex.Lock()
		for _, pid := range d.clearPids {
			d.chMutex.Lock()
			if ch, found := d.channels[pid]; found {
				close(ch)
				delete(d.channels, pid)
			}
			d.chMutex.Unlock()
		}
		d.clearMutex.Unlock()
	}

	d.chMutex.Lock()
	for _, ch := range d.channels {
		close(ch)
	}
	d.chMutex.Unlock()
}
