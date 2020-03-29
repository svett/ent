package schema

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/google/uuid"
)

// Product holds the schema definition for the Product entity.
type Product struct {
	ent.Schema
}

// Fields of the Product.
func (Product) Fields() []ent.Field {
	return []ent.Field{
		field.
			UUID("id", uuid.UUID{}).
			Immutable().
			Default(uuid.New),
		field.
			String("title").
			NotEmpty(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.
			Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Product.
func (Product) Edges() []ent.Edge {
	return nil
}
