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
func readSync(reader *bufio.Reader) (err error) {
	data := make([]byte, 1)
	for i := int64(0); ; i++ {
		_, err = io.ReadFull(reader, data)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err = ErrSyncByteNotFound
		}
		if err != nil {
			break
		}
		if int(data[0]) == 0x47 {
			// check next 188th byte
			rp := bufio.NewReaderSize(reader, PacketSize) // extends only if needed
			var nextData []byte
			nextData, err = rp.Peek(PacketSize)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				err = ErrSyncByteNotFound
			}
			if err != nil {
				break
			}
			if nextData[187] == 0x47 {
				reader.UnreadByte()
				break
			}
		}
	}
	return
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
