package validate

import (
	"fmt"

	"github.com/Neakxs/protocel/options"
	"github.com/google/cel-go/cel"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/rpc/context/attribute_context"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func WithFallbackOverloads() BuildOption {
	return buildOption(func(b *builder) *builder {
		return &builder{
			ob:     &fallbackOverloadBuilder{b},
			opts:   b.opts,
			envOpt: b.envOpt,
		}
	})
}

func WithDescriptors(descs ...protoreflect.Descriptor) BuildOption {
	return buildOption(func(b *builder) *builder {
		lib := &options.Library{}
		for _, desc := range descs {
			fileDesc := desc.ParentFile()
			lib.EnvOpts = append(lib.EnvOpts, cel.TypeDescs(fileDesc))
			for i := 0; i < fileDesc.Imports().Len(); i++ {
				lib.EnvOpts = append(lib.EnvOpts, cel.TypeDescs(fileDesc.Imports().Get(i)))
			}
		}
		if b.envOpt != nil {
			lib.EnvOpts = append(lib.EnvOpts, b.envOpt)
		}
		newBuilder := &builder{
			ob:     b.ob,
			opts:   b.opts,
			envOpt: cel.Lib(lib),
		}
		if _, ok := b.ob.(*fallbackOverloadBuilder); ok {
			newBuilder.ob = &fallbackOverloadBuilder{newBuilder}
		}
		return newBuilder
	})
}

func WithEnvOptions(envOpts ...cel.EnvOption) BuildOption {
	return buildOption(func(b *builder) *builder {
		lib := &options.Library{}
		for _, envOpt := range envOpts {
			if envOpt != nil {
				lib.EnvOpts = append(lib.EnvOpts, envOpt)
			}
		}
		if b.envOpt != nil {
			lib.EnvOpts = append(lib.EnvOpts, b.envOpt)
		}
		newBuilder := &builder{
			ob:     b.ob,
			opts:   b.opts,
			envOpt: cel.Lib(lib),
		}
		if _, ok := b.ob.(*fallbackOverloadBuilder); ok {
			newBuilder.ob = &fallbackOverloadBuilder{newBuilder}
		}
		return newBuilder
	})
}

func WithOptions(optsList ...*Options) BuildOption {
	return buildOption(func(b *builder) *builder {
		opts := &Options{}
		if b.opts != nil {
			opts = proto.Clone(b.opts).(*Options)
		}
		for _, o := range optsList {
			if o != nil {
				proto.Merge(opts, o)
			}
		}
		newBuilder := &builder{
			ob:     b.ob,
			opts:   opts,
			envOpt: b.envOpt,
		}
		if _, ok := b.ob.(*fallbackOverloadBuilder); ok {
			newBuilder.ob = &fallbackOverloadBuilder{newBuilder}
		}
		return newBuilder
	})
}

type BuildOption interface {
	apply(b *builder) *builder
}

type buildOption func(b *builder) *builder

func (opt buildOption) apply(b *builder) *builder { return opt(b) }

type Builder interface {
	WithBuildOptions(opts ...BuildOption) Builder
	BuildServiceRuleValidater(desc protoreflect.ServiceDescriptor) (ServiceRuleValidater, error)
	BuildMessageRuleValidater(desc protoreflect.MessageDescriptor) (MessageRuleValidater, error)
}

type builder struct {
	ob     overloadBuilder
	opts   *Options
	envOpt cel.EnvOption
}

func NewBuilder(opts ...BuildOption) *builder {
	b := &builder{ob: &defaultOverloadBuilder{}}
	for _, opt := range opts {
		b = opt.apply(b)
	}
	return b
}

func (b *builder) WithBuildOptions(opts ...BuildOption) Builder {
	nb := b
	for _, opt := range opts {
		nb = opt.apply(b)
	}
	return nb
}

