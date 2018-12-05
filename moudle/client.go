package main

import (
	"crypto/tls"
	"log"
)

func main() {

	log.SetFlags(log.Lshortfile)

	cer, err := tls.LoadX509KeyPair("./cert.pem", "./key.pem")
	if err != nil {
		log.Fatal(err)
	}

	conf := &tls.Config{
		Certificates: []tls.Certificate{cer},
	}

	conn, err := tls.Dial("tcp", "127.0.0.1:8080", conf)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	n, err := conn.Write([]byte("hello\n"))
	if err != nil {
		log.Println(n, err)
		return
	}

	conn.Close()

}
