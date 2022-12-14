package plugin

import (
	"flag"
	"os"

	"github.com/Neakxs/protocel/authorize"
	"github.com/Neakxs/protocel/options"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoregistry"
	"gopkg.in/yaml.v3"
)

var (
	config         = flag.String("config", "", "global configuration file")
	overrideStdlib = flag.Bool("override_stdlib", false, "override stdlib when protobuf names conflict with cel")
)

func LoadConfig(config string, c *authorize.AuthorizeOptions) error {
	if len(config) > 0 {
		b, err := os.ReadFile(config)
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(b, &c); err != nil {
			return err
		}
	}
	return nil
}

func Run() {
	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gen *protogen.Plugin) error {
		c := &authorize.AuthorizeOptions{}
		if config != nil {
			if err := LoadConfig(*config, c); err != nil {
				return err
			}
		}
		if overrideStdlib != nil && *overrideStdlib {
			if c.Options == nil {
				c.Options = &options.Options{}
			}
			c.Options.OverrideStdlib = true
		}
		var files protoregistry.Files
		for _, file := range gen.Files {
			if err := files.RegisterFile(file.Desc); err != nil {
				return err
			}
		}
		for _, file := range gen.Files {
			if !file.Generate {
				continue
			}
			if err := NewFile(gen, file, c).Generate(); err != nil {
				return err
			}
		}
		return nil
	})
}
