package highlighter

import (
	"ccat/log"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func IsTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

func getTerminalSize() (width, height int, err error) {
	if IsTerminal() {
		return term.GetSize(int(os.Stdout.Fd()))
	}
	// fallback when piping to a file!
	return 80, 24, nil // VT100 terminal size
}

func isITerm2() bool {
	// LC_TERMINAL = iTerm2
	// TERM_PROGRAM = iTerm.app
	if os.Getenv("LC_TERMINAL") == "iTerm2" {
		return IsTerminal()
	}
	return false
}

func clearScreen() {
	if IsTerminal() {
		fmt.Print("\033[H\033[2J")
	}
}

func supportedColors() uint {
	var colors uint
	if !IsTerminal() {
		log.Debugln("stdout is not a terminal, detecting colors anyway...")
	}

	switch {
	case isITerm2():
		colors = 16_000_000
		log.Debugln("  supportedColors: iterm2 -> 16M colors detected")
	case strings.ToLower(os.Getenv("COLORTERM")) == "truecolor":
		colors = 16_000_000
		log.Debugln("  supportedColors: truecolor -> 16M colors detected")
	case os.Getenv("TERM") == "xterm-256color":
		colors = 256
		log.Debugln("  supportedColors: xterm-256color -> 256 colors detected")
	default:
		log.Debugf("  supportedColors: unkown term, $TERM==%s -> 8 colors detected\n", os.Getenv("TERM"))
		return 8
	}

	return colors
}
