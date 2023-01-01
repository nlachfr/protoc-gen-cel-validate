package plugin

import (
	"testing"

	"github.com/Neakxs/protocel/testdata/cmd/plugin"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/pluginpb"
)

func newEnum(desc protoreflect.EnumDescriptor) *protogen.Enum {
	e := &protogen.Enum{
		Desc:   desc,
		Values: []*protogen.EnumValue{},
	}
	for i := 0; i < desc.Values().Len(); i++ {
		desc := desc.Values().Get(i)
		e.Values = append(e.Values, &protogen.EnumValue{
			Desc:    desc,
			GoIdent: protogen.GoIdent{GoName: string(desc.Name())},
			Parent:  e,
		})
	}
	return e
}

func newField(parent *protogen.Message, desc protoreflect.FieldDescriptor) *protogen.Field {
	f := &protogen.Field{
		Desc:   desc,
		GoName: desc.TextName(),
		Parent: parent,
	}
	if v := desc.ContainingOneof(); v != nil {
		f.Oneof = newOneof(f.Parent, v)
	}
	if v := desc.ContainingMessage(); v != nil {
		f.Extendee = parent
	}
	if v := desc.Enum(); v != nil {
		f.Enum = newEnum(v)
	}
	if v := desc.Message(); v != nil {
		f.Message = newMessage(v)
	}
	return f
}

func newOneof(parent *protogen.Message, desc protoreflect.OneofDescriptor) *protogen.Oneof {
	o := &protogen.Oneof{
		Desc:   desc,
		GoName: string(desc.Name()),
		Fields: []*protogen.Field{},
		Parent: parent,
	}
	for i := 0; i < desc.Fields().Len(); i++ {
		o.Fields = append(o.Fields, newField(parent, desc.Fields().Get(i)))
	}
	return o
}

func newMessage(desc protoreflect.MessageDescriptor) *protogen.Message {
	m := &protogen.Message{
		Desc: desc,
		GoIdent: protogen.GoIdent{
			GoName: string(desc.Name()),
		},
		Fields:     []*protogen.Field{},
		Oneofs:     []*protogen.Oneof{},
		Enums:      []*protogen.Enum{},
		Messages:   []*protogen.Message{},
		Extensions: []*protogen.Extension{},
	}
	for i := 0; i < desc.Fields().Len(); i++ {
		m.Fields = append(m.Fields, newField(m, desc.Fields().Get(i)))
	}
	for i := 0; i < desc.Oneofs().Len(); i++ {
		m.Oneofs = append(m.Oneofs, newOneof(m, desc.Oneofs().Get(i)))
	}
	for i := 0; i < desc.Enums().Len(); i++ {
		m.Enums = append(m.Enums, newEnum(desc.Enums().Get(i)))
	}
	for i := 0; i < desc.Messages().Len(); i++ {
		m.Messages = append(m.Messages, newMessage(desc.Messages().Get(i)))
	}
	for i := 0; i < desc.Extensions().Len(); i++ {
		m.Extensions = append(m.Extensions, newExtension(m, desc.Extensions().Get(i)))
	}
	return m
}

func newExtension(parent *protogen.Message, desc protoreflect.ExtensionDescriptor) *protogen.Extension {
	e := &protogen.Field{
		GoName: desc.TextName(),
		Desc:   desc,
		GoIdent: protogen.GoIdent{
			GoName: string(desc.Parent().Name()) + "_" + desc.TextName(),
		},
		Parent: parent,
	}
	if v := desc.ContainingOneof(); v != nil {
		e.Oneof = newOneof(parent, v)
	}
	if v := desc.ContainingMessage(); v != nil {
		e.Extendee = newMessage(v)
	}

	if v := desc.Enum(); v != nil {
		e.Enum = newEnum(v)
	}
	if v := desc.Message(); v != nil {
		e.Message = newMessage(v)
	}
	return e
}

