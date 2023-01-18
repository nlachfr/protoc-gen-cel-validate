package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
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
	listeners := []net.Listener{}
	if c.Serve != nil {
		for _, bindAddr := range c.Serve.Bind {
			parts := strings.SplitN(bindAddr, "://", 2)
			if len(parts) == 1 {
				if listener, err := net.Listen("tcp", bindAddr); err != nil {
					log.Fatalf("cannot bind: %v", err)
				} else {
					listeners = append(listeners, listener)
				}
			} else {
				if listener, err := net.Listen(parts[0], parts[1]); err != nil {
					log.Fatalf("cannot bind: %v", err)
				} else {
					listeners = append(listeners, listener)
				}
			}

		}
	} else {
		log.Fatal("no binding address")
	}
	handler, err := gateway.NewGateway(context.Background(), c)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Starting server with following configuration :")
	b, _ := yaml.Marshal(c)
	fmt.Printf("\n\t%s\n", strings.ReplaceAll(string(b), "\n", "\n\t"))
	wg := sync.WaitGroup{}
	for _, l := range listeners {
		wg.Add(1)
		go func(l net.Listener) {
			if err := http.Serve(l, handler); err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}(l)
	}
	wg.Wait()
}
