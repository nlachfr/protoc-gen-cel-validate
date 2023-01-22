package main

import (
	"flag"
	"os"

	"github.com/nlachfr/protocel/cmd/protoc-gen-go-cel-validate/internal/plugin"
	"github.com/nlachfr/protocel/options"
	"github.com/nlachfr/protocel/validate"
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

func loadConfig(config string, c *validate.Options) error {
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

func main() {
	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gen *protogen.Plugin) error {
		c := &validate.Options{}
		if config != nil {
			if err := loadConfig(*config, c); err != nil {
				return err
			}
		}
		flag.Visit(func(f *flag.Flag) {
			switch f.Name {
			case "stdlib_overriding_enabled":
				if c.Rule == nil {
					c.Rule.Options = &options.Options{}
				}
				c.Rule.Options.StdlibOverridingEnabled = *stdlibOverridingEnabled
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
			if f, err := plugin.NewFile(gen, file, c); err != nil {
				return err
			} else if err = f.Generate(); err != nil {
				return err
			}
		}
		return nil
	})
}
