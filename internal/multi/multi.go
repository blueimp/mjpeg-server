package multi

import (
	"io"
	"sync"
)

// MapWriter implements io.Writer and writes to a map of Writers.
// Unlike io.MultiWriter it ignores errors returned by the individual Writers.
// Writers can be added and removed with the Set and Delete methods.
type MapWriter struct {
	writers map[string]io.Writer
	lock    *sync.RWMutex
}

func (t *MapWriter) Write(p []byte) (int, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	for _, w := range t.writers {
		w.Write(p)
	}
	return len(p), nil
}

// Set adds a Writer with the given ID to the Writers map.
// It returns the new size of the Writers map.
func (t *MapWriter) Set(id string, w io.Writer) int {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.writers[id] = w
	return len(t.writers)
}

// Delete removes a Writer with the given ID from the Writers map.
// It returns the new size of the Writers map.
func (t *MapWriter) Delete(id string) int {
	t.lock.Lock()
	defer t.lock.Unlock()
	delete(t.writers, id)
	return len(t.writers)
}

// Size returns the size of the Writers map.
func (t *MapWriter) Size() int {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return len(t.writers)
}

// NewMapWriter creates a new MapWriter.
func NewMapWriter() *MapWriter {
	writers := make(map[string]io.Writer)
	lock := sync.RWMutex{}
	return &MapWriter{writers, &lock}
}
