package multi

import (
	"io"
	"sync"
)

type empty struct{}

type mapWriter struct {
	writers map[io.Writer]empty
	lock    *sync.RWMutex
}

// MapWriter is an interface to write to a map of Writers.
// Writers can be added and removed with the Add and Remove methods, while the
// Size method returns the current map size.
type MapWriter interface {
	Write(p []byte) (int, error)
	Add(w io.Writer) (size int)
	Remove(w io.Writer) (size int)
	Size() int
}

// Write implements io.Writer but ignores errors by the individual Writers.
func (t *mapWriter) Write(p []byte) (int, error) {
	t.lock.RLock()
	for w := range t.writers {
		w.Write(p)
	}
	t.lock.RUnlock()
	return len(p), nil
}

// Add puts the given Writer into the Writers map.
// It returns the new size of the Writers map.
func (t *mapWriter) Add(w io.Writer) (size int) {
	t.lock.Lock()
	t.writers[w] = empty{}
	size = len(t.writers)
	t.lock.Unlock()
	return
}

// Remove deletes the given Writer from the Writers map.
// It returns the new size of the Writers map.
func (t *mapWriter) Remove(w io.Writer) (size int) {
	t.lock.Lock()
	delete(t.writers, w)
	size = len(t.writers)
	t.lock.Unlock()
	return
}

// Size returns the size of the Writers map.
func (t *mapWriter) Size() int {
	t.lock.RLock()
	size := len(t.writers)
	t.lock.RUnlock()
	return size
}

// NewMapWriter creates a new MapWriter.
func NewMapWriter() MapWriter {
	writers := make(map[io.Writer]empty)
	return &mapWriter{writers, &sync.RWMutex{}}
}
