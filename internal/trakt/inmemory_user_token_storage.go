package trakt

import (
	"context"
	"fmt"
)

type InMemoryUserAccessTokenStorage struct {
	tokens map[string]string
}

func NewInMemoryUserAccessTokenStorage() *InMemoryUserAccessTokenStorage {
	return &InMemoryUserAccessTokenStorage{
		tokens: make(map[string]string),
	}
}

func (s *InMemoryUserAccessTokenStorage) StoreAccessToken(_ context.Context, userID string, token string) error {
	s.tokens[userID] = token

	return nil
}

func (s *InMemoryUserAccessTokenStorage) GetAccessToken(_ context.Context, username string) (string, error) {
	token, ok := s.tokens[username]
	if !ok {
		return "", fmt.Errorf("token for user %s not found", username)
	}

	return token, nil
}
