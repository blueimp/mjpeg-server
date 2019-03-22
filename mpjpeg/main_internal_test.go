package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
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

func TestStreamFiles(t *testing.T) {
	filePaths := []string{"../gopher.jpg"}
	imageData, _ := ioutil.ReadFile(filePaths[0])
	expectedOutput := bytes.Join(
		[][]byte{
			[]byte("--ffmpeg"),
			[]byte("Content-Type: image/jpeg"),
			[]byte(""),
			imageData,
			[]byte("--ffmpeg--"),
			[]byte(""),
		},
		[]byte("\r\n"),
	)
	*noLoop = true
	stdout, stderr := outputHelper(func() {
		streamFiles(filePaths)
	})
	if len(stderr) != 0 {
		t.Errorf("Unexpected stderr: %s", stderr)
	}
	if !bytes.Equal(stdout, expectedOutput) {
		t.Error("Unexpected stdout")
	}
}
