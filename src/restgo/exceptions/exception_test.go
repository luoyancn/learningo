package exceptions

import (
	"testing"
)

func TestNewNotFoundException(t *testing.T) {
	fake_exception := NewNotFoundException("fake_resource_type", "fake_id")
	if fake_exception.Error() !=
		`fake_resource_type with fake_id cannot be found` {
		t.Error("Not expected error messages")
	}
}
