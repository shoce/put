/*
history:
2015-04-28 v1
2015-05-29 v2

usage:
put put-file-test 600 <some-another-file
id | put id.out.text
sudo id | sudo put sudo.id.out.text

GoFmt GoBuildNull GoRelease GoBuild
*/

package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

func log(msg string, args ...interface{}) {
	ts := time.Now().Local().Format("Jan/02;15:04")
	fmt.Fprintf(os.Stderr, ts+" "+msg+"\n", args...)
}

func main() {
	var err error
	var path string
	var mode *os.FileMode

	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		log("usage: put path [mode]")
		os.Exit(1)
	}

	if len(os.Args) > 2 {
		var m uint64
		m, err = strconv.ParseUint(os.Args[2], 8, 32)
		if err != nil {
			log("invalid file mode `%s`", os.Args[2])
			os.Exit(1)
		}
		var m2 os.FileMode
		m2 = os.FileMode(m)
		mode = &m2
	}

	var f *os.File
	f, err = os.Create(path)
	if err != nil {
		log("%v", err)
		os.Exit(1)
	}
	defer f.Close()

	if mode != nil {
		err = f.Chmod(*mode)
		if err != nil {
			log("%v", err)
			os.Exit(1)
		}
	}

	_, err = io.Copy(f, os.Stdin)
	if err != nil {
		log("%v", err)
		os.Exit(1)
	}
}
