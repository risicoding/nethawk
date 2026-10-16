package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type Store struct {
	filePath string
	mu       sync.RWMutex
}

type Data struct {
	DeviceName  string `json:"devicename"`
	DeviceUname string `json:"deviceuname"`
	DeviceAlias string `json:"devicealias"`
}

func getStoreDir() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(homedir, ".local", "share", "nethawk", "store")

	return path, nil
}

func GetStore() (*Store, error) {
	storeDir, err := getStoreDir()
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(storeDir, 0755)
	if err != nil {
		return nil, err
	}

	storeFile := filepath.Join(storeDir, "store.json")

	return &Store{filePath: storeFile}, nil
}

func (s *Store) Write(data Data) error {
	v, err := json.Marshal(data)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	return os.WriteFile(s.filePath, v, 0644)
}

func (s *Store) Read() (Data, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	bytes, err := os.ReadFile(s.filePath)
	if err != nil {
		return Data{}, err
	}

	var data Data

	if err := json.Unmarshal(bytes, &data); err != nil {
		return Data{}, err
	}

	return data, nil
}
