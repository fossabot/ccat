package main

import (
	"bufio"
	"ccat/highlighter"
	"ccat/lockable"
	"ccat/log"
	"ccat/pipedcmd"
	"ccat/scanners"
	"fmt"
	"io"
	"os"
	"regexp"
)

func processFile(path string) {
	log.Debugf("processing %s...\n", path)

	var from io.ReadCloser = os.Stdin
	var err error
	if path != "-" {
		from, err = lockable.FileOpen(path, *argLockIn)
		if err != nil {
			log.Println(err)
			return
		}
		defer lockable.FileClose(from.(*os.File), *argLockIn)
	}
	/*************************************/
	if len(*argExec) > 0 {
		log.Debugf("creating pipedcmd %v...\n", *argExec)
		cmd, err := pipedcmd.New(*argExec)
		log.Debugf("%s", log.Pp(cmd))

		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			log.Debugf("waiting pipedcmd %v...\n", *argExec)

			if err := cmd.Wait(); err != nil {
				log.Println(err)
			}
		}()

		log.Debugf("start pipedcmd %s\n", cmd)

		err = cmd.Start(from)
		if err != nil {
			log.Println(err)
		}
		from = cmd.Stdout.(*os.File)
	}

	if *argHighlight {
		log.Debugln("highlighting...")
		r, w := io.Pipe()
		highlighter.Go(w, from, highlighter.Options{
			FileName:      path,
			StyleHint:     *argStyle,
			FormatterHint: *argFormatter,
			LexerHint:     *argLexer,
		})
		from = r
	}
	log.Debugln("initializing Scanner...")

	scanner := bufio.NewScanner(from)

	splitFn := scanners.ScanBytes
	if len(tokens) > 0 {
		splitFn = scanners.ScanLines
	}
	if *argSplitByWords {
		splitFn = scanners.ScanWords
	}
	scanner.Split(splitFn)
	lineNumber := 1
	log.Debugln("start Scanning...")
	for scanner.Scan() {
		var matched bool
		text := scanner.Text()
		for _, token := range tokens {
			var err error
			var regexpPrefix string
			if *argInsensitive {
				regexpPrefix = "(?i)"
			}

			//fmt.Println("text ", text)
			//fmt.Println("token ", token)
			matched, err = regexp.MatchString(regexpPrefix+token, text)
			if err != nil {
				log.Println(err)
			}
			if matched {
				color := tmap[token]
				text = color.Sprint(text)
				break
			}
		}
		if *argLineNumber {
			fmt.Printf("%d ", lineNumber)
			lineNumber++
		}
		if !*argOnlyMatching || matched && *argOnlyMatching {
			fmt.Print(text)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
	log.Debugln("end Scanning...")
}
