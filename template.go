package main

var messageTemplate = `
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

//func (s *messageDesc) execute(tpl string) string {
//	s.MethodSets = make(map[string]*methodDesc)
//	for _, m := range s.Methods {
//		s.MethodSets[m.Name] = m
//	}
//
//	buf := new(bytes.Buffer)
//	tmpl, err := template.New("http").Parse(strings.TrimSpace(tpl))
//	if err != nil {
//		panic(err)
//	}
//	if err := tmpl.Execute(buf, s); err != nil {
//		panic(err)
//	}
//	return strings.Trim(buf.String(), "\r\n")
//}
