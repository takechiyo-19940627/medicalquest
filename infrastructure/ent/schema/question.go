package schema

import (
    "time"
    
    "entgo.io/ent"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
)

// Question holds the schema definition for the Question entity.
type Question struct {
    ent.Schema
}

// Fields of the Question.
func (Question) Fields() []ent.Field {
    return []ent.Field{
        field.Int("id").
            Positive().
            Immutable(),
        field.String("reference_code").
            Optional().
            Nillable(),
        field.String("title").
            NotEmpty(),
        field.Text("content").
            NotEmpty(),
        field.Time("created_at").
            Default(time.Now).
            Immutable(),
    }
}

// Edges of the Question.
func (Question) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("choices", Choice.Type),
    }
}