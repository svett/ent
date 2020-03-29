package schema

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/google/uuid"
)

// Article holds the schema definition for the Article entity.
type Article struct {
	ent.Schema
}

// Fields of the Article.
func (Article) Fields() []ent.Field {
	return []ent.Field{
		field.
			UUID("id", uuid.UUID{}).
			Immutable().
			Default(uuid.New),
		field.
			String("name").
			NotEmpty().
			Unique(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.
			Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Article.
func (Article) Edges() []ent.Edge {
	return nil
}
