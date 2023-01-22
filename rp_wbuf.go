package replacer

import "os"

const (
	wbufSize = 1024
)

type wbuf struct {
	f          *os.File
	data       []byte
	offset     int
	fileOffset int
}

func newWBuf(f *os.File) *wbuf {
	return &wbuf{
		f:          f,
		data:       make([]byte, wbufSize, wbufSize),
		offset:     0,
		fileOffset: 0,
	}
}

func (wb *wbuf) writeByte(b byte) error {
	if wb.offset >= wbufSize {
		err := wb.flush()

		if err != nil {
			return err
		}
	}

	wb.data[wb.offset] = b
	wb.offset++

	return nil
}

func (wb *wbuf) writeSlice(arr []byte) error {
	for _, b := range arr {
		err := wb.writeByte(b)

		if err != nil {
			return err
		}
	}

	return nil
}

func (wb *wbuf) flush() error {
	if wb.offset > 0 {
		n, err := wb.f.WriteAt(wb.data[0:wb.offset], int64(wb.fileOffset))

		if err != nil {
			return err
		}

		wb.fileOffset += n
		wb.offset = 0
	}

	return nil
}
