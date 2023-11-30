package test

import (
	"testing"

	"github.com/eddie023/wex-tag/ent"
	"github.com/eddie023/wex-tag/ent/enttest"

	_ "github.com/mattn/go-sqlite3"
)

func NewDatabase(t *testing.T) *ent.Client {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")

	return client
}