func newMethod(parent *protogen.Service, desc protoreflect.MethodDescriptor) *protogen.Method {
	m := &protogen.Method{
		Desc:   desc,
		GoName: string(desc.Name()),
		Parent: parent,
	}
	if v := desc.Input(); v != nil {
		m.Input = newMessage(v)
	}
	if v := desc.Output(); v != nil {
		m.Output = newMessage(v)
	}
	return m
}

func newService(desc protoreflect.ServiceDescriptor) *protogen.Service {
	s := &protogen.Service{
		Desc:    desc,
		GoName:  string(desc.Name()),
		Methods: []*protogen.Method{},
	}
	for i := 0; i < desc.Methods().Len(); i++ {
		s.Methods = append(s.Methods, newMethod(s, desc.Methods().Get(i)))
	}
	return s
}

func newFile(desc protoreflect.FileDescriptor) *protogen.File {
	f := &protogen.File{
		Desc:       desc,
		Enums:      []*protogen.Enum{},
		Messages:   []*protogen.Message{},
		Extensions: []*protogen.Extension{},
		Services:   []*protogen.Service{},
		Generate:   true,
	}
	for i := 0; i < desc.Enums().Len(); i++ {
		f.Enums = append(f.Enums, newEnum(desc.Enums().Get(i)))
	}
	for i := 0; i < desc.Messages().Len(); i++ {
		f.Messages = append(f.Messages, newMessage(desc.Messages().Get(i)))
	}
	for i := 0; i < desc.Extensions().Len(); i++ {
		f.Extensions = append(f.Extensions, newExtension(nil, desc.Extensions().Get(i)))
	}
	for i := 0; i < desc.Services().Len(); i++ {
		f.Services = append(f.Services, newService(desc.Services().Get(i)))
	}
	return f
}

func TestGenerate(t *testing.T) {
	tests := []struct {
		Name    string
		Desc    []protoreflect.FileDescriptor
		WantErr bool
	}{
		{
			Name:    "Basic",
			Desc:    []protoreflect.FileDescriptor{plugin.File_testdata_cmd_plugin_basic_proto},
			WantErr: false,
		},
		{
			Name:    "Advanced",
			Desc:    []protoreflect.FileDescriptor{plugin.File_testdata_cmd_plugin_advanced_proto},
			WantErr: false,
		},
		{
			Name:    "Crossref",
			Desc:    []protoreflect.FileDescriptor{plugin.File_testdata_cmd_plugin_crossref_proto},
			WantErr: false,
		},
		{
			Name:    "Fieldmask",
			Desc:    []protoreflect.FileDescriptor{plugin.File_testdata_cmd_plugin_fieldmask_proto},
			WantErr: false,
		},
		{
			Name:    "Reference",
			Desc:    []protoreflect.FileDescriptor{plugin.File_testdata_cmd_plugin_reference_proto},
			WantErr: false,
		},
		{
			Name:    "Method",
			Desc:    []protoreflect.FileDescriptor{plugin.File_testdata_cmd_plugin_method_proto},
			WantErr: false,
		},
		{
			Name:    "Error",
			Desc:    []protoreflect.FileDescriptor{plugin.File_testdata_cmd_plugin_error_proto},
			WantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			gs := []*protogen.File{}
			for i := 0; i < len(tt.Desc); i++ {
				gs = append(gs, newFile(tt.Desc[i]))
				gs[i].GeneratedFilenamePrefix = t.TempDir()
			}
			var i int32 = 0
			f, err := NewFile(&protogen.Plugin{Request: &pluginpb.CodeGeneratorRequest{CompilerVersion: &pluginpb.Version{
				Major: &i,
				Minor: &i,
				Patch: &i,
			}}, Files: gs}, gs[0], nil)
			if err != nil {
				t.Errorf("cannot build file: %v", err)
			} else if err = f.Generate(); err != nil != tt.WantErr {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}
