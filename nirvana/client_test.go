package nirvana

import (
	"os"
	"testing"
)

func TestNewNirvanaClient(t *testing.T) {
	_, err := NewNirvanaClient(nil)
	if err != nil {
		t.Error(err)
	}
}

func TestAuthenticateError(t *testing.T) {
	c, _ := NewNirvanaClient(nil)
	defer c.Close()

	err := c.Authenticate("xxxxxx@gmail.com", "_")
	if err == nil {
		t.Error("expected err not found")
	}
}

func TestAuthenticate(t *testing.T) {
	c, _ := NewNirvanaClient(nil)
	defer c.Close()

	err := c.Authenticate(os.Getenv("NIRVANA_USERNAME"), os.Getenv("NIRVANA_PASSWORD"))
	if err != nil {
		t.Errorf("unexpected err not found : %v", err)
	}
	if len(c.authToken) == 0 {
		t.Error("auth token error")
	}
}

func TestRetrieveSince(t *testing.T) {
	c, _ := NewNirvanaClient(nil)
	defer c.Close()

	err := c.Authenticate(os.Getenv("NIRVANA_USERNAME"), os.Getenv("NIRVANA_PASSWORD"))
	_, err = c.RetrieveSince(0)
	if err != nil {
		t.Errorf("unexpected err not found : %v", err)
	}
}
