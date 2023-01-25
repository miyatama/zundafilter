package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"syscall"
	"unsafe"
	"zundafilter/filters"
	"zundafilter/log"
	"zundafilter/zunda_mecab"

	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	flag.Parse()
}

func main() {
        logger := log.GetLogger()
	defer logger.Sync()
	sugar := logger.Sugar()

	text, err := readFile()
	if err != nil {
		sugar.Errorf("%v" , err)
		os.Exit(1)
	}
	zundaDbRepository := zunda_mecab.ZundaDbRepository{}
	mecabWrapper := zunda_mecab.MecabWrapper{
		Logger:       log.GetLogger(),
	}
	filter := filters.ZundaFilter{
		ZundaDb:      zundaDbRepository,
		MecabWrapper: &mecabWrapper,
		Logger:       log.GetLogger(),
	}
	convertedText, err := filter.Convert(text)
	if err != nil {
		sugar.Errorf("ZundaFilter error: %v" , err)
		os.Exit(1)
	}
	fmt.Print(convertedText)
}

func readFile() (string, error) {
	var filename string
	if args := flag.Args(); len(args) > 0 {
		filename = args[0]
	}
	var r io.Reader
	switch filename {
	case "":
		if terminal.IsTerminal(int(syscall.Stdin)) {
			return "", errors.New("usage: capitalize path")
		}
		r = os.Stdin
	case "-":
		r = os.Stdin
	default:
		f, err := os.Open(filename)
		if err != nil {
			return "", err
		}
		defer f.Close()
		r = f
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return *(*string)(unsafe.Pointer(&b)), nil
}
