package main

import (
	"fmt"
	"github.com/go-woo/protoc-gen-ent/gent"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"os"
	"strings"
)

//type Type uint8
//
//const (
//	TypeInvalid Type = iota
//	TypeBool
//	TypeTime
//	TypeJSON
//	TypeUUID
//	TypeBytes
//	TypeEnum
//	TypeString
//	TypeOther
//	TypeInt8
//	TypeInt16
//	TypeInt32
//	TypeInt
//	TypeInt64
//	TypeUint8
//	TypeUint16
//	TypeUint32
//	TypeUint
//	TypeUint64
//	TypeFloat32
//	TypeFloat64
//	endTypes
//)

// generateSchemaFiles generates ent/scheme/files.
func generateSchemaFiles(gen *protogen.Plugin, file *protogen.File, omitempty bool) []*protogen.GeneratedFile {
	var pgs []*protogen.GeneratedFile
	if len(file.Messages) == 0 {
		return nil
	}
	for _, message := range file.Messages {
		resource := proto.GetExtension(message.Desc.Options(), annotations.E_Resource).(*annotations.ResourceDescriptor)
		if resource == nil {
			continue
		}
		var pg *protogen.GeneratedFile
		if resource.Type == "/ent" {
			pg = generateEntFile(gen, file, message, omitempty, entTemplate)
		}
		if resource.Type == "/mixin" {
			generateMixinFile(gen, file, message, omitempty, mixinTemplate)
		}
		pgs = append(pgs, pg)
	}
	return pgs
}
func generateSchemaFile(gen *protogen.GeneratedFile) *protogen.GeneratedFile {
	return gen
}
func generateEntFile(gen *protogen.Plugin, file *protogen.File, message *protogen.Message, omitempty bool, tpl string) *protogen.GeneratedFile {
	entFile := strings.ToLower(fmt.Sprintf("%v_%v_ent.pb.go",
		file.GeneratedFilenamePrefix, message.Desc.Name()))

	g := gen.NewGeneratedFile(entFile, file.GoImportPath)
	g.P("// Code generated by protoc-gen-ent. DO NOT EDIT.")
	g.P("// versions:")
	g.P(fmt.Sprintf("// - protoc-gen-ent %s", version))
	g.P("// - protoc  ", protocVersion(gen))
	if file.Proto.GetOptions().GetDeprecated() {
		g.P("// ", file.Desc.Path(), " is a deprecated file.")
	} else {
		g.P("// source: ", file.Desc.Path())
	}
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	generateFileContent(message, g, omitempty, entTemplate)
	return g

}
func generateMixinFile(gen *protogen.Plugin, file *protogen.File, message *protogen.Message, omitempty bool, tpl string) *protogen.GeneratedFile {
	mixinFile := strings.ToLower(fmt.Sprintf("%v_%v_mixin.pb.go",
		file.GeneratedFilenamePrefix, message.Desc.Name()))

	g := gen.NewGeneratedFile(mixinFile, file.GoImportPath)
	g.P("// Code generated by protoc-gen-ent. DO NOT EDIT.")
	g.P("// versions:")
	g.P(fmt.Sprintf("// - protoc-gen-ent %s", version))
	g.P("// - protoc  ", protocVersion(gen))
	if file.Proto.GetOptions().GetDeprecated() {
		g.P("// ", file.Desc.Path(), " is a deprecated file.")
	} else {
		g.P("// source: ", file.Desc.Path())
	}
	g.P()
	g.P("package ", file.GoPackageName)

	generateFileContent(message, g, omitempty, mixinTemplate)
	return g

}

