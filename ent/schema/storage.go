package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// Storage holds the schema definition for the Storage entity.
type Storage struct {
	ent.Schema
}

// Fields of the Storage.
func (Storage) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.String("secretID").NotEmpty(),
		field.String("secretKey").NotEmpty(),
		field.String("token"),
		field.String("region").NotEmpty(),
		field.String("endpoint").NotEmpty(),
		field.String("bucket").NotEmpty(),
		field.Int64("capacity").Default(0),
		field.Int64("usage").Default(0),
	}
}

func (Storage) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the Storage.
func (Storage) Edges() []ent.Edge {
	return nil
}
