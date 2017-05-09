package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"zenhack.net/go/socks5"
)

var (
	laddr  = flag.String("laddr", ":1080", "Address to listen on")
	config = flag.String("config", os.Getenv("HOME")+"/.ws-multi-cfg.json",
		"Path to config file.")
)

type Config struct {
	Endpoints map[string]Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Url      string `json:"url"`
	Protocol string `json:"protocol,omitempty"`
	Origin   string `json:"origin,omitempty"`
}

func (c *Config) ReadFrom(r io.Reader) (int64, error) {
	dec := json.NewDecoder(r)
	err := dec.Decode(c)
	// XXX: not returning the right length
	if err != nil {
		return 0, err
	}
	for k, v := range c.Endpoints {
		parsedUrl, err := url.Parse(v.Url)
		if err != nil {
			return 0, err
		}
		if v.Origin == "" {
			v.Origin = parsedUrl.Scheme + "://" + parsedUrl.Host
		}
		c.Endpoints[k] = v
	}
	return 0, nil
}

func (c *Config) Dial(net, addr string) (net.Conn, error) {
	if net != "tcp" {
		return nil, fmt.Errorf("Protocol %s not supported (tcp only)", net)
	}
	endpoint, ok := c.Endpoints[addr]
	if !ok {
		return nil, fmt.Errorf("Endpoint not found: %q", addr)
	}
	return websocket.Dial(endpoint.Url, endpoint.Protocol, endpoint.Origin)
}

func chkfatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	cfg := new(Config)
	flag.Parse()
	file, err := os.Open(*config)
	chkfatal(err)
	_, err = cfg.ReadFrom(file)
	chkfatal(err)
	file.Close()
	chkfatal(socks5.ListenAndServe(cfg, *laddr))
}
