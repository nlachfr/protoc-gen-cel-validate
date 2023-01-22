package gateway

import (
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/Neakxs/protocel/validate"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func removePatternDuplicates(patterns []string) ([]string, error) {
	if len(patterns) <= 1 {
		return patterns, nil
	}
	np := []string{}
	sp := sort.StringSlice(patterns)
	sp.Sort()
	for sp.Len() > 1 {
		ok, err := filepath.Match(sp[0], sp[1])
		if err != nil {
			return nil, err
		}
		if ok {
			sp = append([]string{sp[0]}, sp[2:]...)
		} else {
			np = append(np, sp[0])
			sp = sp[1:]
		}
	}
	return append(np, sp[0]), nil
}

func NewServer(linker *Linker, serverCfg *Configuration_Server, opts *validate.Options) (*Server, error) {
	patterns := []string{}
	for p := range serverCfg.Upstreams {
		patterns = append(patterns, p)
	}
	patterns, err := removePatternDuplicates(patterns)
	if err != nil {
		return nil, err
	}
	handlerChan := make(chan struct {
		rpc     string
		handler http.Handler
	}, 4)
	errChan := make(chan error)
	doneChan := make(chan struct{})
	wg := &sync.WaitGroup{}
	for _, pattern := range patterns {
		upstream, err := NewUpstream(serverCfg.Upstreams[pattern])
		if err != nil {
			return nil, err
		}
		for _, file := range linker.files {
			manager, err := validate.NewManager(file, validate.WithFallbackOverloads(), validate.WithOptions(opts))
			if err != nil {
				return nil, err
			}
			for i := 0; i < file.Services().Len(); i++ {
				sd := file.Services().Get(i)
				if match, err := filepath.Match(pattern, string(sd.FullName())); err != nil {
					return nil, err
				} else if match {
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
				}
			}
		}
	}
	go func() {
		wg.Wait()
		doneChan <- struct{}{}
	}()
	mux := http.NewServeMux()
	for {
		select {
		case err := <-errChan:
			return nil, err
		case handler := <-handlerChan:
			fmt.Printf("\t%s: DONE\n", handler.rpc)
			mux.Handle(handler.rpc, handler.handler)
		case <-doneChan:
			return &Server{addrs: serverCfg.Listen, handler: h2c.NewHandler(mux, &http2.Server{})}, nil
		}
	}
}

type Server struct {
	addrs   []string
	handler http.Handler
}

func (s *Server) ListenAndServe() error {
	if len(s.addrs) == 0 {
		return fmt.Errorf("no binding address")
	}
	listeners := []net.Listener{}
	for _, bindAddr := range s.addrs {
		parts := strings.SplitN(bindAddr, "://", 2)
		if len(parts) == 1 {
			if listener, err := net.Listen("tcp", bindAddr); err != nil {
				return fmt.Errorf("cannot bind: %w", err)
			} else {
				listeners = append(listeners, listener)
			}
		} else {
			if listener, err := net.Listen(parts[0], parts[1]); err != nil {
				return fmt.Errorf("cannot bind: %w", err)
			} else {
				listeners = append(listeners, listener)
			}
		}
	}
	errChan := make(chan error)
	doneChan := make(chan struct{})
	wg := sync.WaitGroup{}
	for _, l := range listeners {
		wg.Add(1)
		go func(l net.Listener) {
			if err := http.Serve(l, s.handler); err != nil {
				errChan <- err
			}
			wg.Done()
		}(l)
	}
	go func() {
		wg.Wait()
		doneChan <- struct{}{}
	}()
	select {
	case err := <-errChan:
		return err
	case <-doneChan:
		return nil
	}
}
