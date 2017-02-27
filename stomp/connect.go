package stomp

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gmallard/stompngo"
)

const (
	// AcceptedVersions are accepted STOMP protocol versions.
	AcceptedVersions = "1.1,1.2"

	// ReadTimeout refers to the read timeout value.
	ReadTimeout = time.Minute / time.Millisecond
	// WriteTimeout refers to the write timeout value.
	WriteTimeout = time.Minute / time.Millisecond
)

// Errors.
var (
	errInvalidURIScheme           = stompConnectionURIError{Msg: "invalid connection URI scheme, must be \"stomp\""}
	errFailedToPassURI            = stompConnectionURIError{Msg: "failed to parse connection URI"}
	errMissingURI                 = stompConnectionURIError{Msg: "missing connection URI, must be set"}
	errURIUsernameWithoutPassword = stompConnectionURIError{Msg: "username supplied on URI without password"}
)

var defaultConnectHeaders = stompngo.Headers{
	stompngo.HK_HEART_BEAT, fmt.Sprintf("%d,%d", WriteTimeout, ReadTimeout),
	stompngo.HK_ACCEPT_VERSION, AcceptedVersions,
}

type stompConnectionURIError struct {
	Msg string
}

func (e stompConnectionURIError) Error() string {
	return e.Msg
}

// CreateStompConnection creates a connection to a STOMP broker from a
// connection URI and returns it.
func CreateStompConnection(uri string) (*stompngo.Connection, error) {
	// Parse connection URI
	if uri == "" {
		return nil, errMissingURI
	}
	u, err := url.Parse(uri)
	if err != nil {
		return nil, errFailedToPassURI
	}
	if u.Scheme != "stomp" {
		return nil, errInvalidURIScheme
	}
	u.Host += ":61613"
	if u.User != nil {
		if _, pSet := u.User.Password(); !pSet {
			return nil, errURIUsernameWithoutPassword
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
	log.WithFields(log.Fields{
		"headers": fmt.Sprintf("%#v", h),
	}).Debug("constructed STOMP connection-specific headers")

	// Create network connection
	n, err := net.Dial(stompngo.NetProtoTCP, u.Host)
	if err != nil {
		return nil, err
	}
	log.Debug("created net.Conn for STOMP client connection")

	// Connect to the STOMP broker
	return stompngo.Connect(n, h.AddHeaders(defaultConnectHeaders))
}
