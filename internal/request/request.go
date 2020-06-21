/*
Package request provides a simple JSON logger for http.Request objects.
*/
package request

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

type logEntry struct {
	ID             string
	Time           time.Time
	RemoteIP       string
	Method         string
	Host           string
	RequestURI     string
	Referrer       string
	UserAgent      string
	ForwardedFor   string
	ForwardedHost  string
	ForwardedProto string
}

// Log prints details for the given request object as JSON to STDOUT.
func Log(req *http.Request, id string) {
	ip, _, _ := net.SplitHostPort(req.RemoteAddr)
	entry := &logEntry{
		ID:             id,
		Time:           time.Now().UTC(),
		RemoteIP:       ip,
		Method:         req.Method,
		Host:           req.Host,
		RequestURI:     req.URL.RequestURI(),
		Referrer:       req.Header.Get("Referer"),
		UserAgent:      req.Header.Get("User-Agent"),
		ForwardedFor:   req.Header.Get("X-Forwarded-For"),
		ForwardedHost:  req.Header.Get("X-Forwarded-Host"),
		ForwardedProto: req.Header.Get("X-Forwarded-Proto"),
	}
	b, _ := json.Marshal(entry)
	fmt.Println(string(b))
}
