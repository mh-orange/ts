package ts

import (
	"bufio"
	"fmt"
	"io"
)

var (
	ErrSyncByteNotFound = fmt.Errorf("Sync Byte (0x47) not found in bit stream")
)

type PacketReader interface {
	Read() (Packet, error)
}

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

type reader struct {
	r    *bufio.Reader
	sync bool
}

func (r *reader) Read() (pkt Packet, err error) {
	if !r.sync {
		err = readSync(r.r)
		if err == nil {
			r.sync = true
		}
	}

	if err == nil {
		pkt = make(Packet, PacketSize)
		_, err = io.ReadFull(r.r, pkt)
	}

	return
}

func NewReader(r io.Reader) PacketReader {
	if _, ok := r.(*bufio.Reader); !ok {
		r = bufio.NewReader(r)
	}

	return &reader{r.(*bufio.Reader), false}
}
