package main

import (
	"bytes"
	"strings"
	"text/template"
)

var entTemplate = `
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

type {{.MessageName}} struct {
	ent.Schema
}

func ({{.MessageName}}) Fields() []ent.Field {
	{{- if .HasFields}}
	return []ent.Field{
		{{- range .Fields}}
		{{.FieldLine}},{{end}}
	}
	{{- else}}
	return nil{{end}}
}

func ({{.MessageName}}) Edges() []ent.Edge {
	{{- if .HasEdges}}
	return []ent.Edge{
		{{- range .Edges}}
		{{.Relation}},{{end}}
	}
	{{- else}}
	return nil{{end}}
}

func ({{.MessageName}}) Indexes() []ent.Index {
	{{- if .HasIndexes}}
	return []ent.Index{
		{{- range .Indexes}}
		{{.Relation}},{{end}}
	}
	{{- else}}
	return nil{{end}}
}

func ({{.MessageName}}) Mixin() []ent.Mixin {
	{{- if .HasMixin}}
	return []ent.Mixin{
		{{- range .Mixin}}
		{{.Relation}},{{end}}
	}
	{{- else}}
	return nil{{end}}
}

func ({{.MessageName}}) Annotations() []schema.Annotation {
	{{- if .HasAnnotations}}
	return []schema.Annotation{
		{{- range .Annotations}}
		{{.Relation}},{{end}}
		
	}
	{{- else}}
	return nil{{end}}
}`

var mixinTemplate = `
import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

var _ time.Month

type {{.MessageName}} struct {
	mixin.Schema
}

func ({{.MessageName}}) Fields() []ent.Field {
	return []ent.Field{
		{{- range .Fields}}
		{{.FieldLine}},{{end}}
	}
}`

type messageDesc struct {
	MessageType    string // User
	MessageName    string // User
	HasFields      bool
	Fields         []*fieldDesc
	HasEdges       bool
	Edges          []*relation
	HasIndexes     bool
	Indexes        []*relation
	HasAnnotations bool
	Annotations    []*relation
	HasMixin       bool
	Mixin          []*relation
}

type fieldDesc struct {
	FieldLine string
}

type relation struct {
	Relation string
}

func (s *messageDesc) execute(tpl string) string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("ent").Parse(strings.TrimSpace(tpl))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}
	return strings.Trim(buf.String(), "\r\n")
}
