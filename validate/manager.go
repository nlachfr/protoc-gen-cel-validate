package validate

import (
	"fmt"
	"path/filepath"
	sync "sync"

	"github.com/Neakxs/protocel/options"
	"github.com/google/cel-go/cel"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var registry = &managerRegistry{
	registry: &sync.Map{}, // map[string]map[Manager]bool
}

func LoadLibrary(pattern string, lib cel.Library) error {
	return registry.LoadLibrary(pattern, lib)
}

type managerRegistry struct {
	registry *sync.Map
}

func (r *managerRegistry) LoadLibrary(pattern string, lib cel.Library) error {
	if lib != nil {
		var err error
		r.registry.Range(func(key, value any) bool {
			var registryErr error
			if ok, _ := filepath.Match(pattern, key.(string)); ok {
				value.(*sync.Map).Range(func(key, value any) bool {
					if e := key.(*Manager).LoadLibrary(lib); e != nil {
						registryErr = e
					}
					return false
				})
			}
			if registryErr != nil {
				err = registryErr
				return false
			}
			return true
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *managerRegistry) Register(m *Manager) error {
	if m == nil {
		return fmt.Errorf("nil manager")
	}
	rawScopedMap, _ := r.registry.LoadOrStore(string(m.file.Package()), &sync.Map{})
	scopedMap := rawScopedMap.(*sync.Map)
	if _, loaded := scopedMap.LoadOrStore(m, true); loaded {
		return fmt.Errorf("manager already registered")
	}
	return nil
}

type ManagerOption interface {
	apply(b *builder)
}

type managerOption func(b *builder)

func (opt managerOption) apply(b *builder) { opt(b) }

func WithFallbackOverloads() ManagerOption {
	return managerOption(func(b *builder) {
		b.ob = &fallbackOverloadBuilder{b}
	})
}

func WithOptions(optsList ...*Options) ManagerOption {
	return managerOption(func(b *builder) {
		opts := &Options{}
		if b.opts != nil {
			opts = proto.Clone(b.opts).(*Options)
		}
		for _, o := range optsList {
			if o != nil {
				proto.Merge(opts, o)
			}
		}
		b.opts = opts
	})
}

func NewManager(file protoreflect.FileDescriptor, opts ...ManagerOption) (*Manager, error) {
	if file == nil {
		return nil, fmt.Errorf("nil file descriptor")
	}
	m := &Manager{
		file:  file,
		onces: &sync.Map{},
		serviceValidaters: make(map[string]struct {
			rv  ServiceRuleValidater
			err error
		}),
		messageValidaters: make(map[string]struct {
			rv  MessageRuleValidater
			err error
		}),
		b: newBuilder(),
	}
	for _, opt := range opts {
		opt.apply(m.b)
	}
	return m, registry.Register(m)
}

type Manager struct {
	file              protoreflect.FileDescriptor
	onces             *sync.Map
	serviceValidaters map[string]struct {
		rv  ServiceRuleValidater
		err error
	}
	messageValidaters map[string]struct {
		rv  MessageRuleValidater
		err error
	}

	b *builder
}

func (m *Manager) LoadLibrary(lib cel.Library) error {
	if len(m.serviceValidaters) > 0 || len(m.messageValidaters) > 0 {
		return fmt.Errorf("cannot load library: manager already used")
	} else if lib == nil {
		return nil
	}
	if m.b.envOpt != nil {
		cel.Lib(&options.Library{EnvOpts: []cel.EnvOption{m.b.envOpt, cel.Lib(lib)}})
	} else {
		m.b.envOpt = cel.Lib(lib)
	}
	return nil
}

func (m *Manager) BuildValidaters() error {
	for i := 0; i < m.file.Services().Len(); i++ {
		if _, err := m.GetServiceRuleValidater(m.file.Services().Get(i)); err != nil {
			return err
		}
	}
	for i := 0; i < m.file.Messages().Len(); i++ {
		if _, err := m.GetMessageRuleValidater(m.file.Messages().Get(i)); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) GetServiceRuleValidater(desc protoreflect.ServiceDescriptor) (ServiceRuleValidater, error) {
	key := string(desc.FullName())
	once, _ := m.onces.LoadOrStore(key, &sync.Once{})
	once.(*sync.Once).Do(func() {
		rv, err := m.b.BuildServiceRuleValidater(desc)
		m.serviceValidaters[key] = struct {
			rv  ServiceRuleValidater
			err error
		}{rv: rv, err: err}
	})
	return m.serviceValidaters[key].rv, m.serviceValidaters[key].err
}

func (m *Manager) GetMessageRuleValidater(desc protoreflect.MessageDescriptor) (MessageRuleValidater, error) {
	key := string(desc.FullName())
	once, _ := m.onces.LoadOrStore(key, &sync.Once{})
	once.(*sync.Once).Do(func() {
		rv, err := m.b.BuildMessageRuleValidater(desc)
		m.messageValidaters[key] = struct {
			rv  MessageRuleValidater
			err error
		}{rv: rv, err: err}
	})
	return m.messageValidaters[key].rv, m.serviceValidaters[key].err
}
