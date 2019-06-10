package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	all := flag.Bool("all", false, "make run scripts for all folders")
	source := flag.String("source", "sources", "relative path to folder containing source code folders")
	args := flag.Args()

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n%s [-options] language_1 [language_N...]\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	var targets []string
	wd, err := os.Getwd()
	if err != nil {
		panic(err) // something crazy happened.
	}
	top, err := os.Open(filepath.Join(wd,*source))
	if err != nil {
		panic(err) // FIXME
	}
	if *all == true {
		targets, err = top.Readdirnames(0)
		if err != nil {
			panic(err)
		}
	} else {
		targets = args
	}
	for _,t := range targets {
		t = slugUp(t)
		l := fetchLang(t)
		l.Path = filepath.Join(top.Name(),t)
		if l.Compiled == true {
			err = mkBuild(l)
			if err != nil {
				panic(err)
			}
		}
		err = mkRunScript(l)
		if err != nil {
			panic(err)
		}
	}

}