// generateFileContent generates the errors definitions, excluding the package statement.
func generateFileContent(message *protogen.Message, g *protogen.GeneratedFile, omitempty bool, tpl string) {
	if message.Desc.Options().(*descriptorpb.MessageOptions).GetDeprecated() {
		g.P(deprecationComment)
	}
	g.P()
	md := &messageDesc{
		MessageType:    "",
		MessageName:    string(message.Desc.Name()),
		HasFields:      false,
		HasEdges:       false,
		HasIndexes:     false,
		HasAnnotations: false,
		HasMixin:       false,
	}
	var fds []*fieldDesc
	for _, field := range message.Fields {
		rules := proto.GetExtension(field.Desc.Options(), gent.E_Field).(*gent.FieldRules)
		if matchType(field) == "Enum" && rules.GetRules() == "" {
			fmt.Fprintf(os.Stderr, "\u001B[31mWARN\u001B[m: %s %s enum must have rule.\n", message.Desc.Name(), field.GoName)
			os.Exit(2)
		}
		if strings.TrimSpace(rules.GetRules()) != "" {
			fds = append(fds, &fieldDesc{
				FieldLine: fmt.Sprintf("field.%v(\"%v\").%v", matchType(field),
					string(field.Desc.Name()), rules.GetRules()),
			})
		} else {
			fds = append(fds, &fieldDesc{
				FieldLine: fmt.Sprintf("field.%v(\"%v\")", matchType(field),
					string(field.Desc.Name())),
			})
		}
		//fmt.Fprintf(os.Stderr, "\tmessage:field-name-type-rules=====%v:%v-%v-%v\n",
		//	message.Desc.Name(), field.GoName, field.Desc.Kind().GoString(), rules.GetRules())
	}
	md.Fields = fds
	if len(message.Fields) != 0 {
		md.HasFields = true
	}

	resource := proto.GetExtension(message.Desc.Options(), annotations.E_Resource).(*annotations.ResourceDescriptor)
	for _, pattern := range resource.GetPattern() {
		splits := strings.Split(strings.ReplaceAll(pattern, "'", "\""), "/")
		if len(splits) < 2 { //no "/" in pattern
			fmt.Fprintf(os.Stderr, "\u001B[31mWARN\u001B[m: %s %s pattern format invalid.\n", message.Desc.Name(), pattern)
			os.Exit(2)
		}
		switch splits[0] {
		case "Edge":
			md.HasEdges = true
			md.Edges = parsePattern(splits)
		case "Indexes":
			md.HasIndexes = true
			md.Indexes = parsePattern(splits)
		case "Annotations":
			md.HasAnnotations = true
			md.Annotations = parsePattern(splits)
		case "Mixin":
			md.HasMixin = true
			md.Mixin = parsePattern(splits)
		}
	}
	g.P(md.execute(tpl))
}

func parsePattern(pattern []string) []*relation {
	var r []*relation
	for i := 1; i < len(pattern); i++ {
		r = append(r, &relation{Relation: pattern[i]})
	}
	return r
}

func matchType(field *protogen.Field) string {
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		return "Bool"
	case protoreflect.EnumKind:
		return "Enum"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "Int32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "Int"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "Uint32"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "Uint64"
	case protoreflect.FloatKind:
		return "Float32"
	case protoreflect.DoubleKind:
		return "Float64"
	case protoreflect.StringKind:
		return "String"
	case protoreflect.BytesKind:
		return "Bytes"
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return parseMessage(field.Message.Desc)
	default:
		return ""
	}
}

func parseMessage(md protoreflect.MessageDescriptor) string {
	switch md.FullName() {
	case "google.protobuf.Timestamp":
		return "Time"
	case "google.protobuf.Duration":
		return "Time"
	case "google.protobuf.DoubleValue":
		return "Float64"
	case "google.protobuf.FloatValue":
		return "Float32"
	case "google.protobuf.Int64Value":
		return "Int"
	case "google.protobuf.Int32Value":
		return "Int32"
	case "google.protobuf.UInt64Value":
		return "Uint64"
	case "google.protobuf.UInt32Value":
		return "Uint32"
	case "google.protobuf.BoolValue":
		return "Bool"
	case "google.protobuf.StringValue":
		return "String"
	case "google.protobuf.BytesValue":
		return "Bytes"
	case "google.protobuf.FieldMask":
		return ""
	case "google.protobuf.Value":
		fmt.Fprintf(os.Stderr, "\u001B[31mWARN\u001B[m: %s field type can't be Value.\n", md.Name())
		os.Exit(2)
		return "Value" // Todo
	case "google.protobuf.Struct":
		fmt.Fprintf(os.Stderr, "\u001B[31mWARN\u001B[m: %s field type can't be struct.\n", md.Name())
		os.Exit(2)
		return "JSON" // Todo
	default:
		fmt.Fprintf(os.Stderr, "\u001B[31mWARN\u001B[m: %s field type was not valid.\n", md.Name())
		os.Exit(2)
		return ""
	}
}

func protocVersion(gen *protogen.Plugin) string {
	v := gen.Request.GetCompilerVersion()
	if v == nil {
		return "(unknown)"
	}
	var suffix string
	if s := v.GetSuffix(); s != "" {
		suffix = "-" + s
	}
	return fmt.Sprintf("v%d.%d.%d%s", v.GetMajor(), v.GetMinor(), v.GetPatch(), suffix)
}

const deprecationComment = "// Deprecated: Do not use."
