package gateway

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/linker"
)

type Linker struct {
	files []linker.File
}

func NewLinker(ctx context.Context, filesConfig *Configuration_Files) (*Linker, error) {
	if filesConfig == nil {
		return nil, nil
	}
	files := []string{}
	for _, path := range filesConfig.Sources {
		if matches, err := filepath.Glob(path); err != nil {
			return nil, err
		} else if matches == nil {
			return nil, fmt.Errorf("no matching file for pattern: %v", path)
		} else {
			files = append(files, matches...)
		}
	}
	imports := []string{}
	for _, path := range filesConfig.Imports {
		if matches, err := filepath.Glob(path); err != nil {
			return nil, err
		} else if matches == nil {
			return nil, fmt.Errorf("no matching file for pattern: %v", path)
		} else {
			imports = append(imports, matches...)
		}
	}
	res, err := (&protocompile.Compiler{
		Resolver: protocompile.WithStandardImports(
			&protocompile.SourceResolver{
				ImportPaths: imports,
			},
		),
	}).Compile(ctx, files...)
	if err != nil {
		return nil, err
	}
	return &Linker{files: res}, nil
}
