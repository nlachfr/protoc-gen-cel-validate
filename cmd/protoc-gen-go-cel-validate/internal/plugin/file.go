package plugin

import (
	"github.com/Neakxs/protocel/cmd/protoc-gen-go-cel-validate/internal/template"
	"github.com/Neakxs/protocel/validate"
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
	imports := []protoreflect.FileDescriptor{f.Desc}
	for i := 0; i < f.Desc.Imports().Len(); i++ {
		imports = append(imports, f.Desc.Imports().Get(i).FileDescriptor)
	}
	svcs := []*Service{}
	for i := 0; i < len(f.Services); i++ {
		svcs = append(svcs, NewService(f.Services[i], cfg, imports...))
	}
	msgs := []*Message{}
	for i := 0; i < len(f.Messages); i++ {
		msgs = append(msgs, NewMessage(f.Messages[i], cfg, imports...))
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

func NewService(s *protogen.Service, cfg *validate.ValidateOptions, imports ...protoreflect.FileDescriptor) *Service {
	return &Service{
		Service: s,
		Config:  cfg,
		Imports: imports,
	}
}

type Service struct {
	*protogen.Service
	Imports []protoreflect.FileDescriptor
	Config  *validate.ValidateOptions
}

func (s *Service) Validate() error {
	if _, err := validate.BuildServiceValidateProgram(s.Config, s.Desc, nil, s.Imports...); err != nil {
		return err
	}
	return nil
}

func NewMessage(m *protogen.Message, cfg *validate.ValidateOptions, imports ...protoreflect.FileDescriptor) *Message {
	return &Message{
		Message: m,
		Imports: imports,
		Config:  cfg,
	}
}

type Message struct {
	*protogen.Message
	Imports []protoreflect.FileDescriptor
	Config  *validate.ValidateOptions
}

func (m *Message) Validate() error {
	if _, err := validate.BuildMessageValidateProgram(m.Config, m.Desc, nil, m.Imports...); err != nil {
		return err
	}
	return nil
}

func (m *Message) ContainsValidatePrograms() bool {
	res, _ := validate.BuildMessageValidateProgram(m.Config, m.Desc, nil, m.Imports...)
	return len(res) != 0
}
