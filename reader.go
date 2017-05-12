package ts

import (
	"bufio"
	"io"

	"github.com/Comcast/gots/packet"
)

func reader(reader *bufio.Reader, outCh chan packet.Packet) {
	_, err := packet.Sync(reader)

	for err == nil {
		pkt := make(packet.Packet, packet.PacketSize)
		if _, err = io.ReadFull(reader, pkt); err != nil {
			continue
		}
		outCh <- pkt
	}
	close(outCh)
}

func Reader(r *bufio.Reader) <-chan packet.Packet {
	outCh := make(chan packet.Packet)
	go reader(r, outCh)
	return outCh
}
