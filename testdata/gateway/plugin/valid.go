package main

import (
	"github.com/google/cel-go/cel"
	"github.com/nlachfr/protocel/validate"
)

func New() interface{} {
	return &Plugin{}
}

type Plugin struct{}

func (*Plugin) Name() string                                     { return "Valid" }
func (*Plugin) Version() string                                  { return "v0.0.0" }
func (*Plugin) BuildLibrary(args ...string) (cel.Library, error) { return &validate.Library{}, nil }
