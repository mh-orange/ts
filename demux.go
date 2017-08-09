package ts

import (
	"sync"
)

type Demux interface {
	Select(pid uint16) <-chan Packet
	Clear(pid uint16)
	Run()
}

type demux struct {
	inCh     <-chan Packet
	channels map[uint16]chan Packet
	chMutex  sync.Mutex
}

func NewDemux(inCh <-chan Packet) Demux {
	return &demux{
		inCh:     inCh,
		channels: make(map[uint16]chan Packet),
	}
}

func (d *demux) Select(pid uint16) <-chan Packet {
	ch := make(chan Packet)
	d.chMutex.Lock()
	d.channels[pid] = ch
	d.chMutex.Unlock()
	return ch
}

func (d *demux) Clear(pid uint16) {
	d.chMutex.Lock()
	if ch, found := d.channels[pid]; found {
		close(ch)
		delete(d.channels, pid)
	}
	d.chMutex.Unlock()
}

func (d *demux) Run() {
	for pkt := range d.inCh {
		if pkt.IsNull() {
			continue
		}

		pid := pkt.PID()
		d.chMutex.Lock()
		if channel, found := d.channels[pid]; found {
			channel <- pkt
		}
		d.chMutex.Unlock()
	}

	d.chMutex.Lock()
	for _, ch := range d.channels {
		close(ch)
	}
	d.chMutex.Unlock()
}
