// Code generated by protoc-gen-ent. DO NOT EDIT.
// versions:
// - protoc-gen-ent v0.0.1
// - protoc  v3.12.4
// source: example/ent/schema/greeter.proto

package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

var _ time.Month

type DetailsMixin struct {
	mixin.Schema
}

func (DetailsMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.Int("age").Positive(),
		field.String("name").NotEmpty(),
	}
}
