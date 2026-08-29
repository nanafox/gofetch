package gofetch

import (
	"testing"
	"time"
)

// TestDefaultNewClient verifies that the configuration expected is used when
// a new client is created.
func TestDefaultNewClient(t *testing.T) {
	c := New(Config{})

	got := c.Config.Debug

	if c.Config.Debug != false {
		t.Fatalf("expected Debug to be %v, but it is %v", false, got)
	}

	expTimeout := 500 * time.Millisecond
	curTimeout := c.Config.Timeout

	if c.Config.Timeout != expTimeout {
		t.Fatalf("expected Timeout to be %v, but it is %v", expTimeout, curTimeout)
	}
}

// TestUserConfigUsedForNewClient ensures that the user-provided configs are used
// instead of the defaults.
func TestUserConfigUsedForNewClient(t *testing.T) {
	expTimeout := 200 * time.Millisecond
	c := New(Config{Timeout: expTimeout, Debug: true})

	got := c.Config.Debug

	if c.Config.Debug != true {
		t.Fatalf("expected Debug to be %v, but it is %v", true, got)
	}

	curTimeout := c.Config.Timeout

	if c.Config.Timeout != expTimeout {
		t.Fatalf("expected Timeout to be %v, but it is %v", expTimeout, curTimeout)
	}
}
