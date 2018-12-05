package main

import (
	"crypto/tls"
	"log"
	"net"
)

func main() {
	cer, err := tls.LoadX509KeyPair("./cert.pem", "./key.pem")
	if err != nil {
		log.Fatalln(err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	server, err := tls.Listen("tcp", "127.0.0.1:8080", config)
	if err != nil {
		log.Println(err)
		return
	}

	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	data := make([]byte, 1024)
	_, err := conn.Read(data)
	if err != nil {
		log.Print(err)
	}
	log.Print(string(data))
	conn.Close()
}