func (b *builder) BuildServiceRuleValidater(desc protoreflect.ServiceDescriptor) (ServiceRuleValidater, error) {
	serviceRule := &ServiceRule{
		Options: &options.Options{},
	}
	if b.opts != nil && b.opts.Rule != nil {
		proto.Merge(serviceRule.Options, b.opts.Rule.Options)
		if sr, ok := b.opts.Rule.ServiceRules[string(desc.FullName())]; ok {
			proto.Merge(serviceRule, sr)
		}
	}
	if fr := GetExtension(desc.ParentFile().Options(), E_File).(*FileRule); fr != nil {
		proto.Merge(serviceRule.Options, fr.Options)
		if sr, ok := fr.ServiceRules[string(desc.FullName())]; ok {
			proto.Merge(serviceRule, sr)
		}
	}
	if sr := GetExtension(desc.Options(), E_Service).(*ServiceRule); sr != nil {
		proto.Merge(serviceRule, sr)
	}
	rule := &Rule{
		Options: &options.Options{},
	}
	proto.Merge(rule.Options, serviceRule.Options)
	if serviceRule.Rule != nil {
		proto.Merge(rule, serviceRule.Rule)
	}
	lib := &options.Library{}
	if b.envOpt != nil {
		lib.EnvOpts = append(lib.EnvOpts, b.envOpt)
	}
	lib.EnvOpts = append(lib.EnvOpts,
		cel.TypeDescs(attribute_context.File_google_rpc_context_attribute_context_proto),
		cel.Variable("attribute_context", cel.ObjectType(string((&attribute_context.AttributeContext{}).ProtoReflect().Descriptor().FullName()))),
	)
	methodDescs := map[string]protoreflect.MethodDescriptor{}
	methodRulesValidaters := map[string]MethodRuleValidater{}
	for i := 0; i < desc.Methods().Len(); i++ {
		methodDesc := desc.Methods().Get(i)
		if methodValidater, err := b.buildMethodRuleValidater(serviceRule, methodDesc, cel.Lib(lib)); err != nil {
			return nil, err
		} else {
			methodDescs[string(methodDesc.FullName())] = methodDesc
			methodRulesValidaters[string(methodDesc.FullName())] = methodValidater
		}
	}
	var ruleValidater RuleValidater
	if len(rule.Programs) > 0 {
		lib.EnvOpts = append(lib.EnvOpts, options.BuildEnvOption(rule.Options))
		if rv, err := BuildRuleValidater(rule, cel.Lib(lib)); err != nil {
			return nil, err
		} else {
			ruleValidater = rv
		}
	}
	return &serviceRuleValidater{ruleValidater: ruleValidater, methodRulesValidaters: methodRulesValidaters}, nil
}

func (b *builder) buildMethodRuleValidater(serviceRule *ServiceRule, desc protoreflect.MethodDescriptor, envOpt cel.EnvOption) (MethodRuleValidater, error) {
	if desc == nil {
		return nil, fmt.Errorf("nil desc")
	}
	rule := &Rule{
		Options: &options.Options{},
	}
	if b.opts != nil && b.opts.Rule != nil {
		proto.Merge(rule.Options, b.opts.Rule.Options)
		if sr, ok := b.opts.Rule.ServiceRules[string(desc.Parent().FullName())]; ok {
			proto.Merge(rule.Options, sr.Options)
			if mr, ok := sr.MethodRules[string(desc.Name())]; ok {
				proto.Merge(rule, mr.Rule)
			}
		}
	}
	if fr := GetExtension(desc.ParentFile().Options(), E_File).(*FileRule); fr != nil {
		proto.Merge(rule.Options, fr.Options)
		if sr, ok := fr.ServiceRules[string(desc.Parent().FullName())]; ok {
			proto.Merge(rule.Options, sr.Options)
			if mr, ok := sr.MethodRules[string(desc.Name())]; ok {
				proto.Merge(rule, mr.Rule)
			}
		}
	}
	if serviceRule != nil {
		proto.Merge(rule.Options, serviceRule.Options)
		if mr, ok := serviceRule.MethodRules[string(desc.Name())]; ok {
			proto.Merge(rule, mr.Rule)
		}
	}
	if mr := GetExtension(desc.Options(), E_Method).(*MethodRule); mr != nil {
		proto.Merge(rule, mr.Rule)
	}
	lib := &options.Library{}
	if envOpt != nil {
		lib.EnvOpts = append(lib.EnvOpts, envOpt)
	}
	lib.EnvOpts = append(lib.EnvOpts,
		cel.TypeDescs(desc.Input().ParentFile()),
		cel.Variable("request", cel.ObjectType(string(desc.Input().FullName()))),
	)
	lib.EnvOpts = append(lib.EnvOpts, options.BuildEnvOption(rule.Options))
	lib.EnvOpts = append(lib.EnvOpts, b.ob.buildOverloads(desc.Input())...)
	if len(rule.Programs) > 0 {
		if rv, err := BuildRuleValidater(rule, cel.Lib(lib)); err != nil {
			return nil, err
		} else {
			return &methodRuleValidater{validater: rv}, nil
		}
	}
	return nil, nil
}

func (b *builder) BuildMessageRuleValidater(desc protoreflect.MessageDescriptor) (MessageRuleValidater, error) {
	messageRule := &MessageRule{
		Options: &options.Options{},
	}
	if b.opts != nil && b.opts.Rule != nil {
		proto.Merge(messageRule.Options, b.opts.Rule.Options)
		if mr, ok := b.opts.Rule.MessageRules[string(desc.FullName())]; ok {
			proto.Merge(messageRule, mr)
		}
	}
	if fr := GetExtension(desc.ParentFile().Options(), E_File).(*FileRule); fr != nil {
		proto.Merge(messageRule.Options, fr.Options)
		if mr, ok := fr.MessageRules[string(desc.FullName())]; ok {
			proto.Merge(messageRule, mr)
		}
	}
	if mr := GetExtension(desc.Options(), E_Message).(*MessageRule); mr != nil {
		proto.Merge(messageRule, mr)
	}
	rule := &Rule{
		Options: &options.Options{},
	}
	proto.Merge(rule.Options, messageRule.Options)
	if messageRule.Rule != nil {
		proto.Merge(rule, messageRule.Rule)
	}
	lib := &options.Library{}
	if b.envOpt != nil {
		lib.EnvOpts = append(lib.EnvOpts, b.envOpt)
	}
	lib.EnvOpts = append(lib.EnvOpts, cel.DeclareContextProto(desc))
	lib.EnvOpts = append(lib.EnvOpts, b.ob.buildOverloads(desc)...)
	fieldRulesValidaters := map[string]FieldRuleValidater{}
	for i := 0; i < desc.Fields().Len(); i++ {
		fieldDesc := desc.Fields().Get(i)
		if fieldValidater, err := b.buildFieldRuleValidater(messageRule, fieldDesc, cel.Lib(lib)); err != nil {
			return nil, err
		} else {
			fieldRulesValidaters[string(fieldDesc.Name())] = fieldValidater
		}
	}
	var ruleValidater RuleValidater
	if len(rule.Programs) > 0 {
		lib.EnvOpts = append(lib.EnvOpts, options.BuildEnvOption(rule.Options, desc))
		if rv, err := BuildRuleValidater(rule, cel.Lib(lib)); err != nil {
			return nil, err
		} else {
			ruleValidater = rv
		}
	}
	return &messageRuleValidater{ruleValidater: ruleValidater, fieldRulesValidaters: fieldRulesValidaters}, nil
}

