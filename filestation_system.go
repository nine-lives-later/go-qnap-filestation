package filestation

import (
	"encoding/base64"
	"fmt"
)

type loginResponse struct {
	Status     int    `json:"status,omitempty"`
	Version    string `json:"version,omitempty"`
	Build      string `json:"build,omitempty"`
	SessionID  string `json:"sid,omitempty"`
	AdminGroup int    `json:"admingroup,omitempty"`
}

// Login perform the authentication against the QNAP storage.
// Any existing session will be logged-out, first.
func (s *FileStationSession) Login(username, password string) error {
	// make sure to close any existing sessions
	s.Logout()

	// perform login
	var result loginResponse

	res, err := s.conn.NewRequest().
		ExpectContentType("application/json").
		SetQueryParam("user", username).
		SetQueryParam("pwd", encodePassword(password)).
		SetResult(&result).
		Get("cgi-bin/filemanager/wfm2Login.cgi")
	if err != nil {
		return fmt.Errorf("failed to perform request: %v", err)
	}
	if res.StatusCode() != 200 {
		return fmt.Errorf("failed to perform request: unexpected HTTP status code: %v", res.StatusCode())
	}

	switch result.Status {
	case 0: // password wrong
		return fmt.Errorf("password or username is invalid")
	case 1: // success
		s.sessionID = result.SessionID
		s.conn.SetQueryParam("sid", s.sessionID)
		return nil
	case 2: // exists
		return fmt.Errorf("already exists")
	case 8: // disabled
		return fmt.Errorf("API is disabled")
	}

	return fmt.Errorf("unknown status code: %v", result.Status)
}

type logoutResponse struct {
	Status  int    `json:"status,omitempty"`
	Version string `json:"version,omitempty"`
	Build   string `json:"build,omitempty"`
}

// Logout invalidates the session.
func (s *FileStationSession) Logout() error {
	// no logged-in?
	if s.sessionID == "" {
		return nil
	}

	var result logoutResponse

	res, err := s.conn.NewRequest().
		ExpectContentType("application/json").
		SetResult(&result).
		Get("cgi-bin/filemanager/wfm2Logout.cgi")
	if err != nil {
		return fmt.Errorf("failed to perform request: %v", err)
	}
	if res.StatusCode() != 200 {
		return fmt.Errorf("failed to perform request: unexpected HTTP status code: %v", res.StatusCode())
	}

	switch result.Status {
	case 1: // success
		s.sessionID = ""
		s.conn.SetQueryParam("sid", "")
		return nil
	case 8: // disabled
		return fmt.Errorf("API is disabled")
	}

	return fmt.Errorf("unknown status code: %v", result.Status)
}

func encodePassword(pwd string) string {
	return base64.StdEncoding.EncodeToString([]byte(pwd))
}
