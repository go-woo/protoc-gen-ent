// Code generated by protoc-gen-ent. DO NOT EDIT.
// versions:
// - protoc-gen-ent v0.0.1
// - protoc  v3.12.4
// source: example/ent/schema/greeter.proto

package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
)

var _ time.Month
var _ entsql.Annotation
var _ index.Descriptor

type Pet struct {
	ent.Schema
}

func (Pet) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("name"),
	}
}

func (Pet) Edges() []ent.Edge {
	return nil
}

func (Pet) Indexes() []ent.Index {
	return nil
}

func (Pet) Mixin() []ent.Mixin {
	return nil
}

func (Pet) Annotations() []schema.Annotation {
	return nil
}
