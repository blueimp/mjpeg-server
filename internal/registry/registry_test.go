package registry

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/blueimp/mjpeg-server/internal/recording"
)

var started int
var stopped int

func startRecordingHelper(command string, args []string, w io.Writer) (
	stop context.CancelFunc,
	wait recording.WaitFunc,
) {
	started++
	stop = func() {
		stopped++
	}
	wait = func() error { return nil }
	return
}

func outputHelper(fn func()) (stdout []byte, stderr []byte) {
	outReader, outWriter, _ := os.Pipe()
	errReader, errWriter, _ := os.Pipe()
	originalOut := os.Stdout
	originalErr := os.Stderr
	os.Stdout = outWriter
	os.Stderr = errWriter
	fn()
	outWriter.Close()
	errWriter.Close()
	stdout, _ = ioutil.ReadAll(outReader)
	stderr, _ = ioutil.ReadAll(errReader)
	os.Stdout = originalOut
	os.Stderr = originalErr
	return
}

func TestNew(t *testing.T) {
	reg := New("go", []string{"version"}, false)
	if reg == nil {
		t.Error("Unexpected: nil")
	}
	_, ok := interface{}(reg).(Registry)
	if !ok {
		t.Error("Unexpected: not a MapWriter")
	}
}

func TestGenerateID(t *testing.T) {
	reg := New("go", []string{"version"}, false)
	id := reg.GenerateID()
	if id != "1" {
		t.Errorf("Unexpected generated ID: %s. Expected: %s", id, "1")
	}
	id = reg.GenerateID()
	if id != "2" {
		t.Errorf("Unexpected generated ID: %s. Expected: %s", id, "2")
	}
	id = reg.GenerateID()
	if id != "3" {
		t.Errorf("Unexpected generated ID: %s. Expected: %s", id, "3")
	}
}

func TestAdd(t *testing.T) {
	started = 0
	stopped = 0
	startRecording = startRecordingHelper
	reg := New("go", []string{"version"}, false)
	if started != 0 {
		t.Errorf("Unexpected started recordings: %d. Expected: %d", started, 0)
	}
	timeBefore := time.Now()
	stdout, stderr := outputHelper(func() {
		reg.Add("1", &bytes.Buffer{})
	})
	timeAfter := time.Now()
	if started != 1 {
		t.Errorf("Unexpected started recordings: %d. Expected: %d", started, 1)
	}
	if string(stderr) != "" {
		t.Errorf("Unexpected stderr: %s", stderr)
	}
	var entry logEntry
	json.Unmarshal(stdout, &entry)
	if entry.ID != "1" {
		t.Errorf("Unexpected 'ID' log: %s. Expected: %s", entry.ID, "1")
	}
	if entry.Time.Before(timeBefore) {
		t.Errorf("Unexpected 'Time' log: %s", entry.Time)
	}
	if entry.Time.After(timeAfter) {
		t.Errorf("Unexpected 'Time' log: %s", entry.Time)
	}
	if entry.Registered != true {
		t.Errorf(
			"Unexpected 'Registered' log: %t. Expected: %t",
			entry.Registered,
			true,
		)
	}
	if entry.NumClients != 1 {
		t.Errorf(
			"Unexpected 'NumClients' log: %d. Expected: %d",
			entry.NumClients,
			1,
		)
	}
	stdout, stderr = outputHelper(func() {
		reg.Add("2", &bytes.Buffer{})
	})
	if started != 1 {
		t.Errorf("Unexpected started recordings: %d. Expected: %d", started, 1)
	}
	if string(stderr) != "" {
		t.Errorf("Unexpected stderr: %s", stderr)
	}
	entry = logEntry{}
	json.Unmarshal(stdout, &entry)
	if entry.ID != "2" {
		t.Errorf("Unexpected 'ID' log: %s. Expected: %s", entry.ID, "2")
	}
	if entry.Registered != true {
		t.Errorf(
			"Unexpected 'Registered' log: %t. Expected: %t",
			entry.Registered,
			true,
		)
	}
	if entry.NumClients != 2 {
		t.Errorf(
			"Unexpected 'NumClients' log: %d. Expected: %d",
			entry.NumClients,
			2,
		)
	}
}

func TestRemove(t *testing.T) {
	started = 0
	stopped = 0
	startRecording = startRecordingHelper
	reg := New("go", []string{"version"}, false)
	var (
		buffer1 bytes.Buffer
		buffer2 bytes.Buffer
	)
	outputHelper(func() {
		reg.Add("1", &buffer1)
		reg.Add("2", &buffer2)
	})
	if stopped != 0 {
		t.Errorf("Unexpected stopped recordings: %d. Expected: %d", stopped, 0)
	}
	timeBefore := time.Now()
	stdout, stderr := outputHelper(func() {
		reg.Remove("2", &buffer2)
	})
	timeAfter := time.Now()
	if stopped != 0 {
		t.Errorf("Unexpected stopped recordings: %d. Expected: %d", stopped, 0)
	}
	if string(stderr) != "" {
		t.Errorf("Unexpected stderr: %s", stderr)
	}
	var entry logEntry
	json.Unmarshal(stdout, &entry)
	if entry.ID != "2" {
		t.Errorf("Unexpected 'ID' log: %s. Expected: %s", entry.ID, "2")
	}
	if entry.Time.Before(timeBefore) {
		t.Errorf("Unexpected 'Time' log: %s", entry.Time)
	}
	if entry.Time.After(timeAfter) {
		t.Errorf("Unexpected 'Time' log: %s", entry.Time)
	}
	if entry.Registered != false {
		t.Errorf(
			"Unexpected 'Registered' log: %t. Expected: %t",
			entry.Registered,
			false,
		)
	}
	if entry.NumClients != 1 {
		t.Errorf(
			"Unexpected 'NumClients' log: %d. Expected: %d",
			entry.NumClients,
			1,
		)
	}
	stdout, stderr = outputHelper(func() {
		reg.Remove("1", &buffer1)
	})
	if stopped != 1 {
		t.Errorf("Unexpected stopped recordings: %d. Expected: %d", stopped, 1)
	}
	if string(stderr) != "" {
		t.Errorf("Unexpected stderr: %s", stderr)
	}
	entry = logEntry{}
	json.Unmarshal(stdout, &entry)
	if entry.ID != "1" {
		t.Errorf("Unexpected 'ID' log: %s. Expected: %s", entry.ID, "1")
	}
	if entry.Registered != false {
		t.Errorf(
			"Unexpected 'Registered' log: %t. Expected: %t",
			entry.Registered,
			false,
		)
	}
	if entry.NumClients != 0 {
		t.Errorf(
			"Unexpected 'NumClients' log: %d. Expected: %d",
			entry.NumClients,
			0,
		)
	}
}

func TestNewWithDirectStart(t *testing.T) {
	started = 0
	stopped = 0
	startRecording = startRecordingHelper
	reg := New("go", []string{"version"}, true)
	if started != 1 {
		t.Errorf("Unexpected started recordings: %d. Expected: %d", started, 1)
	}
	var buffer1 bytes.Buffer
	outputHelper(func() {
		reg.Add("1", &buffer1)
	})
	if started != 1 {
		t.Errorf("Unexpected started recordings: %d. Expected: %d", started, 1)
	}
	outputHelper(func() {
		reg.Remove("1", &buffer1)
	})
	if stopped != 0 {
		t.Errorf("Unexpected stopped recordings: %d. Expected: %d", stopped, 0)
	}
}
