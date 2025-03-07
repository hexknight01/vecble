package client

import (
	"log"
	"readpebble/internal/storage"
	"time"
)

type client struct {
	storage storage.Storage
}

func (c *client) Insert(key string, value []float64) {
	entry := storage.Entry{
		Key:       key,
		Value:     storage.NewObject(value, storage.ObjectTypeArray),
		ShardID:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	c.storage.Insert(entry)
}

func (c *client) Get(key string) []float64 {

	value, err := c.storage.Get([]byte(key))
	if err != nil {
		log.Print(err)
	}
	return value
}
func NewClient(storage storage.Storage) *client {
	return &client{
		storage: storage,
	}
}
