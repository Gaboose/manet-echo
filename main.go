package main

import (
	"fmt"
	manet "github.com/Gaboose/go-multiaddr-net"
	ma "github.com/jbenet/go-multiaddr"
	"io"
	lg "log"
	"os"
	"strings"
)

var log = lg.New(os.Stdout, "", lg.LstdFlags|lg.Lshortfile)

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
	go serve(m, func(s string) string {
		return fmt.Sprintf("you're the \"%s\"\n", s)
	})

	m, err = ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s/ws/notecho", ip, port))
	if err != nil {
		panic(err)
	}
	serve(m, func(s string) string {
		return "not speaking to you\n"
	})
}

func serve(m ma.Multiaddr, handler func(string) string) {
	ln, err := manet.Listen(m)
	if err != nil {
		panic(err)
	}

	for {
		c, err := ln.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		go echo(c, handler)
	}
}

func echo(rwc io.ReadWriteCloser, handler func(string) string) {
	for {
		buf := make([]byte, 256)
		n, err := rwc.Read(buf)
		if err != nil {
			log.Println(err)
			rwc.Close()
			return
		}
		s := strings.TrimRight(string(buf[:n]), "\n")
		s = handler(s)
		_, err = rwc.Write([]byte(s))
	}
}
