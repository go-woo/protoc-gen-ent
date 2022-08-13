package main

import (
	"bytes"
	"strings"
	"text/template"
)

var entTemplate = `
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// {{.MessageName}} holds the schema definition for the {{.MessageName}} entity.
type {{.MessageName}} struct {
	ent.Schema
}

// Fields of the {{.MessageName}}.
func (Todo) Fields() []ent.Field {
	return []ent.Field{
		field.Bytes("asd"),
	}
}

// Edges of the {{.MessageName}}.
func ({{.MessageName}}) Edges() []ent.Edge {
	return nil
}`
var mixinTemplate = `
import (
	"encoding/base64"
	"net/http"
	"os"
	"strconv"
	"time"
	

)
var _ strconv.NumError
var _ time.Time
var _ base64.CorruptInputError
`

type messageDesc struct {
	MessageType string // User
	MessageName string // User
	SchemaType  string // "Ent"/"Mixin"
	Metadata    string // example/v1/greeter.proto
	HasField    bool
	Fields      []*fieldDesc
}

type fieldDesc struct {
	FieldName  string
	FieldType  string
	FieldRules string
}

type fieldRule struct {
	Rule string
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
