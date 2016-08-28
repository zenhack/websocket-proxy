package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"net"
	"net/url"
	"os"
)

var (
	listen   = flag.String("listen", ":8080", "Local address to listen on")
	rawUrl   = flag.String("url", "", "websocket url to connect to")
	protocol = flag.String("protocol", "", "Value of Sec-WebSocket-Protocol")
)

func copyClose(w io.Writer, r io.ReadCloser) {
	io.Copy(w, r)
	r.Close()
}

func main() {
	flag.Parse()
	if *rawUrl == "" {
		fmt.Println("Please provide a websocket url (via -url).")
		os.Exit(1)
	}
	parsedUrl, err := url.Parse(*rawUrl)
	if err != nil {
		fmt.Println("Invalid URL: ", err)
		os.Exit(1)
	}

	ln, err := net.Listen("tcp", *listen)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening listen socket: ", err)
		os.Exit(1)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error in accept: ", err)
			continue
		}
		ws, err := websocket.Dial(*rawUrl, *protocol, parsedUrl.Host)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error connecting to websocket: ", err)
			conn.Close()
			continue
		}
		go copyClose(conn, ws)
		go copyClose(ws, conn)
	}
}
