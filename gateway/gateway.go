package gateway

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/linker"
)

func CompileFiles(ctx context.Context, filesConfig *Configuration_Files) ([]linker.File, error) {
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
	return (&protocompile.Compiler{
		Resolver: protocompile.WithStandardImports(
			&protocompile.SourceResolver{
				ImportPaths: imports,
			},
		),
	}).Compile(ctx, files...)
}

func NewGateway(ctx context.Context, config *Configuration) (http.Handler, error) {
	mux := http.NewServeMux()
	if config == nil {
		return nil, fmt.Errorf("nil config")
	} else if files, err := CompileFiles(ctx, config.Files); err != nil {
		return nil, err
	} else {
		for _, file := range files {
			for i := 0; i < file.Services().Len(); i++ {

			}
		}
	}
	return mux, nil
}
