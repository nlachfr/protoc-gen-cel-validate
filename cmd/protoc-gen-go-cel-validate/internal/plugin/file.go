package plugin

import (
	"github.com/nlachfr/protoc-gen-cel-validate/cmd/protoc-gen-go-cel-validate/internal/template"
	"github.com/nlachfr/protoc-gen-cel-validate/validate"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func NewFile(p *protogen.Plugin, f *protogen.File, c *validate.Configuration) (*File, error) {
	g := p.NewGeneratedFile(f.GeneratedFilenamePrefix+".pb.cel.validate.go", f.GoImportPath)
	cfg := &validate.Configuration{}
	proto.Merge(cfg, c)
	manager, err := validate.NewManager(f.Desc, validate.WithConfiguration(cfg))
	if err != nil {
		return nil, err
	}
	svcs := []*Service{}
	for i := 0; i < len(f.Services); i++ {
		svcs = append(svcs, NewService(manager, f.Services[i]))
	}
	msgs := []*Message{}
	for i := 0; i < len(f.Messages); i++ {
		msgs = append(msgs, NewMessage(manager, f.Messages[i]))
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
	Config   *validate.Configuration
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

func NewService(b *validate.Manager, s *protogen.Service) *Service {
	return &Service{
		Service: s,
		Manager: b,
	}
}

type Service struct {
	*protogen.Service
	Manager *validate.Manager
}

func (s *Service) Validate() error {
	if _, err := s.Manager.GetServiceRuleValidater(s.Desc); err != nil {
		return err
	}
	return nil
}

func NewMessage(b *validate.Manager, m *protogen.Message) *Message {
	return &Message{
		Message: m,
		Manager: b,
	}
}

type Message struct {
	*protogen.Message
	Manager *validate.Manager
}

func (m *Message) Validate() error {
	if _, err := m.Manager.GetMessageRuleValidater(m.Desc); err != nil {
		return err
	}
	return nil
}

func (m *Message) ContainsValidatePrograms() bool {
	v, _ := m.Manager.GetMessageRuleValidater(m.Desc)
	return v.HasValidaters()
}
