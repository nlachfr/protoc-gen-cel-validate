package plugin

import (
	"github.com/Neakxs/protocel/cmd/protoc-gen-go-cel-validate/internal/template"
	"github.com/Neakxs/protocel/validate"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func NewFile(p *protogen.Plugin, f *protogen.File, c *validate.Options) (*File, error) {
	g := p.NewGeneratedFile(f.GeneratedFilenamePrefix+".pb.cel.validate.go", f.GoImportPath)
	cfg := &validate.Options{}
	proto.Merge(cfg, c)
	builder := validate.NewBuilder(
		validate.WithOptions(cfg),
		validate.WithDescriptors(f.Desc),
	)
	svcs := []*Service{}
	for i := 0; i < len(f.Services); i++ {
		svcs = append(svcs, NewService(builder, f.Services[i]))
	}
	msgs := []*Message{}
	for i := 0; i < len(f.Messages); i++ {
		msgs = append(msgs, NewMessage(builder, f.Messages[i]))
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
	Config   *validate.Options
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

func NewService(b validate.Builder, s *protogen.Service) *Service {
	return &Service{
		Service: s,
		Builder: b,
	}
}

type Service struct {
	*protogen.Service
	Builder validate.Builder
}

func (s *Service) Validate() error {
	if _, err := s.Builder.BuildServiceRuleValidater(s.Desc); err != nil {
		return err
	}
	return nil
}

func NewMessage(b validate.Builder, m *protogen.Message) *Message {
	return &Message{
		Message: m,
		Builder: b,
	}
}

type Message struct {
	*protogen.Message
	Builder validate.Builder
}

func (m *Message) Validate() error {
	if _, err := m.Builder.BuildMessageRuleValidater(m.Desc); err != nil {
		return err
	}
	return nil
}

func (m *Message) ContainsValidatePrograms() bool {
	v, _ := m.Builder.BuildMessageRuleValidater(m.Desc)
	return v.HasValidaters()
}
