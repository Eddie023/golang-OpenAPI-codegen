package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Transaction holds the schema definition for the Transaction entity.
type Transaction struct {
	ent.Schema
}

// Fields of the Transaction.
func (Transaction) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()),
		field.Time("date").Default(time.Now),
		field.Float("amount_in_usd").Positive(),
		field.String("description").MaxLen(50),
	}
}

// Edges of the Transaction.
func (Transaction) Edges() []ent.Edge {
	return nil
}
