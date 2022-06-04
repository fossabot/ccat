package highlighter

import (
	"ccat/log"
	"io"
	"math/rand"
	"time"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

const (
	DEFAULT_STYLE     = "github"
	DEFAULT_FORMATTER = "terminal256"
)

type Chroma struct {
	style     string
	formatter string
	lexer     string
}

func (h *Chroma) HighLight(w io.WriteCloser, r io.ReadCloser, o Options) error {
	log.Debugln(" highlighter: start chroma Highlighter")
	log.Debugln(log.Pp(o))

	var filename string = o.FileName

	someSourceCode, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	log.Debugf(" highlighter: read %v bytes\n", len(someSourceCode))

	log.Debugf(" highlighter: registered lexers are: %v\n", lexers.Names(true))
	lexersList := lexers.Names(true)
	var lexer chroma.Lexer
	if len(o.LexerHint) > 0 && stringInSlice(o.LexerHint, lexersList) {
		log.Debugf(" highlighter: setting the lexer to %v\n", o.LexerHint)
		lexer = lexers.Get(o.LexerHint)
	} else {
		lexer = lexers.Match(filename)
		if lexer == nil {
			log.Debugf(" highlighter: filename did not help to find a lexer, analyzing content...\n")
			lexer = lexers.Analyse(string(someSourceCode))
			if lexer == nil {
				log.Debugf(" highlighter: fallbacking the lexer\n")
				lexer = lexers.Fallback
			}
		}
		lexer = chroma.Coalesce(lexer)

	}
	//log.Debugf(" highlighter: lexers %v\n", lexer.Config().Name)

	log.Debugf(" highlighter: chosen Lexer is %v\n", lexer.Config().Name)

	log.Debugf(" highlighter: registered styles are: %v\n", styles.Names())
	//registered styles are: [abap algol algol_nu arduino autumn base16-snazzy borland bw colorful doom-one doom-one2 dracula emacs friendly fruity github hr_high_contrast hrdark igor lovelace manni monokai monokailight murphy native nord onesenterprise paraiso-dark paraiso-light pastie perldoc pygments rainbow_dash rrt solarized-dark solarized-dark256 solarized-light swapoff tango trac vim vs vulcan witchhazel xcode xcode-dark]

	stylesList := styles.Names()
	if o.StyleHint == "random" {
		rand.Seed(time.Now().UnixNano())
		randStyle := rand.Intn(len(stylesList))
		h.style = stylesList[randStyle]
	} else if len(o.StyleHint) > 0 && stringInSlice(o.StyleHint, stylesList) {
		h.style = o.StyleHint
	} else {
		h.style = DEFAULT_STYLE
	}

	style := styles.Get(h.style) // or monokai
	if style == nil {
		style = styles.Fallback
	}
	log.Debugf(" highlighter: style is %+v\n", style.Name)

	log.Debugf(" highlighter: registered formatters are: %v\n", formatters.Names())

	formattersList := formatters.Names()
	if len(o.FormatterHint) > 0 && stringInSlice(o.FormatterHint, formattersList) {
		h.formatter = o.FormatterHint
	} else {
		h.formatter = DEFAULT_FORMATTER
	}
	formatter := formatters.Get(h.formatter)
	if formatter == nil {
		formatter = formatters.Fallback
	}
	log.Debugf(" highlighter: formatter is %v\n", h.formatter)

	iterator, err := lexer.Tokenise(nil, string(someSourceCode))
	if err != nil {
		return err
	}

	err = formatter.Format(w, style, iterator)
	if err != nil {
		return err
	}

	log.Debugln(" highlighter: end chroma Highlight")
	return nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			log.Debugf("%v == %v", a, b)
			return true
		}
	}
	return false
}
