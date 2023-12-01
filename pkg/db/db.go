package db

import (
	"fmt"
	"log/slog"
	"testing"

	"github.com/eddie023/wex-tag/ent"
	"github.com/eddie023/wex-tag/ent/enttest"
	"github.com/eddie023/wex-tag/pkg/config"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Client *ent.Client
}

// Initiate a new db connection
func NewConnection(cfg *config.ApiConfig) (*DB, error) {
	client, err := connectDb(cfg)
	if err != nil {
		return nil, err
	}

	return &DB{
		Client: client,
	}, nil
}

func connectDb(cfg *config.ApiConfig) (*ent.Client, error) {
	connectionURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", cfg.Db.Host, cfg.Db.Port, cfg.Db.User, cfg.Db.Dbname, cfg.Db.Password)

	slog.Debug("connecting to db", "database", "postgres", "connection-url", connectionURL)

	client, err := ent.Open("postgres", connectionURL)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// CreateTestDatabase creates a dummy sqlite3 ent client which can be used for writing test purposes.
func CreateTestDatabase(t *testing.T) *ent.Client {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")

	return client
}
