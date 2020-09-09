package filestation

import (
	"encoding/base64"
	"fmt"
)

type loginResponse struct {
	Status     FileStationStatus `json:"status,omitempty"`
	Version    string            `json:"version,omitempty"`
	Build      string            `json:"build,omitempty"`
	SessionID  string            `json:"sid,omitempty"`
	AdminGroup int               `json:"admingroup,omitempty"`
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
	case WFM2_SUCCESS: // success
		s.sessionID = result.SessionID
		s.conn.SetQueryParam("sid", s.sessionID)
		return nil
	}

	return result.Status
}

type logoutResponse struct {
	Status  FileStationStatus `json:"status,omitempty"`
	Version string            `json:"version,omitempty"`
	Build   string            `json:"build,omitempty"`
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
	case WFM2_SUCCESS: // success
		s.sessionID = ""
		s.conn.SetQueryParam("sid", "")
		return nil
	}

	return result.Status
}

func encodePassword(pwd string) string {
	return base64.StdEncoding.EncodeToString([]byte(pwd))
}
