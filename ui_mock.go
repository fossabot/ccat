//go:build !crappy
// +build !crappy

package main

import (
	"ccat/globalctx"
	"ccat/log"
	"os"
)

type mocked struct{}

var textView *mocked
var argUi = new(bool)

func ui(c chan struct{}) {
	log.Debugln("NOT starting ui")
	globalctx.Set("textview", os.Stdout)

	log.Debugln("tell channel we're ready")
	c <- struct{}{}
	log.Debugln("close channel")
	close(c)
}

func (m mocked) ScrollToBeginning() {
}
