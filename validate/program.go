package validate

import (
	"context"
	"fmt"

	options "github.com/Neakxs/protocel/options"
	"github.com/google/cel-go/cel"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type Validater interface {
	Validate(ctx context.Context) error
	ValidateWithMask(ctx context.Context, fm *fieldmaskpb.FieldMask) error
}

type ValidateProgram struct {
	Id      string
	Expr    string
	Program cel.Program
}

type RuleValidater interface {
	Programs() []*ValidateProgram
}

type ruleValidater struct {
	programs []*ValidateProgram
}

func (v *ruleValidater) Programs() []*ValidateProgram { return v.programs }

func BuildRuleValidater(rule *Rule, envOpt cel.EnvOption) (RuleValidater, error) {
	envOpts := []cel.EnvOption{cel.Types(&fieldmaskpb.FieldMask{})}
	if envOpt != nil {
		envOpts = append(envOpts, envOpt)
	}
	validater := &ruleValidater{}
	if rule != nil {
		for _, rawProgram := range rule.Programs {
			customEnvOpts := envOpts
			if rule.Options != nil {
				if macros, err := options.BuildMacros(rule.Options, rawProgram.Expr, customEnvOpts); err != nil {
					return nil, fmt.Errorf("build macros error: %v", err)
				} else {
					customEnvOpts = append(customEnvOpts, cel.Macros(macros...))
				}
			}
			env, err := cel.NewCustomEnv(customEnvOpts...)
			if err != nil {
				return nil, fmt.Errorf("new env error: %w", err)
			}
			ast, issues := env.Compile(rawProgram.Expr)
			if issues != nil && issues.Err() != nil {
				return nil, fmt.Errorf("compile error: %w", issues.Err())
			}
			if !ast.OutputType().IsAssignableType(cel.BoolType) {
				return nil, fmt.Errorf("output type not bool")
			}
			pgr, err := env.Program(ast, cel.EvalOptions(cel.OptOptimize))
			if err != nil {
				return nil, fmt.Errorf("program error: %w", err)
			}
			validater.programs = append(validater.programs, &ValidateProgram{
				Id:      rawProgram.Id,
				Expr:    rawProgram.Expr,
				Program: pgr,
			})
		}
	}
	return validater, nil
}
