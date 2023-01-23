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
	imports := []string{}
	for _, path := range filesConfig.Imports {
		if matches, err := filepath.Glob(path); err != nil {
			return nil, err
		} else if matches == nil {
			return nil, fmt.Errorf("no matching file for import pattern: %v", path)
		} else {
			imports = append(imports, matches...)
		}
	}
	if len(imports) == 0 {
		imports = append(imports, ".")
	}
	files := []string{}
	for _, path := range filesConfig.Sources {
		var matches []string
		for _, imp := range filesConfig.Imports {
			if impMatches, err := filepath.Glob(filepath.Join(imp, path)); err != nil {
				return nil, err
			} else if impMatches != nil {
				for _, v := range impMatches {
					detectedPath, _ := filepath.Rel(imp, v)
					matches = append(matches, detectedPath)
				}
				break
			}
		}
		if len(matches) == 0 {
			return nil, fmt.Errorf("no matching file for pattern: %v", path)
		}
		files = append(files, matches...)
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
