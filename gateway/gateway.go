package gateway

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"sync"

	"github.com/Neakxs/protocel/validate"
	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/linker"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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

func BuildUpstreamClient(cfg *Configuration_Upstream) (*http.Client, error) {
	if cfg == nil || cfg.Address == "" {
		return nil, fmt.Errorf("")
	}
	addrUrl, err := url.Parse(cfg.Address)
	if err != nil {
		return nil, err
	}
	var upstreamNet, upstreamAddr string = "tcp", addrUrl.Host
	httpTransport := &http2.Transport{
		AllowHTTP: true,
	}
	if cfg.Server != "" {
		if srvUrl, err := url.Parse(cfg.Server); err != nil {
			return nil, err
		} else {
			switch srvUrl.Scheme {
			case "tcp", "tcp4", "tcp6", "udp", "udp4", "udp6", "unix", "unixgram", "unixpacket":
				upstreamNet = srvUrl.Scheme
				upstreamAddr = srvUrl.Host
			default:
				return nil, fmt.Errorf(`unknown scheme: "%s"`, srvUrl.Scheme)
			}
		}
	}
	switch addrUrl.Scheme {
	case "http":
		httpTransport.DialTLSContext = func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.DefaultResolver.Dial(ctx, upstreamNet, upstreamAddr)
		}
	case "https":
		httpTransport.DialTLSContext = func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
			return tls.Dial(upstreamNet, upstreamAddr, httpTransport.TLSClientConfig)
		}
	default:
		return nil, fmt.Errorf(`unknown scheme: "%s"`, addrUrl.Scheme)
	}
	return &http.Client{
		Transport: httpTransport,
	}, nil
}

func NewGateway(ctx context.Context, config *Configuration) (http.Handler, error) {
	mux := http.NewServeMux()
	if config == nil {
		return nil, fmt.Errorf("nil config")
	} else if files, err := CompileFiles(ctx, config.Files); err != nil {
		return nil, err
	} else {
		if config.Upstreams == nil {
			config.Upstreams = map[string]*Configuration_Upstream{}
		}
		handlerChan := make(chan struct {
			rpc     string
			handler http.Handler
		}, 4)
		errChan := make(chan error)
		doneChan := make(chan struct{})
		wg := &sync.WaitGroup{}
		fmt.Println("Compiling validation rules...")
		for _, file := range files {
			manager, err := validate.NewManager(file, validate.WithFallbackOverloads(), validate.WithOptions(config.Validate))
			if err != nil {
				return nil, fmt.Errorf("cannot build manager for '%s': %w", file.Path(), err)
			}
			for i := 0; i < file.Services().Len(); i++ {
				sd := file.Services().Get(i)
				if upstreamConfig, ok := config.Upstreams[string(sd.FullName())]; ok {
					upstream, err := NewUpstream(upstreamConfig)
					if err != nil {
						return nil, err
					}
					wg.Add(1)
					go func() {
						defer wg.Done()
						svr, err := manager.GetServiceRuleValidater(sd)
						if err != nil {
							errChan <- err
						}
						rpc, handler := NewServiceHandler(sd, svr, upstream)
						handlerChan <- struct {
							rpc     string
							handler http.Handler
						}{
							rpc: rpc, handler: handler,
						}
					}()
				} else {
					fmt.Printf("Found no matching upstream for service: %s\n", sd.FullName())
				}
			}
		}
		go func() {
			wg.Wait()
			doneChan <- struct{}{}
		}()
		for {
			select {
			case err := <-errChan:
				return nil, err
			case handler := <-handlerChan:
				fmt.Printf("\t%s: DONE\n", handler.rpc)
				mux.Handle(handler.rpc, handler.handler)
			case <-doneChan:
				return h2c.NewHandler(mux, &http2.Server{}), nil
			}
		}
	}
}
