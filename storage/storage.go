package storage

import (
	"errors"
	"sync"
)

type Storage struct {
	urls map[string]string
	mux  sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		urls: map[string]string{},
	}
}

func (s *Storage) Store(shortURL, longURL string) error {
	if shortURL == "" || longURL == "" {
		return errors.New("shortURL & longURL are required")
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	s.urls[shortURL] = longURL
	return nil
}

func (s *Storage) Get(shortURL string) (string, error) {
	if shortURL == "" {
		return "", errors.New("shortURL are required")
	}

	s.mux.RLock()
	defer s.mux.RUnlock()
	longURL, ok := s.urls[shortURL]
	if !ok {
		return "", errors.New("not found any realated longURL with this shortURL")
	}

	return longURL, nil
}
