package translation

import (
	"context"
	"fmt"

	"github.com/jon-at-github/hello-api/config"
	"github.com/jon-at-github/hello-api/handlers/rest"
	"github.com/redis/go-redis/v9"
)

var _ rest.Translator = &Database{}

// Database has a Redis client representing a pool of zero or more underlying connections.
type Database struct {
	conn *redis.Client
}

// NewDatabaseService creates a new instance of a redis database.
func NewDatabaseService(cfg config.Configuration) *Database {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.DatabaseURL, cfg.DatabasePort),
		Password: "",
		DB:       0,
	})
	return &Database{
		conn: rdb,
	}
}

// Close closes the client, releasing any open resources.
func (s *Database) Close() error {
	err := s.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

// Translate will take a given word and language and find a translation for it.
func (s *Database) Translate(word string, language string) string {
	out := s.conn.Get(context.Background(), fmt.Sprintf("%s:%s", word, language))
	return out.Val()
}
