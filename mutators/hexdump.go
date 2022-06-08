package mutators

import (
	"ccat/log"
	"encoding/hex"
	"fmt"
	"io"
)

var hexName = "hex"
var hexDescription = "dump in Hex"

type hexMutator struct {
	genericMutator
}

type hexFactory struct {
}

func init() {
	f := new(hexFactory)
	register(hexName, f)
}

func (f *hexFactory) New(logger *log.Logger) (Mutator, error) {
	logger.Println("hex: new")
	return &hexMutator{
		genericMutator: newGeneric(logger),
	}, nil
}

func (m *hexMutator) Start(w io.WriteCloser, r io.ReadCloser) error {
	m.mu.Lock()
	if m.started {
		m.mu.Unlock()
		return fmt.Errorf("hex: mutator has already started.")
	}
	m.started = true
	m.mu.Unlock()
	m.logger.Printf("hex: start %v\n", w)

	go func() {
		m.logger.Printf("hex: dumping from %v to %v\n", r, w)
		written, err := hexDump(w, r)
		m.logger.Printf("hex: done\n")
		if err != nil {
			m.logger.Println(err)
		}
		m.logger.Printf("hex: written %d bytes\n", written)
		m.logger.Printf("hex: closing %v\n", w)
		w.Close()
		if err != nil {
			m.logger.Println(err)
		}
		close(m.done)
	}()

	return nil
}
func (m *hexMutator) Wait() error {
	m.logger.Printf("hex: wait called\n")
	m.mu.Lock()
	if !m.started {
		m.mu.Unlock()
		return fmt.Errorf("hex: mutator is not started")
	}
	if m.waited {
		m.mu.Unlock()
		return fmt.Errorf("hex: mutator is already waited")
	}
	m.waited = true
	m.mu.Unlock()
	<-m.done
	return nil
}

func (m *hexMutator) Name() string {
	return hexName
}
func (m *hexMutator) Description() string {
	return hexDescription
}
func (f *hexFactory) Name() string {
	return hexName
}
func (f *hexFactory) Description() string {
	return hexDescription
}

func hexDump(w io.WriteCloser, r io.ReadCloser) (int64, error) {

	dumper := hex.Dumper(w)
	defer dumper.Close()

	return io.Copy(dumper, r)
}
