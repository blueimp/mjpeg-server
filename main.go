package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/blueimp/mjpeg-server/internal/multi"
	"github.com/blueimp/mjpeg-server/internal/recording"
	"github.com/blueimp/mjpeg-server/internal/request"
)

var (
	// Version provides the program version information.
	// It is provided at build time via -ldflags="-X main.Version=VERSION".
	Version        = "dev"
	showVersion    = flag.Bool("v", false, "Output version and exit")
	addr           = flag.String("a", ":9000", "TCP listen address")
	urlPath        = flag.String("p", "/", "URL path")
	boundary       = flag.String("b", "ffmpeg", "Multipart boundary")
	command        string
	args           []string
	clients        multi.MapWriter
	startRecording recording.StartFunc
	stopRecording  context.CancelFunc
)

func parseArgs() {
	flag.Parse()
	command = flag.Arg(0)
	if command != "" {
		args = flag.Args()[1:]
	}
}

func registerClient(w io.Writer) {
	if clients.Add(w) == 1 {
		// First client added, start the recording.
		stopRecording, _ = startRecording(command, args, clients)
	}
}

func deregisterClient(w io.Writer) {
	if clients.Remove(w) == 0 {
		// Last client removed, stop the recording.
		stopRecording()
	}
}

func setHeaders(header http.Header) {
	// Provide the multipart boundary via MJPEG over HTTP content-type header.
	// See also:
	// - https://en.wikipedia.org/wiki/Motion_JPEG#M-JPEG_over_HTTP
	// - https://tools.ietf.org/html/rfc2046#section-5.1.1
	header.Set(
		"Content-Type",
		fmt.Sprintf("multipart/x-mixed-replace;boundary=%s", *boundary),
	)
	// Prevent client caches from storing the response.
	// See also: https://tools.ietf.org/html/rfc7234#section-5.2.1.5
	header.Set("Cache-Control", "no-store")
	// Signal to the client that the connection will be closed after completion of
	// the response.
	// See also: https://tools.ietf.org/html/rfc2616#section-14.10
	header.Set("Connection", "close")
}

func requestHandler(res http.ResponseWriter, req *http.Request) {
	request.Log(req)
	if req.Method != "GET" {
		res.Header().Set("Allow", "GET")
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if req.URL.Path != *urlPath {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	setHeaders(res.Header())
	registerClient(res)
	// Wait until the client connection is closed.
	<-req.Context().Done()
	deregisterClient(res)
}

func main() {
	log.SetOutput(os.Stderr)
	parseArgs()
	if *showVersion {
		fmt.Println(Version)
		os.Exit(0)
	}
	clients = multi.NewMapWriter()
	startRecording = recording.Start
	log.Fatalln(http.ListenAndServe(*addr, http.HandlerFunc(requestHandler)))
}
