package plugin

import (
	"flag"
	"os"

	"github.com/Neakxs/protocel/options"
	"github.com/Neakxs/protocel/validate"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoregistry"
	"gopkg.in/yaml.v3"
)

var (
	config                           = flag.String("config", "", "global configuration file")
	stdlibOverridingEnabled          = flag.Bool("stdlib_overriding_enabled", false, "override stdlib when protobuf names conflict with cel")
	requiredSupportDisabled          = flag.Bool("required_support_disabled", false, "disable google.protobuf.field_behavior.REQUIRED support")
	resourceReferenceSupportDisabled = flag.Bool("resource_reference_support_disabled", false, "disable google.protobuf.resource_reference rules generation")
)

func LoadConfig(config string, c *validate.ValidateOptions) error {
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
		c := &validate.ValidateOptions{}
		if config != nil {
			if err := LoadConfig(*config, c); err != nil {
				return err
			}
		}
		flag.Visit(func(f *flag.Flag) {
			switch f.Name {
			case "stdlib_overriding_enabled":
				if c.Options == nil {
					c.Options = &options.Options{}
				}
				c.Options.StdlibOverridingEnabled = *stdlibOverridingEnabled
			case "required_support_disabled":
				c.RequiredSupportDisabled = *requiredSupportDisabled
			case "resource_reference_support_disabled":
				c.ResourceReferenceSupportDisabled = *resourceReferenceSupportDisabled
			}
		})
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
			if f, err := NewFile(gen, file, c); err != nil {
				return err
			} else if err = f.Generate(); err != nil {
				return err
			}
		}
		return nil
	})
}
