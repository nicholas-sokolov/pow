package message

import (
	"encoding/json"
	"fmt"
	"github.com/nicholas-sokolov/test_task/config"
	"github.com/nicholas-sokolov/test_task/pow"
	"github.com/nicholas-sokolov/test_task/quote"
	"log"
)

const (
	requestChallenge = iota
	responseChallenge
	requestData
	responseData
)

type Message struct {
	RequestType int
	Header      pow.HashcashData
	Quote       quote.Quote
}

func (m *Message) ChallengeMessage(address string) {
	m.RequestType = responseChallenge
	m.Header = *pow.NewHashcashData(address)
}

func (m *Message) ResponseDataMessage() error {
	// check expiration
	if m.Header.IsExpired() {
		return fmt.Errorf("hashcash is expired")
	}

	// check if difficulty was overridden
	if m.Header.Difficulty != config.GetDifficulty() {
		return fmt.Errorf("difficulty was overridden")
	}

	if !m.Header.IsValid() {
		return fmt.Errorf("hashcash is not valid")
	}

	m.RequestType = responseData
	q, err := quote.GetQuote()
	if err != nil {
		return fmt.Errorf("error to get a quote: %v", err)
	}

	m.Quote = *q

	return nil
}

func (m *Message) ProcessServerMessage(addr string) ([]byte, error) {
	switch {
	case m.isRequestChallenge():
		m.ChallengeMessage(addr)

	case m.isRequestData():
		log.Println("validating challenge")

		err := m.ResponseDataMessage()
		if err != nil {
			return nil, err
		}
	}

	msg, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("error to marshal callenge, %v", err)
	}

	// add '\n' to the message
	msg = append(msg, byte('\n'))

	return msg, nil
}

func (m *Message) ProcessClientMessage() ([]byte, error) {
	switch {
	case m.isResponseChallenge():
		hashcash := pow.SolveChallenge(m.Header)

		m.RequestType = requestData
		m.Header = *hashcash

		return m.preparePayload()

	case m.isResponseData():
		log.Printf("Quote: %s Author: %s", m.Quote.Text, m.Quote.Author)
		return nil, nil

	default:
		return nil, fmt.Errorf("unable to recognize message: %v", m)
	}

}

func NewRequestChallenge() ([]byte, error) {
	msg := Message{
		RequestType: requestChallenge,
	}

	return msg.preparePayload()
}

func (m *Message) preparePayload() ([]byte, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("error to marshal hashcash: %v", err)
	}

	data = append(data, byte('\n'))

	return data, nil
}

func (m *Message) isRequestChallenge() bool {
	return m.RequestType == requestChallenge
}

func (m *Message) isResponseChallenge() bool {
	return m.RequestType == responseChallenge
}

func (m *Message) isRequestData() bool {
	return m.RequestType == requestData
}

func (m *Message) isResponseData() bool {
	return m.RequestType == responseData
}
