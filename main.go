package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/blueimp/mjpeg-server/internal/registry"
	"github.com/blueimp/mjpeg-server/internal/request"
)

var (
	// Version provides the program version information.
	// It is provided at build time via -ldflags="-X main.Version=VERSION".
	Version     = "dev"
	showVersion = flag.Bool("v", false, "Output version and exit")
	directStart = flag.Bool("d", false, "Start command directly")
	addr        = flag.String("a", ":9000", "TCP listen address")
	urlPath     = flag.String("p", "/", "URL path")
	boundary    = flag.String("b", "ffmpeg", "Multipart boundary")
	command     string
	args        []string
	reg         registry.Registry
)

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
	id := reg.GenerateID()
	request.Log(req, id)
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
	reg.Add(id, res)
	// Wait until the client connection is closed.
	<-req.Context().Done()
	reg.Remove(id, res)
}

func parseArgs() {
	flag.Parse()
	command = flag.Arg(0)
	if command != "" {
		args = flag.Args()[1:]
	}
}

func main() {
	log.SetOutput(os.Stderr)
	parseArgs()
	if *showVersion {
		fmt.Println(Version)
		os.Exit(0)
	}
	reg = registry.New(command, args, *directStart)
	log.Fatalln(http.ListenAndServe(*addr, http.HandlerFunc(requestHandler)))
}
