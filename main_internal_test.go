package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blueimp/mjpeg-server/internal/multi"
	"github.com/blueimp/mjpeg-server/internal/recording"
)

func TestRequestHandler(t *testing.T) {
	clients = multi.NewMapWriter()
	size := clients.Size()
	if size != 0 {
		t.Errorf("Unexpected clients size: %d. Expected: %d", size, 0)
	}
	started := 0
	stopped := 0
	startRecording = func(command string, args []string, w io.Writer) (
		stop context.CancelFunc,
		wait recording.WaitFunc,
	) {
		size := clients.Size()
		if size != 1 {
			t.Errorf("Unexpected clients size: %d. Expected: %d", size, 1)
		}
		started++
		stop = func() {
			size := clients.Size()
			if size != 0 {
				t.Errorf("Unexpected clients size: %d. Expected: %d", size, 0)
			}
			stopped++
		}
		wait = func() error { return nil }
		return
	}
	rec := httptest.NewRecorder()
	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(
		"GET",
		"http://localhost:9000/",
		nil,
	).WithContext(ctx)
	go func() {
		rec := httptest.NewRecorder()
		ctx, cancel2 := context.WithCancel(context.Background())
		req := httptest.NewRequest(
			"GET",
			"http://localhost:9000/",
			nil,
		).WithContext(ctx)
		go func() {
			time.Sleep(100 * time.Millisecond)
			size := clients.Size()
			if size != 2 {
				t.Errorf("Unexpected clients size: %d. Expected: %d", size, 2)
			}
			cancel2()
		}()
		requestHandler(rec, req)
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()
	requestHandler(rec, req)
	if started != 1 {
		t.Errorf(
			"Unexpected number of recording starts: %d. Expected: %d",
			started,
			1,
		)
	}
	if stopped != 1 {
		t.Errorf(
			"Unexpected number of recording stops: %d. Expected: %d",
			stopped,
			1,
		)
	}
	if rec.Code != http.StatusOK {
		t.Errorf(
			"Unexpected response status: %d. Expected: %d",
			rec.Code,
			http.StatusOK,
		)
	}
	header := rec.Header().Get("Content-Type")
	expectedHeader := "multipart/x-mixed-replace;boundary=ffmpeg"
	if header != expectedHeader {
		t.Errorf(
			"Unexpected Content-Type header: %s. Expected: %s",
			header,
			expectedHeader,
		)
	}
	header = rec.Header().Get("Cache-Control")
	expectedHeader = "no-store"
	if header != expectedHeader {
		t.Errorf(
			"Unexpected Cache-Control header: %s. Expected: %s",
			header,
			expectedHeader,
		)
	}
	header = rec.Header().Get("Connection")
	expectedHeader = "close"
	if header != expectedHeader {
		t.Errorf(
			"Unexpected Connection header: %s. Expected: %s",
			header,
			expectedHeader,
		)
	}
}

func TestRequestHandlerWithInvalidMethod(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		"POST",
		"http://localhost:9000/",
		nil,
	)
	requestHandler(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf(
			"Unexpected response status: %d. Expected: %d",
			rec.Code,
			http.StatusMethodNotAllowed,
		)
	}
}

func TestRequestHandlerWithInvalidPath(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		"GET",
		"http://localhost:9000/invalid/",
		nil,
	)
	requestHandler(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf(
			"Unexpected response status: %d. Expected: %d",
			rec.Code,
			http.StatusNotFound,
		)
	}
}

func TestRequestHandlerWithCustomPath(t *testing.T) {
	*urlPath = "/banana"
	rec := httptest.NewRecorder()
	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(
		"GET",
		"http://localhost:9000/banana",
		nil,
	).WithContext(ctx)
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()
	requestHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf(
			"Unexpected response status: %d. Expected: %d",
			rec.Code,
			http.StatusOK,
		)
	}
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(
		"GET",
		"http://localhost:9000/",
		nil,
	)
	requestHandler(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf(
			"Unexpected response status: %d. Expected: %d",
			rec.Code,
			http.StatusNotFound,
		)
	}
	*urlPath = "/"
}

func TestRequestHandlerWithCustomBoundary(t *testing.T) {
	*boundary = "banana"
	rec := httptest.NewRecorder()
	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(
		"GET",
		"http://localhost:9000/",
		nil,
	).WithContext(ctx)
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()
	requestHandler(rec, req)
	header := rec.Header().Get("Content-Type")
	expectedHeader := "multipart/x-mixed-replace;boundary=banana"
	if header != expectedHeader {
		t.Errorf(
			"Unexpected Content-Type header: %s. Expected: %s",
			header,
			expectedHeader,
		)
	}
	*boundary = "ffmpeg"
}
