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
	"path"
	"strconv"
	"time"
)

func log(msg string, args ...interface{}) {
	const Beat = time.Duration(24) * time.Hour / 1000
	tzBiel := time.FixedZone("Biel", 60*60)
	t := time.Now().In(tzBiel)
	ty := t.Sub(time.Date(t.Year(), 1, 1, 0, 0, 0, 0, tzBiel))
	td := t.Sub(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, tzBiel))
	ts := fmt.Sprintf(
		"%d/%d@%d",
		t.Year()%1000,
		int(ty/(time.Duration(24)*time.Hour))+1,
		int(td/Beat),
	)
	fmt.Fprintf(os.Stderr, ts+" "+msg+"\n", args...)
}

func main() {
	var err error
	var fpath string
	var mode *os.FileMode

	if len(os.Args) > 1 {
		fpath = os.Args[1]
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

	dirpath := path.Dir(fpath)
	dirstat, err := os.Stat(dirpath)
	if err == nil && !dirstat.IsDir() {
		log("%s is not a dir", dirpath)
		os.Exit(1)
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirpath, os.FileMode(0755))
		if err != nil {
			log("%v", err)
			os.Exit(1)
		}
	}

	var f *os.File
	f, err = os.Create(fpath)
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
