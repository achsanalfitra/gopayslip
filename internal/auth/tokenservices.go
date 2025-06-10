package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

// simple tokenizer, use cache-based storage for prod

// reasonable
const (
	tokenLength = 32
)

// this has to be instantiated because it stores the token data
type Tokenizer struct {
	tokenStore map[string]map[string]string
}

func NewTokenizer() *Tokenizer {
	return &Tokenizer{
		tokenStore: make(map[string]map[string]string),
	}
}

func (t *Tokenizer) GenerateToken(user string) (Access, Refresh string, err error) {
	accessBytes, err := t.generateRandomBytes(tokenLength)
	if err != nil {
		return "", "", errors.New("failed to generate access token bytes")
	}
	aTokenHash := sha256.Sum256(accessBytes)
	Access = hex.EncodeToString(aTokenHash[:])

	refreshBytes, err := t.generateRandomBytes(tokenLength)
	if err != nil {
		return "", "", errors.New("failed to generate refresh token bytes")
	}
	rTokenHash := sha256.Sum256(refreshBytes)
	Refresh = hex.EncodeToString(rTokenHash[:])

	if _, ok := t.tokenStore[user]; !ok {
		t.tokenStore[user] = make(map[string]string)
	}
	t.tokenStore[user][Refresh] = Access

	return Access, Refresh, nil
}

// helper for GenerateToken
func (t *Tokenizer) generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
