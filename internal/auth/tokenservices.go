package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"
)

// simple tokenizer, use cache-based storage for prod

// reasonable
const (
	tokenLength = 32
	accessTTL   = 15 * time.Minute
	refreshTTL  = 7 * 24 * time.Hour
)

// this has to be instantiated because it stores the token data)
type Tokenizer struct {
	userRefreshTokens map[string]string
	accessExpiry      map[string]time.Time
	refreshExpiry     map[string]time.Time
}

func NewTokenizer() *Tokenizer {
	return &Tokenizer{
		userRefreshTokens: make(map[string]string),
		accessExpiry:      make(map[string]time.Time),
		refreshExpiry:     make(map[string]time.Time),
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

	t.userRefreshTokens[user] = Refresh

	t.accessExpiry[Access] = time.Now().Add(accessTTL)
	t.refreshExpiry[Refresh] = time.Now().Add(refreshTTL)

	return Access, Refresh, nil
}

func (t *Tokenizer) AuthorizeToken(user, access string) error {
	refresh, ok := t.userRefreshTokens[user]
	if !ok {
		return errors.New("no active session")
	}

	expiry, ok := t.accessExpiry[access]
	if !ok {
		return errors.New("no access token for this user")
	}

	if time.Now().After(expiry) {
		return errors.New("token expired")
	}

	refreshExpiry, ok := t.refreshExpiry[refresh]
	if !ok || time.Now().After(refreshExpiry) {
		delete(t.userRefreshTokens, user)
		if ok {
			delete(t.refreshExpiry, refresh)
		}
		return errors.New("refresh token invalid or expired")
	}

	return nil
}

func (t *Tokenizer) RefreshToken(user, oldRefreshToken string) (Access, Refresh string, err error) {
	refresh, ok := t.userRefreshTokens[user]
	if !ok || refresh != oldRefreshToken {
		return "", "", errors.New("invalid refresh token")
	}

	expiry, ok := t.refreshExpiry[oldRefreshToken]
	if !ok {
		return "", "", errors.New("refresh token expiry not found")
	}

	if time.Now().After(expiry) {
		// delete expired sessions
		delete(t.userRefreshTokens, user)
		delete(t.refreshExpiry, oldRefreshToken)
		return "", "", errors.New("refresh token expired")
	}

	delete(t.userRefreshTokens, user)
	delete(t.refreshExpiry, oldRefreshToken)

	return t.GenerateToken(user)
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
