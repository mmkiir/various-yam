package main

import (
	"encoding/json"
	"os"
	"sync"
)

// FileStorage is a simple key-value storage that persists data to a file.
type FileStorage struct {
	filepath string
	mu       sync.RWMutex
}

// NewFileStorage creates a new FileStorage instance.
func NewFileStorage(filepath string) *FileStorage {
	return &FileStorage{
		filepath: filepath,
	}
}

// WriteFile writes the data to the file.
func (fs *FileStorage) readFile() (map[string]string, error) {
	file, err := os.Open(fs.filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := make(map[string]string)
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

// writeFile writes the data to the file.
func (fs *FileStorage) writeFile(data map[string]string) error {
	file, err := os.Create(fs.filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(data); err != nil {
		return err
	}

	return nil
}

// GetItem returns the value for the given key.
func (fs *FileStorage) GetItem(key string) (string, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	data, err := fs.readFile()
	if err != nil {
		return "", err
	}

	return data[key], nil
}

// SetItem sets the value for the given key.
func (fs *FileStorage) SetItem(key, value string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	data, err := fs.readFile()
	if err != nil {
		data = make(map[string]string)
	}

	data[key] = value

	return fs.writeFile(data)
}

// RemoveItem removes the value for the given key.
func (fs *FileStorage) RemoveItem(key string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	data, err := fs.readFile()
	if err != nil {
		return err
	}

	delete(data, key)

	return fs.writeFile(data)
}

// Clear removes all items from the storage.
func (fs *FileStorage) Clear() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	return fs.writeFile(make(map[string]string))
}

// Close closes the storage.
func (fs *FileStorage) Key(index int) (string, error) {
	data, err := fs.readFile()
	if err != nil {
		return "", err
	}

	i := 0
	for key := range data {
		if i == index {
			return key, nil
		}
		i++
	}

	return "", nil
}
