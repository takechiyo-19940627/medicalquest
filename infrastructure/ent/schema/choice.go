package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Choice holds the schema definition for the Choice entity.
type Choice struct {
	ent.Schema
}

// Fields of the Choice.
func (Choice) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Immutable(),
		field.String("uid").
			NotEmpty().
			Immutable().
			Unique(),
		field.Int("question_id").
			Positive(),
		field.String("content").
			NotEmpty().
			Comment("選択肢の内容"),
		field.Bool("is_correct").
			Default(false),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the Choice.
func (Choice) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("question", Question.Type).
			Ref("choices").
			Unique().
			Required().
			Field("question_id"),
	}
}
