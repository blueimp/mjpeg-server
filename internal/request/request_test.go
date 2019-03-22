package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func outputHelper(fn func()) (stdout []byte, stderr []byte) {
	outReader, outWriter, _ := os.Pipe()
	errReader, errWriter, _ := os.Pipe()
	originalOut := os.Stdout
	originalErr := os.Stderr
	defer func() {
		os.Stdout = originalOut
		os.Stderr = originalErr
	}()
	os.Stdout = outWriter
	os.Stderr = errWriter
	fn()
	outWriter.Close()
	errWriter.Close()
	stdout, _ = ioutil.ReadAll(outReader)
	stderr, _ = ioutil.ReadAll(errReader)
	return
}

func TestLog(t *testing.T) {
	req := httptest.NewRequest(
		"GET",
		"http://localhost:9000/mjpeg",
		nil,
	)
	req.Header.Set("Referer", "http://example.org/")
	req.Header.Set("User-Agent", "Examplebot/1.0 (+http://example.org)")
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("X-Forwarded-Host", "example")
	req.Header.Set("X-Forwarded-For", "127.0.0.1")
	timeBefore := time.Now()
	stdout, stderr := outputHelper(func() {
		Log(req)
	})
	timeAfter := time.Now()
	if string(stderr) != "" {
		t.Errorf("Unexpected stderr: %s", stderr)
	}
	var entry logEntry
	json.Unmarshal(stdout, &entry)
	if entry.Time.Before(timeBefore) {
		t.Errorf("Unexpected 'Time' log: %s", entry.Time)
	}
	if entry.Time.After(timeAfter) {
		t.Errorf("Unexpected 'Time' log: %s", entry.Time)
	}
	// httptest.NewRequest always uses the RemoteAddr 192.0.2.1:1234
	if entry.RemoteIP != "192.0.2.1" {
		t.Errorf(
			"Unexpected 'IP' log: %s. Expected: %s",
			entry.RemoteIP,
			"192.0.2.1",
		)
	}
	if entry.Method != "GET" {
		t.Errorf("Unexpected 'Method' log: %s. Expected: %s", entry.Method, "GET")
	}
	if entry.Host != "localhost:9000" {
		t.Errorf(
			"Unexpected 'Host' log: %s. Expected: %s",
			entry.Host,
			"localhost:9000",
		)
	}
	if entry.RequestURI != "/mjpeg" {
		t.Errorf(
			"Unexpected 'RequestURI' log: %s. Expected: %s",
			entry.RequestURI,
			"/mjpeg",
		)
	}
	if entry.Referrer != "http://example.org/" {
		t.Errorf(
			"Unexpected 'Referrer' log: %s. Expected: %s",
			entry.Referrer,
			"http://example.org/",
		)
	}
	if entry.UserAgent != "Examplebot/1.0 (+http://example.org)" {
		t.Errorf(
			"Unexpected 'UserAgent' log: %s. Expected: %s",
			entry.UserAgent,
			"Examplebot/1.0 (+http://example.org)",
		)
	}
	if entry.ForwardedFor != "127.0.0.1" {
		t.Errorf(
			"Unexpected 'ForwardedFor' log: %s. Expected: %s",
			entry.ForwardedFor,
			"127.0.0.1",
		)
	}
	if entry.ForwardedHost != "example" {
		t.Errorf(
			"Unexpected 'ForwardedHost' log: %s. Expected: %s",
			entry.ForwardedHost,
			"example",
		)
	}
	if entry.ForwardedProto != "https" {
		t.Errorf(
			"Unexpected 'ForwardedProto' log: %s. Expected: %s",
			entry.ForwardedProto,
			"https",
		)
	}
}
