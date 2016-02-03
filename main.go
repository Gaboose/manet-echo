package main

import (
	"fmt"
	manet "github.com/Gaboose/go-multiaddr-net"
	ma "github.com/jbenet/go-multiaddr"
	"io"
	"os"
	"strings"
)

func main() {
	ip := os.Getenv("OPENSHIFT_GO_IP")
	if ip == "" {
		ip = "127.0.0.1"
	}

	port := os.Getenv("OPENSHIFT_GO_PORT")
	if port == "" {
		port = "8000"
	}

	m, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s/ws/echo", ip, port))
	if err != nil {
		panic(err)
	}

	ln, err := manet.Listen(m)
	if err != nil {
		panic(err)
	}

	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go echo(c)
	}
}

func echo(rwc io.ReadWriteCloser) {
	for {
		buf := make([]byte, 256)
		n, err := rwc.Read(buf)
		if err != nil {
			rwc.Close()
			return
		}
		s := strings.TrimRight(string(buf[:n]), "\n")
		s = fmt.Sprintf("you're the \"%s\"\n", s)
		_, err = rwc.Write([]byte(s))
	}
}
