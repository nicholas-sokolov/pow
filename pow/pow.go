package pow

import (
	"crypto/sha256"
	"fmt"
	"github.com/nicholas-sokolov/test_task/config"
	"log"
	"math/rand"
	"strings"
)

type HashcashData struct {
	Version    int
	Difficulty int
	Address    string
	Nonce      int
	Counter    uint
}

func NewHashcashData(addr string) *HashcashData {
	return &HashcashData{
		Version:    1,
		Difficulty: config.GetDifficulty(),
		Address:    addr,
		Nonce:      rand.Intn(9999),
		Counter:    0,
	}
}

func (h *HashcashData) String() string {
	return fmt.Sprintf(
		"%d:%d:%s:%d:%d",
		h.Version, h.Difficulty, h.Address, h.Nonce, h.Counter,
	)
}

func (h *HashcashData) getSha256() string {
	sum := sha256.Sum256([]byte(h.String()))
	return fmt.Sprintf("%x", sum)
}

func (h *HashcashData) IsValid() bool {
	prefix := strings.Repeat("0", h.Difficulty)
	return strings.HasPrefix(h.getSha256(), prefix)
}

func SolveChallenge(h HashcashData) *HashcashData {
	for {
		if h.IsValid() {
			log.Println("challenge solved, attempts:", h.Counter)
			return &h
		}

		h.Counter++
	}
}
