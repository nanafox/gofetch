package gofetch

import (
	"testing"
	"time"
)

// TestDefaultNewClient verifies that the configuration expected is used when
// a new client is created.
func TestDefaultNewClient(t *testing.T) {
	client := New()

	got := client.Config.Debug

	if client.Config.Debug != false {
		t.Fatalf("expected Debug to be %v, but it is %v", false, got)
	}

	expectedTimeout := 500 * time.Millisecond
	currentTimeout := client.Config.Timeout

	if client.Config.Timeout != expectedTimeout {
		t.Fatalf("expected Timeout to be %v, but it is %v", expectedTimeout, currentTimeout)
	}
}

// TestUserConfigUsedForNewClient ensures that the user-provided configs are used
// instead of the defaults.
func TestUserConfigUsedForNewClient(t *testing.T) {
	expectedTimeout := 200 * time.Millisecond
	client := New(Config{Timeout: expectedTimeout, Debug: true})

	got := client.Config.Debug

	if client.Config.Debug != true {
		t.Fatalf("expected Debug to be %v, but it is %v", true, got)
	}

	currentTimeout := client.Config.Timeout

	if client.Config.Timeout != expectedTimeout {
		t.Fatalf("expected Timeout to be %v, but it is %v", expectedTimeout, currentTimeout)
	}
}
