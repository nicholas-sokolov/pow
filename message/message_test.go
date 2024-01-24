package message

import (
	"encoding/json"
	"testing"
)

func TestMessage_ChallengeMessage(t *testing.T) {
	m := Message{}
	address := "test_address"

	m.ChallengeMessage(address)

	if !m.isResponseChallenge() {
		t.Errorf("RequestType = %d; want %d", m.RequestType, responseChallenge)
	}

}

func TestNewRequestChallenge(t *testing.T) {
	data, err := NewRequestChallenge()
	if err != nil {
		t.Fatalf("NewRequestChallenge() error = %v", err)
	}

	if len(data) == 0 {
		t.Error("NewRequestChallenge() returned empty data")
	}

	var m Message

	err = json.Unmarshal(data, &m)
	if err != nil {
		t.Errorf("error to unmarshal data, %v", err)
	}

	if !m.isRequestChallenge() {
		t.Errorf("RequestType = %d; want %d", m.RequestType, requestChallenge)
	}
}
