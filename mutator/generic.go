package mutator

import (
	"log"
	"sync"
)

type genericMutator struct {
	sync.Mutex
	logger      *log.Logger
	name        string
	description string

	started bool
	waited  bool
	done    chan struct{}
}

func newGeneric(logger *log.Logger) genericMutator {
	return genericMutator{
		Mutex:       sync.Mutex{},
		logger:      logger,
		name:        name,
		description: description,
		started:     false,
		waited:      false,
		done:        make(chan struct{}),
	}
}
