package main

import (
	"net"
	"testing"
)

func Testmain(t *testing.T) {
	go func() {
		listen, _ := net.Listen("tcp", "127.0.0.1:7777")
		for {
			listen.Accept()
		}
	}()
	main()
}
