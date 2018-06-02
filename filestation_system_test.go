package filestation

import (
	"os"
	"strings"
	"testing"
)

func createTestSession(t *testing.T) *FileStationSession {
	return createTestSessionEx(t, "", "")
}

func createTestSessionEx(t *testing.T, username, password string) *FileStationSession {
	// retrieve api auth if undefined
	if username == "" {
		username = os.Getenv("QNAP_USER")

		if username == "" {
			username = "unitest-user"
		}
	}

	if password == "" {
		password = os.Getenv("QNAP_PWD")

		if password == "" {
			password = "t3st!!!"
		}
	}

	// retrieve hostname
	host := os.Getenv("QNAP_HOSTNAME")

	if host == "" {
		host = "storage1.konsorten.net:8443"
	}

	// create the session
	session, err := Connect(host, username, password, nil)

	if err != nil {
		t.Fatalf("Failed to connect to QNAP File Station API: %v", err)
	}

	return session
}

func TestPasswordEncode(t *testing.T) {
	input := "admin"
	expected := "YWRtaW4="

	if encodePassword(input) != expected {
		t.Fatal("password encoding failed")
	}
}

func TestInvalidHost(t *testing.T) {
	_, err := Connect("d0esn0tex1st", "inva1dus3r", "inval1dAp1K3y", nil)
	if err == nil {
		t.Fatal("Error expected")
	}

	if !strings.Contains(err.Error(), "no such host") {
		t.Fatalf("Wrong error message returned: %v", err)
	}
}

func TestConnect(t *testing.T) {
	s := createTestSession(t)

	// real logout
	err := s.Logout()
	if err != nil {
		t.Fatalf("Failed to logout: %v", err)
	}

	// redundant logout
	err = s.Logout()
	if err != nil {
		t.Fatalf("Failed to logout: %v", err)
	}
}

func TestConnect_InvalidLogin(t *testing.T) {
	host := os.Getenv("QNAP_HOSTNAME")

	if host == "" {
		host = "storage1.konsorten.net:8443"
	}

	_, err := Connect(host, "unkn0wnUs3r", "!nval1dP@ssw0rd", nil)
	if err == nil {
		t.Fatal("Error expected")
	}

	if !strings.Contains(err.Error(), "Password or Username is invalid") {
		t.Fatalf("Wrong error message returned: %v", err)
	}
}
