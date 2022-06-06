package openers

import (
	"ccat/log"
	"errors"
	"io"
	"os"
	"sync"
)

var (
	// register() is called from init() so this has to be global
	globalCollection = NewCollection("global")
)

type Opener interface {
	Open(s string, lock bool) (io.ReadCloser, error)
	Evaluate(s string) float32 //score the ability to open
	Name() string
	Description() string
}

type OpenerCollection struct {
	sync.Mutex
	Name    string
	openers []Opener
}

func NewCollection(name string) *OpenerCollection {
	//log.Printf("openers collection %s ready.\n", name)
	return &OpenerCollection{
		Name: name,
	}
}

func register(opener Opener) error {
	globalCollection.Lock()
	globalCollection.openers = append(globalCollection.openers, opener)
	globalCollection.Unlock()
	log.SetDebug(os.Stderr)
	log.Debugf(" opener \"%s\" registered (%s)\n", opener.Name(), opener.Description())
	return nil
}

func Open(s string, lock bool) (io.ReadCloser, error) {
	log.Debugf(" openers: request to open %s\n", s)

	var eMax float32
	var oChosen Opener
	for _, o := range globalCollection.openers {
		e := o.Evaluate(s)
		log.Debugf(" openers: evaluate %s with \"%s\": %v\n", s, o.Name(), e)
		if e > eMax {
			eMax = e
			oChosen = o
		}

	}
	if eMax == 0.0 {
		return nil, errors.New("No adequate opener found.")
	}
	log.Debugf(" openers: chosen one is \"%s\"\n", oChosen.Name())
	return oChosen.Open(s, lock)
}
