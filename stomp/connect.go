package discordgatewayclient_stompmanager

import (
	"net/url"
	"strings"

	"github.com/go-stomp/stomp"
)

// Errors
var (
	errInvalidURIScheme = stompConnectionError{Msg: "invalid connection URI scheme, must be \"stomp\""}
	errFailedToPassURI  = stompConnectionError{Msg: "failed to parse connection URI"}
	errMissingURI       = stompConnectionError{Msg: "missing connection URI, must be set"}
)

type stompConnectionError struct {
	Msg string
}

func (e stompConnectionError) Error() string {
	return e.Msg
}

// CreateStompConnection creates a connection to a STOMP broker from a
// connection URI and returns it.
func CreateStompConnection(uri string) (*stomp.Conn, error) {
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
	if u.Port() == "" {
		u.Host += ":61613"
	}

	// Decide connection options
	opts := []func(*stomp.Conn) error{}
	if u.User != nil {
		user := u.User.Username()
		if pass, pSet := u.User.Password(); pSet {
			opts = append(opts, stomp.ConnOpt.Login(user, pass))
		}
	}
	path := u.EscapedPath()
	if path != "" && path != "/" {
		opts = append(opts, stomp.ConnOpt.Host(strings.Split(path, "/")[0]))
	}
	opts = append(opts, stomp.ConnOpt.AcceptVersion(stomp.V11))
	opts = append(opts, stomp.ConnOpt.AcceptVersion(stomp.V12))

	// Create connection... or not
	return stomp.Dial("tcp", u.Host, opts...)
}
