/*
 *   Copyright (c) 2025 Vecble
 *   All rights reserved.
 */
package storage

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"

	"github.com/cockroachdb/pebble"
)

type Storage interface {
	Search(key []byte) ([]byte, error)
	Get(key []byte) ([]float64, error)
	Insert(data Entry) error
}

type storage struct {
	db *pebble.DB
}

// Vector Search
func (s *storage) Search(key []byte) ([]byte, error) {
	return nil, nil
}

func (s *storage) Get(key []byte) ([]float64, error) {
	res, closer, err := s.db.Get(key)
	if err != nil {
		return nil, err
	}
	defer closer.Close()
	resFloat, err := deserializeFloat64Array(res)
	if err != nil {
		return nil, err
	}
	return resFloat, nil
}

func serializeFloat64Array(arr []float64) ([]byte, error) {
	size := len(arr) * 8
	bytes := make([]byte, size)
	for i, val := range arr {
		binary.LittleEndian.PutUint64(bytes[i*8:(i+1)*8], uint64(math.Float64bits(val)))
	}
	return bytes, nil
}

func deserializeFloat64Array(bytes []byte) ([]float64, error) {
	if len(bytes)%8 != 0 {
		return nil, fmt.Errorf("invalid byte slice length for float64 array")
	}
	count := len(bytes) / 8
	arr := make([]float64, count)
	for i := range arr {
		arr[i] = math.Float64frombits(uint64(binary.LittleEndian.Uint64(bytes[i*8 : (i+1)*8])))
	}
	return arr, nil
}

func calculateDistance(v1, v2 []float64) float64 {
	if len(v1) != len(v2) {
		log.Fatal("Vectors must be of the same dimension")
	}
	var sum float64
	for i := range v1 {
		diff := v1[i] - v2[i]
		sum += diff * diff
	}
	return math.Sqrt(sum)
}

// SetValue stores a generic slice of numbers (int, float32, float64) as bytes in Pebble
func (s *storage) Insert(entry Entry) error {
	if entry.Value.ObjectType == ObjectTypeArray {
		data := entry.Value.Value.([]float64)
		dataToInsert, err := serializeFloat64Array(data)
		if err != nil {
			log.Print(err)
		}
		err = s.db.Set([]byte(entry.Key), dataToInsert, &pebble.WriteOptions{
			Sync: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func NewStorage(db *pebble.DB) storage {
	return storage{
		db: db,
	}
}
