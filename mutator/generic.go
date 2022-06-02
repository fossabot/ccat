package mutator

import (
	"ccat/log"
	"sync"
)

type genericMutator struct {
	mu          sync.Mutex
	logger      *log.Logger
	name        string
	description string

	started bool
	waited  bool
	done    chan struct{}
}

func newGeneric(logger *log.Logger) genericMutator {
	return genericMutator{
		mu:          sync.Mutex{},
		logger:      logger,
		name:        name,
		description: description,
		started:     false,
		waited:      false,
		done:        make(chan struct{}),
	}
}
