package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Burrow holds the schema definition for the Burrow entity.
type Burrow struct {
	ent.Schema
}

// Fields of the Burrow.
func (Burrow) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Unique(),
		field.String("name").
			Comment("Name of the burrow"),
		field.Float("depth").
			Default(0.0).
			Comment("Current depth of the burrow in meters"),
		field.Float("width").
			Default(0.0).
			Comment("Width of the burrow in meters"),
		field.Bool("is_occupied").
			Default(false).
			Comment("Whether the burrow is currently occupied"),
		field.Int("age"),
		field.Time("updated_at"),
	}
}

func (Burrow) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
	}
}

// Edges of the Burrow.
func (Burrow) Edges() []ent.Edge {
	return nil
}
