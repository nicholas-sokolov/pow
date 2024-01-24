package pow

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestHashcashData_String(t *testing.T) {
	now := int64(1234)
	h := HashcashData{
		Version:    1,
		Difficulty: 4,
		Timestamp:  now,
		Address:    "test",
		Nonce:      123,
		Counter:    0,
	}

	expected := fmt.Sprintf("1:4:%d:test:123:0", now)
	if !strings.HasPrefix(h.String(), expected) {
		t.Errorf("String() = %v, want prefix %v", h.String(), expected)
	}
}

func TestHashcashData_IsValid(t *testing.T) {
	h := HashcashData{
		Version:    1,
		Difficulty: 1,
		Timestamp:  1706126745,
		Address:    "127.0.0.1:56974",
		Nonce:      8237,
		Counter:    3,
	}

	if !h.IsValid() {
		t.Errorf("IsValid() = false, want true")
	}
	h.Counter = 2
	if h.IsValid() {
		t.Errorf("IsValid() = true, want false")
	}
}

func TestHashcashData_IsExpired(t *testing.T) {
	h := HashcashData{
		Timestamp: time.Now().Unix() - 100,
	}
	err := os.Setenv("HASHCASH_EXPIRATION", "60")
	if err != nil {
		fmt.Println("error to set HASHCASH_EXPIRATION to env", err)
	}

	if !h.IsExpired() {
		t.Errorf("IsExpired() = false, want true")
	}

	err = os.Setenv("HASHCASH_EXPIRATION", "101")
	if err != nil {
		t.Errorf("error to set HASHCASH_EXPIRATION to env, %v", err)
	}

	if h.IsExpired() {
		t.Errorf("IsExpired() = true, want false")
	}
}
