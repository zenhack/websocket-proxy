package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"net"
	"os"
)

var (
	listen = flag.String("listen", ":8080", "Local address to listen on")
	url    = flag.String("url", "", "websocket url to connect to")
)

func main() {
	flag.Parse()
	if *url == "" {
		fmt.Println("Please provide a websocket url (via -url).")
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
			fmt.Fprintf(os.Stderr, "Error in accept: ", err)
			continue
		}
		ws, err := websocket.Dial(*url, "", "http://localhost/")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error connecting to websocket: ", err)
			conn.Close()
		}
		go io.Copy(conn, ws)
		go io.Copy(ws, conn)
	}
}
