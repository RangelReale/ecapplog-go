package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/RangelReale/ecapplog-go"
)

var (
	app         = flag.String("app", "eclcat", "ecapplog app name")
	category    = flag.String("category", "ALL", "ecapplog category name")
	split       = flag.String("split", "\\n", "split string (default=\\n")
	lineNumbers = flag.Bool("n", false, "show line numbers")
)

func main() {
	flag.Parse()

	client := ecapplog.NewClient(ecapplog.WithAppName(*app),
		ecapplog.WithFlushOnClose(true))
	client.Open()
	defer client.Close()

	if flag.NArg() == 0 {
		err := output(client, os.Stdin)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	} else {
		for _, fname := range flag.Args() {
			fh, err := os.Open(fname)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			defer fh.Close()

			err = output(client, fh)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}
	}
}

func output(client *ecapplog.Client, r io.Reader) error {
	scanner := bufio.NewScanner(r)
	switch *split {
	case "", "\\n":
		scanner.Split(bufio.ScanLines)
	case " ":
		scanner.Split(bufio.ScanWords)
	default:
		scanner.Split(SplitAt(*split))
	}
	lineno := 1
	for scanner.Scan() {
		text := scanner.Text()
		var opt []ecapplog.LogOption
		if *lineNumbers {
			opt = append(opt, ecapplog.WithSource(text))
			text = fmt.Sprintf("%04d %s", lineno, text)
		}
		client.Log(time.Now(), ecapplog.Priority_INFORMATION, *category, text, opt...)
		lineno++
	}
	return scanner.Err()
}

func SplitAt(substring string) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	searchBytes := []byte(substring)
	searchLen := len(searchBytes)
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		dataLen := len(data)

		// Return nothing if at end of file and no data passed
		if atEOF && dataLen == 0 {
			return 0, nil, nil
		}

		// Find next separator and return token
		if i := bytes.Index(data, searchBytes); i >= 0 {
			return i + searchLen, data[0:i], nil
		}

		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return dataLen, data, nil
		}

		// Request more data.
		return 0, nil, nil
	}
}
