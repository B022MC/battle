package plaza

import (
	"sync"
)

type RecvBuf struct {
	buffer  [1024 * 1024]byte
	size    int
	handler chan []byte
	mut     sync.Mutex
}

func NewRecvBuf(handler chan []byte) *RecvBuf {
	rf := &RecvBuf{}
	rf.handler = handler
	return rf
}

func (that *RecvBuf) Add(data []byte) {
	that.mut.Lock()
	defer that.mut.Unlock()
	defer func() {
		_ = recover()
	}()

	for _, d := range data {
		that.buffer[that.size] = d
		that.size++
	}

	for {
		if sz := that.getPacketSize(); sz <= that.size {
			packet := make([]byte, sz)
			for i := 0; i < sz; i++ {
				packet[i] = that.buffer[i]
			}
			for i := sz; i < that.size; i++ {
				that.buffer[i-sz] = that.buffer[i]
			}
			that.size -= sz

			that.handler <- packet
		} else {
			break
		}
	}

}

func (that *RecvBuf) getPacketSize() int {
	lo := uint16(that.buffer[2])
	hi := uint16(that.buffer[3])
	return int(hi<<8 + lo)
}
