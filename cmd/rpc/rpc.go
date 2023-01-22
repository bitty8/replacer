package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/bitty8/replacer"
)

const (
	inSep = ','
)

func parseInput(input []byte) []string {
	arr := make([]string, 0)

	n := len(input)

	i := 0
	li := i

	for ; i < n; i++ {
		if input[i] == inSep {
			p := string(input[li:i])

			if len(p) > 0 {
				arr = append(arr)
			}

			li = i + 1
		}
	}

	if i > li {
		arr = append(arr, string(input[li:i]))
	}

	return arr
}

func main() {
	var (
		input      string
		outDir     string
		force      bool
		paramsfile string
		params     string
	)

	flag.StringVar(&input, "input", "", "the list of input stub files separated by ,")
	flag.StringVar(&outDir, "outdir", "", "build directory")
	flag.BoolVar(&force, "force", false, "replace with an empty value if the key is not found either abort process")
	flag.StringVar(&paramsfile, "mapfile", "", "file than contains a map of values (only json)")
	flag.StringVar(&params, "params", "", "list of params in key=vale view separated by , (lame=x,fname=y)")

	flag.Parse()

	if len(input) == 0 {
		fmt.Fprintln(os.Stderr, "the input flag is required")
		os.Exit(1)
	}

	if len(outDir) == 0 {
		outDir, ok := os.LookupEnv("PWD")

		if !ok {
			outDir, ok = os.LookupEnv("HOME")

			if !ok {
				fmt.Fprintln(os.Stderr, "cannot detect build dir, the outdir flag is not defined and vars such as PWD and HOME is not found")
				os.Exit(2)
			}

			outDir += "/.replacer"
		}

		outDir += "/" + strconv.Itoa(int(time.Now().Unix())) + "_build"
		fmt.Println(outDir)
	}

	rpl, err := replacer.NewReplacer(
		parseInput([]byte(input)),
		outDir,
		force,
		paramsfile,
		params,
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}

	if !rpl.Exec() {
		os.Exit(4)
	}
}
