package ts

import (
	"sync"
)

type Demux interface {
	Select(pid uint16, handler PacketHandler)
	Clear(pid uint16)
	Run()
}

type demux struct {
	reader   PacketReader
	handlers map[uint16]PacketHandler
	chMutex  sync.Mutex
}

func NewDemux(reader PacketReader) Demux {
	return &demux{
		reader:   reader,
		handlers: make(map[uint16]PacketHandler),
	}
}

func (d *demux) Select(pid uint16, handler PacketHandler) {
	d.chMutex.Lock()
	d.handlers[pid] = handler
	d.chMutex.Unlock()
}

func (d *demux) Clear(pid uint16) {
	d.chMutex.Lock()
	delete(d.handlers, pid)
	d.chMutex.Unlock()
}

func (d *demux) Run() {
	for pkt, err := d.reader.Read(); err == nil; pkt, err = d.reader.Read() {
		if pkt.IsNull() {
			continue
		}

		pid := pkt.PID()
		d.chMutex.Lock()
		if handler, found := d.handlers[pid]; found {
			handler.Handle(pkt)
		}
		d.chMutex.Unlock()
	}
}
