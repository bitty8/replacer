package replacer

import (
	"io/ioutil"
	"os"

	"github.com/buger/jsonparser"
)

const (
	assignChar = '='
	sepChar    = ','
)

func parseParamsString(gm *gmap, params string) {
	var (
		bv  []byte
		i   int
		n   int
		li  int
		key string
	)

	bv = []byte(params)
	n = len(bv)

	i = 0
	li = i

	key = ""

	for ; i < n; i++ {
		switch bv[i] {
		case assignChar:
			key = string(bv[li:i])
			li = i + 1
			break
		case sepChar:
			if len(key) > 0 {
				gm.set(key, bv[li:i])
				key = ""
			}

			li = i + 1
		}
	}

	if i > li && len(key) > 0 {
		gm.set(key, bv[li:i])
	}
}

func parseParamsFile(gm *gmap, filePath string) error {
	f, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)

	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(f)

	if err != nil {
		return err
	}

	return jsonparser.ObjectEach(data, func(key, value []byte, dataType jsonparser.ValueType, offset int) error {
		if dataType == jsonparser.String {
			gm.set(string(key), value)
		}

		return nil
	})
}
