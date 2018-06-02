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
func (s *FileStationSession) Login(username, password string) error {
	var result *loginResponse

	err := s.getForEntity(&result, "cgi-bin/filemanager/wfm2Login.cgi", QueryParameters{
		"user": username,
		"pwd":  encodePassword(password),
	})
	if err != nil {
		return err
	}

	switch result.Status {
	case 0: // password wrong
		return fmt.Errorf("Password or Username is invalid")
	case 1: // success
		s.SessionID = result.SessionID
		return nil
	case 8: // disabled
		return fmt.Errorf("API is disabled")
	}

	return fmt.Errorf("Unknown status code: %v", result.Status)
}

type logoutResponse struct {
	Status  int    `json:"status,omitempty"`
	Version string `json:"version,omitempty"`
	Build   string `json:"build,omitempty"`
}

// Logout invalidates the session.
func (s *FileStationSession) Logout() error {
	// no logged-in?
	if s.SessionID == "" {
		return nil
	}

	var result *logoutResponse

	err := s.getForEntity(&result, "cgi-bin/filemanager/wfm2Logout.cgi", QueryParameters{
		"sid": s.SessionID,
	})
	if err != nil {
		return err
	}

	switch result.Status {
	case 1: // success
		s.SessionID = ""
		return nil
	case 8: // disabled
		return fmt.Errorf("API is disabled")
	}

	return fmt.Errorf("Unknown status code: %v", result.Status)
}

func encodePassword(pwd string) string {
	return base64.StdEncoding.EncodeToString([]byte(pwd))
}
