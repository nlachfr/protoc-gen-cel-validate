package gateway

import (
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nlachfr/protocel/validate"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func findServiceUpstream(sd protoreflect.ServiceDescriptor, upstreams map[string]*Upstream) (upstream *Upstream, err error) {
	name := string(sd.FullName())
	patternParts := []string{}
	for p, u := range upstreams {
		if matched, err := filepath.Match(p, name); err != nil {
			return nil, err
		} else if matched {
			selectUpstream := false
			newPatternParts := strings.Split(p, ".")
			if len(patternParts) < len(newPatternParts) {
				selectUpstream = true
			} else if len(patternParts) == len(newPatternParts) {
				for i := 0; i < len(patternParts); i++ {
					lp := patternParts[i]
					rp := newPatternParts[i]
					if strings.Contains(lp, "*") {
						if !strings.Contains(rp, "*") {
							selectUpstream = true
							break
						}
					} else if strings.Contains(rp, "*") {
						break
					}
					if rp[i] > lp[i] {
						selectUpstream = true
						break
					}
				}
			}
			if selectUpstream {
				patternParts = newPatternParts
				upstream = u
			}
		}
	}
	return upstream, nil
}

func NewServer(linker *Linker, serverCfg *Configuration_Server, opts *validate.Options) (*Server, error) {
	upstreams := map[string]*Upstream{}
	for pattern, upstreamCfg := range serverCfg.Upstreams {
		upstream, err := NewUpstream(upstreamCfg)
		if err != nil {
			return nil, err
		}
		upstreams[pattern] = upstream
	}
	handlerChan := make(chan struct {
		rpc     string
		handler http.Handler
	}, 4)
	errChan := make(chan error)
	doneChan := make(chan struct{})
	wg := &sync.WaitGroup{}
	for _, file := range linker.files {
		manager, err := validate.NewManager(file, validate.WithFallbackOverloads(), validate.WithOptions(opts))
		if err != nil {
			return nil, err
		}
		for i := 0; i < file.Services().Len(); i++ {
			sd := file.Services().Get(i)
			if upstream, err := findServiceUpstream(sd, upstreams); err != nil {
				return nil, err
			} else if upstream != nil {
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
