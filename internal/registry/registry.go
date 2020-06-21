/*
Package registry manages the handling of recording clients.
*/
package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/blueimp/mjpeg-server/internal/multi"
	"github.com/blueimp/mjpeg-server/internal/recording"
)

var startRecording = recording.Start

type logEntry struct {
	ID         string
	Time       time.Time
	Registered bool
	NumClients int
}

type registry struct {
	command       string
	args          []string
	clients       multi.MapWriter
	counter       uint64
	stopRecording context.CancelFunc
	waitForStop   recording.WaitFunc
}

// Registry is an interface to manage the handling of recording clients.
// Clients can be added and removed with the Add and Remove methods, while the
// GenerateID method returns an auto-incrementing ID.
type Registry interface {
	GenerateID() string
	Add(id string, w io.Writer) (num int)
	Remove(id string, w io.Writer) (num int)
}

func log(id string, registered bool, numClients int) {
	entry := &logEntry{
		ID:         id,
		Time:       time.Now().UTC(),
		Registered: registered,
		NumClients: numClients,
	}
	b, _ := json.Marshal(entry)
	fmt.Println(string(b))
}

// GenerateID returns an auto-incrementing ID.
func (t *registry) GenerateID() string {
	return strconv.FormatUint(atomic.AddUint64(&t.counter, 1), 10)
}

// Add puts the given Writer into the Registry.
// It returns the new number of clients in the Registry.
func (t *registry) Add(id string, w io.Writer) (num int) {
	num = t.clients.Add(w)
	if num == 1 {
		// First client added, start the recording.
		t.stopRecording, t.waitForStop = startRecording(
			t.command,
			t.args,
			t.clients,
		)
	}
	log(id, true, num)
	return
}

// Add deletes the given Writer from the Registry.
// It returns the new number of clients in the Registry.
func (t *registry) Remove(id string, w io.Writer) (num int) {
	num = t.clients.Remove(w)
	if num == 0 {
		// Last client removed, stop the recording.
		t.stopRecording()
	}
	log(id, false, num)
	return
}

// New creates a new Registry.
func New(command string, args []string) Registry {
	return &registry{
		command,
		args,
		multi.NewMapWriter(),
		0,
		nil,
		nil,
	}
}
