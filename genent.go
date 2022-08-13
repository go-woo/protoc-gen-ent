package main

import (
	"encoding/json"
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

//var methodSets = make(map[string]int)

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
			rj, _ := json.Marshal(resource)
			fmt.Fprintf(os.Stderr, "ent=====%v\n", string(rj))
			for _, field := range message.Fields {
				fmt.Fprintf(os.Stderr, "\tent message:field-name-type=====%v:%v-%v\n",
					message.Desc.Name(), field.GoName, field.Desc.Kind().String())
			}
		}
		if resource.Type == "/mixin" {
			generateMixinFile(gen, file, message, omitempty, mixinTemplate)
			rj, _ := json.Marshal(resource)
			fmt.Fprintf(os.Stderr, "mixin=====%v\n", string(rj))
			for _, field := range message.Fields {
				fmt.Fprintf(os.Stderr, "\tmixin message:field-name-type=====%v:%v-%v\n",
					message.Desc.Name(), field.GoName, field.Desc.Kind().String())
			}
		}
		pgs = append(pgs, pg)
	}
	return pgs
}
func generateSchemaFile(gen *protogen.GeneratedFile) *protogen.GeneratedFile {
	return gen
}
func generateEntFile(gen *protogen.Plugin, file *protogen.File, message *protogen.Message, omitempty bool, tpl string) *protogen.GeneratedFile {
	entFile := fmt.Sprintf("%v_%v_ent.pb.go", file.GeneratedFilenamePrefix, message.Desc.Name())

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
	generateEntFileContent(message, g, omitempty, entTemplate)
	return g

}
func generateMixinFile(gen *protogen.Plugin, file *protogen.File, message *protogen.Message, omitempty bool, tpl string) *protogen.GeneratedFile {
	mixinFile := fmt.Sprintf("%v_%v_mixin.pb.go", file.GeneratedFilenamePrefix, message.Desc.Name())

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
	g.P()
	generateMixinFileContent(message, g, omitempty, mixinTemplate)
	return g

}

// generateEntFileContent generates the errors definitions, excluding the package statement.
func generateEntFileContent(message *protogen.Message, g *protogen.GeneratedFile, omitempty bool, tpl string) {
	g.P()
	md := messageDesc{
		MessageName: string(message.Desc.Name()),
		SchemaType:  "Ent",
		Metadata:    "",
		HasField:    false,
		Fields:      nil,
	}
	var fds []*fieldDesc
	for _, field := range message.Fields {
		rules := proto.GetExtension(field.Desc.Options(), gent.E_Field).(*gent.FieldRules)
		fds = append(fds, &fieldDesc{
			FieldName:  string(field.Desc.Name()),
			FieldType:  matchType(field.Desc.Kind().String()),
			FieldRules: rules.GetRules(),
		})
		fmt.Fprintf(os.Stderr, "\tmixin message:field-name-type-rules=====%v:%v-%v-%v\n",
			message.Desc.Name(), field.GoName, field.Desc.Kind().GoString(), rules.GetRules())

	}

	if len(message.Fields) != 0 {
		md.HasField = true
	}
	g.P(md.execute(tpl))
}

// generateMixinFileContent generates the errors definitions, excluding the package statement.
func generateMixinFileContent(message *protogen.Message, g *protogen.GeneratedFile, omitempty bool, tpl string) {
	g.P()
	md := messageDesc{
		MessageName: string(message.Desc.Name()),
		SchemaType:  "Mixin",
		Metadata:    "",
		HasField:    false,
		Fields:      nil,
	}
	var fds []*fieldDesc
	for _, field := range message.Fields {
		rules := proto.GetExtension(field.Desc.Options(), gent.E_Field).(*gent.FieldRules)
		fds = append(fds, &fieldDesc{
			FieldName:  string(field.Desc.Name()),
			FieldType:  matchType(field.Desc.Kind().String()),
			FieldRules: rules.GetRules(),
		})
		fmt.Fprintf(os.Stderr, "\tmixin message:field-name-type-rules=====%v:%v-%v-%v\n",
			message.Desc.Name(), field.GoName, field.Desc.Kind().GoString(), rules.GetRules())

	}

	if len(message.Fields) != 0 {
		md.HasField = true
	}
	g.P(md.execute(tpl))
}

func matchType(kind string) string {
	return "Int"
}
func genMessage(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, message *protogen.Message, omitempty bool, tpl string) {
	if message.Desc.Options().(*descriptorpb.MessageOptions).GetDeprecated() {
		g.P("//")
		g.P(deprecationComment)
	}
	// if message miss id field will exit

}

//func buildMessageDesc(g *protogen.GeneratedFile, m *protogen.Method, method, path string, host string, scopes []string) *methodDesc {
//	return nil
//}

