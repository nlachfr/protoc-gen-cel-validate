package plugin

import (
	"fmt"

	"github.com/Neakxs/protocel/cmd/protoc-gen-go-cel-validate/internal/template"
	"github.com/Neakxs/protocel/validate"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func NewFile(p *protogen.Plugin, f *protogen.File, c *validate.ValidateOptions) (*File, error) {
	g := p.NewGeneratedFile(f.GeneratedFilenamePrefix+".pb.cel.validate.go", f.GoImportPath)
	cfg := &validate.ValidateOptions{}
	proto.Merge(cfg, c)
	fileRule := proto.GetExtension(f.Desc.Options(), validate.E_File).(*validate.ValidateOptions)
	if fileRule != nil {
		proto.Merge(cfg, fileRule)
	}
	imports := []protoreflect.FileDescriptor{}
	for _, imp := range p.Files {
		imports = append(imports, imp.Desc)
	}
	resourceMap, err := validate.GenerateResourceTypePatternMapping(f.Desc, imports...)
	if err != nil {
		return nil, err
	}
	svcs := []*Service{}
	for i := 0; i < len(f.Services); i++ {
		svcs = append(svcs, NewService(f.Services[i], cfg, p.Files...))
	}
	msgs := []*Message{}
	for i := 0; i < len(f.Messages); i++ {
		msgs = append(msgs, NewMessage(f.Messages[i], resourceMap, cfg, p.Files...))
	}
	return &File{
		p:        p,
		g:        g,
		File:     f,
		Services: svcs,
		Messages: msgs,
		Config:   cfg,
	}, nil
}

type File struct {
	p *protogen.Plugin
	g *protogen.GeneratedFile
	*protogen.File
	Services []*Service
	Messages []*Message
	Config   *validate.ValidateOptions
}

func (f *File) Generate() error {
	if err := f.Validate(); err != nil {
		return err
	}
	if tmpl, err := template.GenerateTemplate(f.p.Request.CompilerVersion, f.g); err != nil {
		return err
	} else {
		return tmpl.Execute(f.g, f)
	}
}

func (f *File) Validate() error {
	for i := 0; i < len(f.Services); i++ {
		if err := f.Services[i].Validate(); err != nil {
			return err
		}
	}
	for i := 0; i < len(f.Messages); i++ {
		if err := f.Messages[i].Validate(); err != nil {
			return err
		}
	}
	return nil
}

func NewService(s *protogen.Service, cfg *validate.ValidateOptions, imports ...*protogen.File) *Service {
	methods := []*Method{}
	for i := 0; i < len(s.Methods); i++ {
		methods = append(methods, NewMethod(s.Methods[i], cfg, imports...))
	}
	return &Service{
		Service: s,
		Config:  cfg,
		Imports: imports,
		Methods: methods,
	}
}

type Service struct {
	*protogen.Service
	Imports []*protogen.File
	Config  *validate.ValidateOptions
	Methods []*Method
}

func (s *Service) Validate() error {
	for i := 0; i < len(s.Methods); i++ {
		if err := s.Methods[i].Validate(); err != nil {
			return err
		}
	}
	return nil
}

func NewMethod(m *protogen.Method, cfg *validate.ValidateOptions, imports ...*protogen.File) *Method {
	return &Method{
		Method:  m,
		Config:  cfg,
		Imports: imports,
	}
}

type Method struct {
	*protogen.Method
	Imports []*protogen.File
	Config  *validate.ValidateOptions
}

func (m *Method) MethodRule() *validate.ValidateRule {
	return proto.GetExtension(m.Desc.Options(), validate.E_Method).(*validate.ValidateRule)
}

func (m *Method) Validate() error {
	imports := []protoreflect.FileDescriptor{}
	for i := 0; i < len(imports); i++ {
		imports = append(imports, m.Imports[i].Desc)
	}
	rule := m.MethodRule()
	if rule == nil {
		return nil
	}
	exprs := []string{}
	if rule.Expr != "" {
		exprs = append(exprs, rule.Expr)
	}
	if len(rule.Exprs) > 0 {
		exprs = append(exprs, rule.Exprs...)
	}
	if _, err := validate.BuildMethodValidateProgram(exprs, m.Config, m.Desc, nil, imports...); err != nil {
		return err
	}
	return nil
}

func NewMessage(m *protogen.Message, resourceMap map[string]string, cfg *validate.ValidateOptions, imports ...*protogen.File) *Message {
	fields := []*Field{}
	for i := 0; i < len(m.Fields); i++ {
		fields = append(fields, NewField(m.Fields[i], resourceMap, cfg, imports...))
	}
	return &Message{
		Message: m,
		Imports: imports,
		Config:  cfg,
		Fields:  fields,
	}
}

type Message struct {
	*protogen.Message
	Imports []*protogen.File
	Config  *validate.ValidateOptions
	Fields  []*Field
}

func (m *Message) MessageRule() *validate.ValidateRule {
	return proto.GetExtension(m.Desc.Options(), validate.E_Message).(*validate.ValidateRule)
}

func (m *Message) Validate() error {
	for i := 0; i < len(m.Fields); i++ {
		if err := m.Fields[i].Validate(); err != nil {
			return err
		}
	}
	return nil
}

func NewField(f *protogen.Field, resourceMap map[string]string, cfg *validate.ValidateOptions, imports ...*protogen.File) *Field {
	return &Field{Imports: imports, Field: f, ResourceMap: resourceMap, Config: cfg}
}

type Field struct {
	Imports     []*protogen.File
	ResourceMap map[string]string
	Config      *validate.ValidateOptions
	*protogen.Field
}

func (f *Field) FieldRule() *validate.ValidateRule {
	return proto.GetExtension(f.Desc.Options(), validate.E_Field).(*validate.ValidateRule)
}

func (f *Field) HasFieldBehaviorRequired() bool {
	behaviors := proto.GetExtension(f.Desc.Options(), annotations.E_FieldBehavior).([]annotations.FieldBehavior)
	for i := 0; i < len(behaviors); i++ {
		if behaviors[i] == annotations.FieldBehavior_REQUIRED {
			return true
		}
	}
	return false
}

func (f *Field) HasResourceReference() bool {
	return proto.GetExtension(f.Desc.Options(), annotations.E_ResourceReference).(*annotations.ResourceReference) != nil
}

func (f *Field) GetResourceReferenceValidate() string {
	var regexp string
	if f.Config != nil && f.Config.ResourceReferenceSupportDisabled {
		return ""
	}
	if ref := proto.GetExtension(f.Desc.Options(), annotations.E_ResourceReference).(*annotations.ResourceReference); ref != nil {
		if ref.Type != "" {
			if ref.Type != "*" {
				regexp = fmt.Sprintf("^%s$", f.ResourceMap[ref.Type])
			}
		} else if ref.ChildType != "" {
			regexp = fmt.Sprintf("^%s", f.ResourceMap[ref.ChildType])
		}
	}
	if regexp != "" {
		if f.Desc.IsList() {
			return fmt.Sprintf(`%s.all(s, s.matches("%s"))`, f.Desc.TextName(), regexp)
		} else if f.Desc.Kind() == protoreflect.StringKind {
			return fmt.Sprintf(`%s.matches("%s")`, f.Desc.TextName(), regexp)
		}
	}
	return ""
}

func (f *Field) Validate() error {
	imports := []protoreflect.FileDescriptor{}
	for i := 0; i < len(imports); i++ {
		imports = append(imports, f.Imports[i].Desc)
	}
	rule := f.FieldRule()
	exprs := []string{}
	if f.HasResourceReference() {
		if s := f.GetResourceReferenceValidate(); s == "" {
			return fmt.Errorf("cannot build resource reference validate")
		} else {
			exprs = append(exprs, s)
		}
	}
	if rule != nil {
		if rule.Expr != "" {
			exprs = append(exprs, rule.Expr)
		}
		if len(rule.Exprs) > 0 {
			exprs = append(exprs, rule.Exprs...)
		}
	}
	if len(exprs) > 0 {
		if _, err := validate.BuildValidateProgram(exprs, f.Config, f.Parent.Desc, nil, imports...); err != nil {
			return err
		}
	}
	return nil
}
