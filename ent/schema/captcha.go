package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// Captcha holds the schema definition for the Captcha entity.
type Captcha struct {
	ent.Schema
}

// Fields of the Captcha.
func (Captcha) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("key", uuid.UUID{}).Unique().Default(uuid.New),
		field.Text("text").NotEmpty(),
		field.Time("expired_time").Default(func() time.Time {
			return time.Now().Add(time.Minute * 5)
		}),
	}

}

func (Captcha) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the Captcha.
func (Captcha) Edges() []ent.Edge {
	return nil
}
