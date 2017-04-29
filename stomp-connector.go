package main

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gmallard/stompngo"
)

// Accepted STOMP protocol versions.
const AcceptedVersions = "1.1,1.2"

// Read and write timeout values.
const (
	ReadTimeout  = time.Minute / time.Millisecond
	WriteTimeout = time.Minute / time.Millisecond
)

// defaultConnectHeaders is a set of default connection headers for any connections to STOMP brokers.
var defaultConnectHeaders = stompngo.Headers{
	stompngo.HK_HEART_BEAT, fmt.Sprintf("%d,%d", WriteTimeout, ReadTimeout),
	stompngo.HK_ACCEPT_VERSION, AcceptedVersions,
}

// createStompConnection creates a connection to a STOMP broker from a connection URI and returns it.
func createStompConnection(uri string) (*stompngo.Connection, error) {
	// Parse connection URI
	if uri == "" {
		return nil, errors.New("missing connection URI, must be set")
	}
	u, err := url.Parse(uri)
	if err != nil {
		return nil, errors.New("failed to parse connection URI")
	}
	if u.Scheme != "stomp" {
		return nil, errors.New("invalid connection URI scheme, must be \"stomp\"")
	}
	if u.Port() == "" {
		u.Host += ":61613"
	}
	if u.User != nil {
		if _, pSet := u.User.Password(); !pSet {
			return nil, errors.New("username supplied in URI without password")
		}
	}
	log.Debug("validated STOMP connection URI")

	// Decide connection-specific headers
	h := stompngo.Headers{}
	if path := u.EscapedPath(); path != "" && path != "/" {
		h = h.Add(stompngo.HK_HOST, strings.Split(path, "/")[0])
	} else {
		h = h.Add(stompngo.HK_HOST, "")
	}
	if u.User != nil {
		h = h.Add(stompngo.HK_LOGIN, u.User.Username())
		pass, _ := u.User.Password()
		h = h.Add(stompngo.HK_PASSCODE, pass)
	}
	log.WithField("headers", fmt.Sprintf("%#v", h)).Debug("constructed STOMP connection-specific headers")

	// Create network connection
	n, err := net.Dial(stompngo.NetProtoTCP, u.Host)
	if err != nil {
		return nil, err
	}
	log.Debug("created net.Conn for STOMP client connection")

	// Connect to the STOMP broker
	return stompngo.Connect(n, h.AddHeaders(defaultConnectHeaders))
}