func buildExpr(protoName, goName string, fd protoreflect.FieldDescriptor) string {
	if fd.IsMap() || fd.IsList() { //google http rule do not support
		return "return http.ErrNotSupported"
	}
	switch fd.Kind() {
	case protoreflect.BoolKind:
		return fmt.Sprintf(`if cv, err := strconv.ParseBool(v); err != nil {
			return http.ErrNotSupported
		}else{
			req.%v = cv
		}`, goName)
	case protoreflect.EnumKind: //Todo
		return fmt.Sprintf("return http.ErrNotSupported")

	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return fmt.Sprintf(`if cv, err := strconv.ParseInt(v, 10, 32); err != nil {
			return http.ErrNotSupported
		}else{
			req.%v = int32(cv)
		}`, goName)
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return fmt.Sprintf(`if cv, err := strconv.ParseInt(v, 10, 64); err != nil {
			return http.ErrNotSupported
		}else{
			req.%v = cv
		}`, goName)
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return fmt.Sprintf(`if cv, err := strconv.ParseUint(v, 10, 32); err != nil {
			return http.ErrNotSupported
		}else{
			req.%v = uint32(cv)
		}`, goName)
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return fmt.Sprintf(`if cv, err := strconv.ParseUint(v, 10, 64); err != nil {
			return http.ErrNotSupported
		}else{
			req.%v = cv
		}`, goName)
	case protoreflect.FloatKind:
		return fmt.Sprintf(`if cv, err := strconv.ParseFloat(v, 32); err != nil {
			return http.ErrNotSupported
		}else{
			req.%v = float32(cv)
		}`, goName)
	case protoreflect.DoubleKind:
		return fmt.Sprintf(`if cv, err := strconv.ParseFloat(v, 64); err != nil {
			return http.ErrNotSupported
		}else{
			req.%v = float32(cv)
		}`, goName)
	case protoreflect.StringKind:
		return fmt.Sprintf("req.%v = v", goName)
	case protoreflect.BytesKind:
		return fmt.Sprintf(`if cv, err := base64.StdEncoding.DecodeString(v); err != nil {
			return http.ErrNotSupported
		}else{
			req.%v = cv
		}`, goName)
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return parseMessage(fd.Message(), goName)
	default:
		return fmt.Sprintf("return http.ErrNotSupported")
	}
}

func parseMessage(md protoreflect.MessageDescriptor, goName string) string {
	switch md.FullName() {
	case "google.protobuf.Timestamp":
		return fmt.Sprintf(`if cv, err := time.Parse(time.RFC3339Nano, v); err != nil {
			return http.ErrNotSupported
		}else{
			req.%v = cv
		}`, goName)
	case "google.protobuf.Duration":
		return fmt.Sprintf(`if cv, err := time.ParseDuration(v); err != nil {
			return http.ErrNotSupported
		}else{
			req.%v = cv
		}`, goName)
	case "google.protobuf.DoubleValue":
		return fmt.Sprintf(`if cv, err := strconv.ParseFloat(v, 64); err != nil {
			return http.ErrNotSupported
		}else{
			req.%v = cv
		}`, goName)
	case "google.protobuf.FloatValue":
		return fmt.Sprintf(`if cv, err := strconv.ParseFloat(v, 32); err != nil {
			return http.ErrNotSupported
		}else{
			req.%v = float32(cv)
		}`, goName)
	case "google.protobuf.Int64Value":
		return fmt.Sprintf(`if cv, err := strconv.ParseInt(v, 10, 64); err != nil {
			return http.ErrNotSupported
		}else{
			req.%v = cv
		}`, goName)
	case "google.protobuf.Int32Value":
		return fmt.Sprintf(`if cv, err := strconv.ParseInt(v, 10, 32); err != nil {
			return http.ErrNotSupported
		}else{
			req.%v = int32(cv)
		}`, goName)
	case "google.protobuf.UInt64Value":
		return fmt.Sprintf(`if cv, err := strconv.ParseUint(v, 10, 64); err != nil {
			return http.ErrNotSupported
		}else{
			req.%v = cv
		}`, goName)
	case "google.protobuf.UInt32Value":
		return fmt.Sprintf(`if cv, err := strconv.ParseUint(v, 10, 32); err != nil {
			return http.ErrNotSupported
		}else{
			req.%v = uint32(cv)
		}`, goName)
	case "google.protobuf.BoolValue":
		return fmt.Sprintf(`if cv, err := strconv.ParseBool(v); err != nil {
			return http.ErrNotSupported
		}else{
			req.%v = cv
		}`, goName)
	case "google.protobuf.StringValue":
		return fmt.Sprintf("req.%v = v", goName)
	case "google.protobuf.BytesValue":
		return fmt.Sprintf(`if cv, err := base64.StdEncoding.DecodeString(v); err != nil {
			return http.ErrNotSupported
		}else{
			req.%v = cv
		}`, goName)
	case "google.protobuf.FieldMask":
		return "return http.ErrNotSupported"
	case "google.protobuf.Value":
		return "return http.ErrNotSupported"
	case "google.protobuf.Struct":
		return "return http.ErrNotSupported"
	default:
		return "return http.ErrNotSupported"
	}

	return "return http.ErrNotSupported"
}

func camelCaseVars(s string) string {
	subs := strings.Split(s, ".")
	vars := make([]string, 0, len(subs))
	for _, sub := range subs {
		vars = append(vars, camelCase(sub))
	}
	return strings.Join(vars, ".")
}

// camelCase returns the CamelCased name.
// If there is an interior underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
// There is a remote possibility of this rewrite causing a name collision,
// but it's so remote we're prepared to pretend it's nonexistent - since the
// C++ generator lowercases names, it's extremely unlikely to have two fields
// with different capitalizations.
// In short, _my_field_name_2 becomes XMyFieldName_2.
func camelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
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
