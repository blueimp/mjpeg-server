package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/textproto"
	"os"
	"time"
)

var (
	boundary = flag.String("b", "ffmpeg", "Multipart boundary")
	interval = flag.Duration("s", 100*time.Millisecond, "Sleep interval")
	noLoop   = flag.Bool("n", false, "Do not loop forever")
)

func streamFiles(filePaths []string) {
	imageContents := make([][]byte, len(filePaths))
	for i, filePath := range filePaths {
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		imageContents[i] = content
	}
	multipartWriter := multipart.NewWriter(os.Stdout)
	multipartWriter.SetBoundary(*boundary)
	header := make(textproto.MIMEHeader)
	header.Add("Content-Type", "image/jpeg")
	for {
		for i := range filePaths {
			writer, _ := multipartWriter.CreatePart(header)
			writer.Write(imageContents[i])
			time.Sleep(*interval)
		}
		if *noLoop == true {
			break
		}
	}
	multipartWriter.Close()
}

func main() {
	flag.Parse()
	streamFiles(flag.Args())
}
