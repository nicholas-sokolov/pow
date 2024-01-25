package message

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nicholas-sokolov/test_task/config"
	"github.com/nicholas-sokolov/test_task/pow"
	"github.com/nicholas-sokolov/test_task/quote"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
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

func getChallengeMessage(addr string) *Message {
	return &Message{
		RequestType: responseChallenge,
		Header:      *pow.NewHashcashData(addr),
	}
}

func getResponseDataMessage(msg Message) (*Message, error) {
	var respMessage Message

	if !msg.Header.IsValid() {
		return nil, fmt.Errorf("hashcash is not valid")
	}

	msg.RequestType = responseData
	q, err := quote.GetQuote()
	if err != nil {
		return nil, fmt.Errorf("error to get a quote: %v", err)
	}

	respMessage.RequestType = responseData
	respMessage.Quote = *q

	return &respMessage, nil
}

func (m *Message) ProcessServerMessage(ctx context.Context, addr string) ([]byte, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.GetRedisHost(),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	var msg *Message

	switch {
	case m.isRequestChallenge():

		msg = getChallengeMessage(addr)

		cachedData, err := getCacheValue(*msg)
		if err != nil {
			return nil, fmt.Errorf("error to convert value for cache: %v", err)
		}

		expiration := time.Duration(config.GetExpiration()) * time.Second
		err = rdb.SetEx(ctx, addr, cachedData, expiration).Err()
		if err != nil {
			return nil, fmt.Errorf("error to set cache: %v", err)
		}

	case m.isRequestData():
		log.Println("validating challenge")

		val, err := rdb.GetDel(ctx, addr).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return nil, fmt.Errorf("no cached value, probably expired")
			}
			return nil, fmt.Errorf("error to get cache: %v", err)
		}

		var cachedHeader pow.HashcashData
		err = json.Unmarshal([]byte(val), &cachedHeader)
		if err != nil {
			return nil, fmt.Errorf("error to unmarshal cached value: %v", err)
		}

		cachedHeader.Counter = m.Header.Counter
		if cachedHeader != m.Header {
			return nil, fmt.Errorf("cached and value from request are not equal")
		}

		msg, err = getResponseDataMessage(*m)
		if err != nil {
			return nil, fmt.Errorf("failed to get response data: %v", err)
		}
	default:
		return nil, fmt.Errorf("unknonw message type: %d", m.RequestType)
	}

	return msg.preparePayload()
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

func getCacheValue(msg Message) ([]byte, error) {
	msg.Header.Counter = 0
	data, err := json.Marshal(msg.Header)
	if err != nil {
		return nil, fmt.Errorf("error to marshal header")
	}
	return data, nil
}