func (b *builder) buildFieldRuleValidater(messageRule *MessageRule, desc protoreflect.FieldDescriptor, envOpt cel.EnvOption) (FieldRuleValidater, error) {
	if desc == nil {
		return nil, fmt.Errorf("nil desc")
	}
	rule := &Rule{
		Options: &options.Options{},
	}
	if b.opts != nil && b.opts.Rule != nil {
		proto.Merge(rule.Options, b.opts.Rule.Options)
		if mr, ok := b.opts.Rule.MessageRules[string(desc.Parent().FullName())]; ok {
			proto.Merge(rule.Options, mr.Options)
			if fr, ok := mr.FieldRules[string(desc.Name())]; ok {
				proto.Merge(rule, fr.Rule)
			}
		}
	}
	if fr := GetExtension(desc.ParentFile().Options(), E_File).(*FileRule); fr != nil {
		proto.Merge(rule.Options, fr.Options)
		if mr, ok := fr.MessageRules[string(desc.Parent().FullName())]; ok {
			proto.Merge(rule.Options, mr.Options)
			if fr, ok := mr.FieldRules[string(desc.Name())]; ok {
				proto.Merge(rule, fr.Rule)
			}
		}
	}
	if messageRule != nil {
		proto.Merge(rule.Options, messageRule.Options)
		if fr, ok := messageRule.FieldRules[string(desc.Name())]; ok {
			proto.Merge(rule, fr.Rule)
		}
	}
	if fr := GetExtension(desc.Options(), E_Field).(*FieldRule); fr != nil {
		proto.Merge(rule, fr.Rule)
	}
	lib := &options.Library{}
	if envOpt != nil {
		lib.EnvOpts = append(lib.EnvOpts, envOpt)
	}
	lib.EnvOpts = append(lib.EnvOpts, options.BuildEnvOption(rule.Options, desc.ContainingMessage()))
	envOpt = cel.Lib(lib)
	resourceReferenceMap := GenerateResourceTypePatternMapping(desc)
	if b.opts == nil || !b.opts.ResourceReferenceSupportDisabled {
		if resourceReference := proto.GetExtension(desc.Options(), annotations.E_ResourceReference).(*annotations.ResourceReference); resourceReference != nil {
			var ref string
			if resourceReference.Type != "" {
				if resourceReference.ChildType != "" {
					return nil, fmt.Errorf(`resource reference error: type and child_type are defined`)
				} else if resourceReference.Type != "*" {
					ref = resourceReference.Type
				}
			} else if resourceReference.ChildType != "" {
				ref = resourceReference.ChildType
			}
			if regexp, ok := resourceReferenceMap[ref]; ok {
				expr := ""
				if desc.IsList() {
					expr = fmt.Sprintf(`%s.all(s, s.matches("%s"))`, desc.TextName(), regexp)
				} else if desc.Kind() == protoreflect.StringKind {
					expr = fmt.Sprintf(`%s.matches("%s")`, desc.TextName(), regexp)
				}
				if expr != "" {
					rule.Programs = append(rule.Programs, &Rule_Program{
						Id:   ref,
						Expr: expr,
					})
				}
			} else {
				return nil, fmt.Errorf(`cannot find type "%s"`, ref)
			}
		}
	}
	required := false
	if b.opts != nil && !b.opts.RequiredSupportDisabled {
		for _, behavior := range proto.GetExtension(desc.Options(), annotations.E_FieldBehavior).([]annotations.FieldBehavior) {
			if behavior == annotations.FieldBehavior_REQUIRED {
				required = true
			}
		}
	}
	var ruleValidater RuleValidater
	if len(rule.Programs) > 0 {
		if rv, err := BuildRuleValidater(rule, envOpt); err != nil {
			return nil, err
		} else {
			ruleValidater = rv
		}
	}
	return &fieldRuleValidater{validater: ruleValidater, required: required}, nil
}
