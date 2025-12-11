package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// Mailbox holds the schema definition for the Mailbox entity.
type Mailbox struct {
	ent.Schema
}

// Fields of the Mailbox.
func (Mailbox) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").NotEmpty(),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Mailbox.
func (Mailbox) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("mailboxes").Unique().Required(),
		edge.To("envelopes", Envelope.Type),
	}
}
