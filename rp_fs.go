package replacer

import (
	"io/ioutil"
	"os"
	"path"
)

func cutStubExt(path string) string {
	n := len(path)

	maxI := n - 1

	if path[maxI] == '.' {
		maxI--
	}

	sf := true
	li := 0
	ri := maxI + 1

	for i := maxI; i > 0 && li == 0; i-- {
		switch path[i] {
		case '.':
			if sf {
				if path[i+1:ri] == "stub" {
					ri = i
					sf = false
				}
			}

			break
		case '/':
			li = i + 1
		}
	}

	return path[li:ri]
}

func getAllTemplates(paths []string) ([]string, error) {
	resSet := make([]string, 0)

	for _, p := range paths {
		n := len(p)

		if n == 0 {
			continue
		}

		f, err := os.Stat(p)

		if err != nil {
			return nil, err
		}

		if f.IsDir() {
			files, err := ioutil.ReadDir(p)

			if err != nil {
				return nil, err
			}

			for _, f := range files {
				if f.IsDir() {
					continue
				}

				resSet = append(resSet, path.Join(p, f.Name()))
			}

			continue
		}

		resSet = append(resSet, p)
	}

	return resSet, nil
}
