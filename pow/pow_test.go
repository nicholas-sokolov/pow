package pow

import (
	"testing"
)

func TestHashcashData_String(t *testing.T) {
	h := HashcashData{
		Version:    1,
		Difficulty: 4,
		Address:    "test",
		Nonce:      123,
		Counter:    0,
	}

	expected := "1:4:test:123:0"
	if h.String() != expected {
		t.Errorf("String() = %v, want %v", h.String(), expected)
	}
}

func TestHashcashData_IsValid(t *testing.T) {
	h := HashcashData{
		Version:    1,
		Difficulty: 4,
		Address:    "127.0.0.1:63910",
		Nonce:      8671,
		Counter:    29862,
	}

	if !h.IsValid() {
		t.Errorf("IsValid() = false, want true")
	}
	h.Counter = 2
	if h.IsValid() {
		t.Errorf("IsValid() = true, want false")
	}
}
