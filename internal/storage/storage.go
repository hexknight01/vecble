/*
 *   Copyright (c) 2025 Vecble
 *   All rights reserved.
 */
package storage

import "github.com/cockroachdb/pebble"

type Storage interface {
}

type storage struct {
	db *pebble.DB
}

// Vector Search
func (s *storage) Search(key []byte) ([]byte, error) {
	return nil, nil
}

func (s *storage) Get(key []byte) ([]byte, error) {
	res, closer, err := s.db.Get(key)
	if err != nil {
		return nil, err
	}
	defer closer.Close()
	return res, nil
}

// SetValue stores a generic slice of numbers (int, float32, float64) as bytes in Pebble
func (s *storage) Set(data Entry) error {

	err := s.db.Set(data.Data.Get(), value, &pebble.WriteOptions{Sync: true})
	if err != nil {
		return err
	}
	return nil
}

func NewStorage(db *pebble.DB) storage {
	return storage{
		db: db,
	}
}
