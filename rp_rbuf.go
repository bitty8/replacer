package replacer

import (
	"golang.org/x/exp/mmap"
)

type rbuf struct {
	mm     *mmap.ReaderAt
	n      int
	offset int
}

func newRBuf(mm *mmap.ReaderAt) *rbuf {
	return &rbuf{
		mm:     mm,
		n:      mm.Len(),
		offset: 0,
	}
}

func (rb *rbuf) canAdvance(value int) bool {
	return rb.offset+value < rb.n
}

func (rb *rbuf) advanceAt(value int) {
	rb.offset += value
}

func (rb *rbuf) canRecvByte() bool {
	return rb.canAdvance(0)
}

func (rb *rbuf) recvByte() byte {
	if !rb.canAdvance(0) {
		return 0
	}

	b := rb.mm.At(rb.offset)
	rb.advanceAt(1)

	return b
}

func (rb *rbuf) readAt(offset int, n int) []byte {
	mem := make([]byte, n)

	n, err := rb.mm.ReadAt(mem, int64(offset))

	if err != nil || n == 0 {
		return nil
	}

	return mem
}

func (rb *rbuf) setOffset(offset int) {
	rb.offset = offset
}
