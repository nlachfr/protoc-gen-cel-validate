package gateway

import (
	"fmt"
	"plugin"

	"github.com/google/cel-go/cel"
	"github.com/nlachfr/protocel/validate"
)

type Plugin interface {
	Name() string
	Version() string
	BuildLibrary(args ...string) (cel.Library, error)
}

func LoadPlugins(cfgs ...*Configuration_Plugin) (cel.EnvOption, error) {
	plugins := map[*plugin.Plugin]*Configuration_Plugin{}
	for _, cfg := range cfgs {
		if p, err := plugin.Open(cfg.Path); err != nil {
			return nil, err
		} else {
			plugins[p] = cfg
		}
	}
	libs := &validate.Library{}
	for p, cfg := range plugins {
		symbol, err := p.Lookup("New")
		if err != nil {
			return nil, err
		} else if new, ok := symbol.(func() interface{}); !ok {
			return nil, fmt.Errorf("invalid symbol")
		} else if gwplugin, ok := new().(Plugin); !ok {
			return nil, fmt.Errorf("invalid plugin")
		} else {
			args := []string{cfg.Path}
			for k, v := range cfg.Args {
				args = append(args, "-"+k, v)
			}
			if lib, err := gwplugin.BuildLibrary(args...); err != nil {
				return nil, err
			} else if lib != nil {
				libs.EnvOpts = append(libs.EnvOpts, cel.Lib(lib))
			}
		}
	}
	return cel.Lib(libs), nil
}
