package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/Neakxs/protocel/gateway"
	"gopkg.in/yaml.v3"
)

var (
	config = flag.String("config", "config.yml", "global configuration file")
)

func loadConfig(config string, c *gateway.Configuration) error {
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
	flag.Parse()
	c := &gateway.Configuration{}
	if config != nil {
		if err := loadConfig(*config, c); err != nil {
			log.Fatalf("cannot load config: %v", err)
		}
	}
	linker, err := gateway.NewLinker(context.TODO(), c.Files)
	if err != nil {
		log.Fatal(err)
	}
	srvs := []*gateway.Server{}
	for _, srvCfg := range c.Servers {
		srv, err := gateway.NewServer(linker, srvCfg, c.Validate)
		if err != nil {
			log.Fatal(err)
		}
		srvs = append(srvs, srv)
	}
	fmt.Println("Starting server with following configuration :")
	b, _ := yaml.Marshal(c)
	fmt.Printf("\n\t%s\n", strings.ReplaceAll(string(b), "\n", "\n\t"))
	wg := &sync.WaitGroup{}
	errChan := make(chan error)
	doneChan := make(chan struct{})
	for _, srv := range srvs {
		wg.Add(1)
		go func(srv *gateway.Server) {
			defer wg.Done()
			if err := srv.ListenAndServe(); err != nil {
				errChan <- err
			}
		}(srv)
	}
	go func() {
		wg.Wait()
		doneChan <- struct{}{}
	}()
	select {
	case err := <-errChan:
		log.Fatal(err)
	case <-doneChan:
		return
	}
}
