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
	for _, w := range t.writers {
		w.Write(p)
	}
	t.lock.RUnlock()
	return len(p), nil
}

// Set adds a Writer with the given ID to the Writers map.
// It returns the new size of the Writers map.
func (t *MapWriter) Set(id string, w io.Writer) (size int) {
	t.lock.Lock()
	t.writers[id] = w
	size = len(t.writers)
	t.lock.Unlock()
	return
}

// Delete removes a Writer with the given ID from the Writers map.
// It returns the new size of the Writers map.
func (t *MapWriter) Delete(id string) (size int) {
	t.lock.Lock()
	delete(t.writers, id)
	size = len(t.writers)
	t.lock.Unlock()
	return
}

// Size returns the size of the Writers map.
func (t *MapWriter) Size() int {
	t.lock.RLock()
	size := len(t.writers)
	t.lock.RUnlock()
	return size
}

// NewMapWriter creates a new MapWriter.
func NewMapWriter() *MapWriter {
	writers := make(map[string]io.Writer)
	lock := sync.RWMutex{}
	return &MapWriter{writers, &lock}
}
