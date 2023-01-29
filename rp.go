package replacer

import (
	"fmt"
	"os"
	"path"
	"sync"
	"sync/atomic"

	"golang.org/x/exp/mmap"
)

type Replacer struct {
	templatesPaths []string
	outDir         string
	outName        string
	isForce        bool
	gm             *gmap
	wg             *sync.WaitGroup
	hasErrFlag     int32
}

func NewReplacer(inPaths []string, outDir string, force bool, paramsFile string, params string, outName string) (*Replacer, error) {
	if len(inPaths) > 1 && len(outName) > 0 {
		return nil, fmt.Errorf("the flag \"outname\" is applied if the number of input files is equal to one")
	}

	inFilesPaths, err := getAllTemplates(inPaths)

	if err != nil {
		return nil, err
	}

	gm := newGMap()

	if len(paramsFile) > 0 {
		err = parseParamsFile(gm, paramsFile)

		if err != nil {
			return nil, err
		}
	}

	parseParamsString(gm, params)

	return &Replacer{
		templatesPaths: inFilesPaths,
		outDir:         outDir,
		outName:        outName,
		isForce:        force,
		gm:             gm,
		wg:             &sync.WaitGroup{},
		hasErrFlag:     0,
	}, nil
}

func (r *Replacer) setErrFlag() {
	atomic.StoreInt32(&r.hasErrFlag, 1)
}

func (r *Replacer) replace(wb *wbuf, rb *rbuf) error {
	var (
		b     byte
		kpath []byte
		keyLi int

		//escape flag
		ef  bool
		err error
	)

	ef = false

	for rb.canRecvByte() {
		b = rb.recvByte()

		if ef {
			if err = wb.writeByte(b); err != nil {
				goto write_err
			}
			ef = false
			continue
		}

		switch b {
		case '\\':
			ef = true
			break
		case '{':
			keyLi = rb.offset

			for rb.canRecvByte() {
				b = rb.recvByte()

				if b == '}' || b <= 32 || b == '{' {
					break
				}
			}

			if b != '}' {
				ef = true
				rb.setOffset(keyLi - 1)
			} else {
				kpath = rb.readAt(keyLi, rb.offset-keyLi-1)

				if kpath == nil {
					return fmt.Errorf("empty key at offset %d in stub file", keyLi)
				}

				val := r.gm.get(string(kpath))

				if val != nil {
					wb.writeSlice(val)
				} else if !r.isForce {
					return fmt.Errorf("cant replace at position %d, key {%s} is not found", keyLi-1, kpath)
				}
			}

			break
		default:
			if err = wb.writeByte(b); err != nil {
				goto write_err
			}
		}
	}

	wb.flush()

	return nil

write_err:
	return fmt.Errorf("write to out file error, %s", err.Error())
}

func (r *Replacer) runTask(inFile string) {
	var (
		outF *os.File
		wb   *wbuf
		rb   *rbuf
		mm   *mmap.ReaderAt
		err  error
	)

	defer r.wg.Done()

	outPath := ""

	if len(r.outName) > 0 {
		outPath = path.Join(r.outDir, r.outName)
	} else {
		outPath = path.Join(r.outDir, cutStubExt(inFile))
	}
	// make out file path
	os.Remove(outPath)

	// create out file path
	if outF, err = os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE, 0666); err != nil {
		r.setErrFlag()
		fmt.Fprintf(os.Stderr, "open out file %s error %s\n", outPath, err.Error())
		return
	}

	defer outF.Close()
	// create new write buffer for writing data to out file
	wb = newWBuf(outF)

	// mmap on stub filee
	if mm, err = mmap.Open(inFile); err != nil {
		r.setErrFlag()
		fmt.Fprintf(os.Stderr, "open mmap on file %s error, %s\n", inFile, err.Error())
		return
	}
	defer mm.Close()

	rb = newRBuf(mm)

	if err = r.replace(wb, rb); err != nil {
		r.setErrFlag()
		outF.Close()
		os.Remove(outPath)
		fmt.Fprintf(os.Stderr, "error: %s, file: %s\n", err, inFile)
		return
	}

	fmt.Fprintf(os.Stdout, "file %s replaced successfuly to %s\n", inFile, outPath)
}

func (r *Replacer) Exec() bool {
	err := os.MkdirAll(r.outDir, os.ModePerm)

	if err != nil {
		fmt.Fprintf(os.Stderr, "mkdir %s error, %s", r.outDir, err.Error())
		return false
	}

	for _, f := range r.templatesPaths {
		r.wg.Add(1)
		go r.runTask(f)
	}

	r.wg.Wait()

	if atomic.LoadInt32(&r.hasErrFlag) > 0 {
		return false
	}

	return true
}
