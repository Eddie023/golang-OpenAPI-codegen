package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Transaction holds the schema definition for the Transaction entity.
type Transaction struct {
	ent.Schema
}

// Fields of the Transaction.
func (Transaction) Fields() []ent.Field {
	return []ent.Field{
		// this creates a unique identifier by default
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.Time("date").Default(time.Now),
		// update the default float type by passing custom go type of decimal
		// in postgres we will create numeric column type for more precision
		field.Float("amount_in_usd").GoType(decimal.Decimal{}).SchemaType(map[string]string{
			dialect.Postgres: "numeric",
		}),
		field.String("description").MaxLen(50),
	}
}

// Edges of the Transaction.
func (Transaction) Edges() []ent.Edge {
	return nil
}
