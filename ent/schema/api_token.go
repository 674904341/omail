package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// APIToken holds the schema definition for the APIToken entity.
type APIToken struct {
	ent.Schema
}

// Fields of the APIToken.
func (APIToken) Fields() []ent.Field {
	return []ent.Field{
		field.String("token").Unique().NotEmpty(),
		field.String("name").NotEmpty(),
		field.Time("created_at").Default(time.Now),
		field.Time("last_used_at").Optional(),
		field.Bool("revoked").Default(false),
	}
}

// Edges of the APIToken.
func (APIToken) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("api_tokens").Unique().Required(),
	}
}
