package ts

import (
	"bufio"
	"fmt"
	"io"
)

var (
	ErrSyncByteNotFound = fmt.Errorf("Sync Byte (0x47) not found in bit stream")
)

// Credit: https://github.com/Comcast/gots/blob/master/packet/io.go
func readSync(reader *bufio.Reader) error {
	data := make([]byte, 1)
	for i := int64(0); ; i++ {
		_, err := io.ReadFull(reader, data)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			return err
		}
		if int(data[0]) == 0x47 {
			// check next 188th byte
			rp := bufio.NewReaderSize(reader, PacketSize) // extends only if needed
			nextData, err := rp.Peek(PacketSize)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			if err != nil {
				return err
			}
			if nextData[187] == 0x47 {
				reader.UnreadByte()
				return nil
			}
		}
	}
	return ErrSyncByteNotFound
}

func reader(reader *bufio.Reader, outCh chan Packet) {
	err := readSync(reader)

	for err == nil {
		pkt := make(Packet, PacketSize)
		if _, err = io.ReadFull(reader, pkt); err != nil {
			continue
		}
		outCh <- pkt
	}
	close(outCh)
}

func Reader(r io.Reader) <-chan Packet {
	outCh := make(chan Packet)
	if _, ok := r.(*bufio.Reader); !ok {
		r = bufio.NewReader(r)
	}

	go reader(r.(*bufio.Reader), outCh)
	return outCh
}
