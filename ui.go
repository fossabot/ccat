//go:build crappy
// +build crappy

package main

import (
	"ccat/globalctx"
	"ccat/log"
	"flag"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var textView *tview.TextView
var argUi = flag.Bool("ui", false, "display file in a minimal ui")

func ui(c chan struct{}) {
	if !*argUi {
		log.Debugln("not starting ui")
		globalctx.Set("textview", os.Stdout)
		c <- struct{}{}
		close(c)
		return
	}
	log.Debugln("starting ui")
	app := tview.NewApplication()
	textView = tview.NewTextView().SetDynamicColors(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	globalctx.Set("textview", tview.ANSIWriter(textView))
	c <- struct{}{}

	textView.SetDoneFunc(func(key tcell.Key) {
		app.Stop()
		close(c)
	})
	textView.SetBorder(true)
	textView.ScrollToBeginning()

	if err := app.SetRoot(textView, true).SetFocus(textView).Run(); err != nil {
		panic(err)
	}

}
