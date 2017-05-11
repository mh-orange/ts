package ts

import (
	"bufio"
	"io"
	"sync"

	"github.com/Comcast/gots/packet"
)

type Demux interface {
	Select(pid uint16) <-chan packet.Packet
	Clear(pid uint16)
	Run() error
}

type demux struct {
	reader   *bufio.Reader
	channels map[uint16]chan packet.Packet
	chMutex  sync.Mutex
}

func NewDemux(reader *bufio.Reader) Demux {
	return &demux{
		reader:   reader,
		channels: make(map[uint16]chan packet.Packet),
	}
}

func (d *demux) Select(pid uint16) <-chan packet.Packet {
	d.chMutex.Lock()
	d.channels[pid] = make(chan packet.Packet)
	d.chMutex.Unlock()
	return d.channels[pid]
}

func (d *demux) Clear(pid uint16) {
	d.chMutex.Lock()
	if ch, found := d.channels[pid]; found {
		close(ch)
		delete(d.channels, pid)
	}
	d.chMutex.Unlock()
}

func (d *demux) Run() (err error) {
	_, err = packet.Sync(d.reader)
	if err != nil {
		return err
	}

	for {
		pkt := make(packet.Packet, packet.PacketSize)
		if _, err = io.ReadFull(d.reader, pkt); err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				err = nil
			}
			break
		}

		if ok, _ := packet.IsNull(pkt); ok {
			continue
		}

		pid, _ := packet.Pid(pkt)
		d.chMutex.Lock()
		if channel, ok := d.channels[pid]; ok {
			channel <- pkt
		}
		d.chMutex.Unlock()
	}

	d.chMutex.Lock()
	for _, ch := range d.channels {
		close(ch)
	}
	d.chMutex.Unlock()
	return
}
