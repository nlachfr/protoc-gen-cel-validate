package template

import (
	_ "embed"
	"fmt"
	"text/template"

	"github.com/Neakxs/protocel/cmd/protoc-gen-go-cel-validate/internal/version"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

//go:embed template.go.tmpl
var tmpl string

func GenerateTemplate(v *pluginpb.Version, g *protogen.GeneratedFile) (*template.Template, error) {
	return template.New("").Funcs(template.FuncMap{
		"PluginVersion": func() string {
			return fmt.Sprintf("v%d.%d.%d", version.Major, version.Minor, version.Patch)
		},
		"ProtocVersion": func() string {
			return fmt.Sprintf("v%d.%d.%d", *v.Major, *v.Minor, *v.Patch)
		},
	}).Funcs(template.FuncMap{
		"QualifiedGoIdent": func(imp protogen.GoImportPath, s string) string {
			return g.QualifiedGoIdent(imp.Ident(s))
		},
		"proto": func(s string) string {
			return g.QualifiedGoIdent(protogen.GoImportPath("google.golang.org/protobuf/proto").Ident(s))
		},
		"context": func(s string) string {
			return g.QualifiedGoIdent(protogen.GoImportPath("context").Ident(s))
		},
		"sync": func(s string) string {
			return g.QualifiedGoIdent(protogen.GoImportPath("sync").Ident(s))
		},
		"validate": func(s string) string {
			return g.QualifiedGoIdent(protogen.GoImportPath("github.com/Neakxs/protocel/validate").Ident(s))
		},
		"options": func(s string) string {
			return g.QualifiedGoIdent(protogen.GoImportPath("github.com/Neakxs/protocel/options").Ident(s))
		},
		"cel": func(s string) string {
			return g.QualifiedGoIdent(protogen.GoImportPath("github.com/google/cel-go/cel").Ident(s))
		},
		"fieldmaskpb": func(s string) string {
			return g.QualifiedGoIdent(protogen.GoImportPath("google.golang.org/protobuf/types/known/fieldmaskpb").Ident(s))
		},
		"protoMarshal": func(m proto.Message) []byte {
			if raw, err := proto.Marshal(m); err == nil {
				return raw
			}
			return nil
		},
	}).Parse(tmpl)
}
