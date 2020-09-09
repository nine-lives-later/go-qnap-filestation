package filestation

import (
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strings"
	"time"
)

var defaultConfigOptions = ConfigOptions{
	APICallTimeout:              60 * time.Second,
	IgnoreInvalidSSLCertificate: true,
}

// ConfigOptions contains some advanced settings on server communication.
type ConfigOptions struct {
	APICallTimeout              time.Duration
	IgnoreInvalidSSLCertificate bool
}

// FileStationSession is a container for our session state.
type FileStationSession struct {
	host      string
	sessionID string
	conn      *resty.Client
	options   *ConfigOptions
}

// String returns the session's hostname.
func (s *FileStationSession) String() string {
	return s.host
}

// Connect sets up our connection to the QNAP system.
func Connect(host, username, password string, configOptions *ConfigOptions) (*FileStationSession, error) {
	if !strings.HasPrefix(host, "http") {
		host = fmt.Sprintf("https://%s", host)
	}
	if configOptions == nil {
		configOptions = &defaultConfigOptions
	}

	// create the session
	session := &FileStationSession{
		host:    host,
		conn:    resty.New().SetHostURL(host).SetTimeout(configOptions.APICallTimeout),
		options: configOptions,
	}

	// setup SSL certificate handling
	if configOptions.IgnoreInvalidSSLCertificate {
		session.conn.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	// perform login
	err := session.Login(username, password)
	if err != nil {
		return nil, err
	}

	// done
	return session, nil
}

func (s *FileStationSession) Close() error {
	return s.Logout()
}
