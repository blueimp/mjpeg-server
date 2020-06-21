package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blueimp/mjpeg-server/internal/registry"
)

func TestRequestHandler(t *testing.T) {
	reg = registry.New(command, args)
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
	reg = registry.New(command, args)
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
	reg = registry.New(command, args)
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
	reg = registry.New(command, args)
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
	reg = registry.New(command, args)
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
