package main

import (
	"../def"
	"../moudle"
	"log"
	"net"
)

func main() {

	for i := 0; i < 10; i++ {
		conn, err := net.Dial("tcp", "127.0.0.1:8080")
		if err != nil {
			log.Fatal(err)
		}

		socket := moudle.NewSocket(conn)

		submit := def.Submit{
			SubmitID:   i,
			ProblemID:  1000,
			CodeSource: []byte("#include <stdio.h> \n int main(void) { printf(\"1 2 3 4 5\"); return 0;}"),
			Language:   def.CLanguage,
		}
		err = socket.WriteStruct(&submit)

		if err != nil {
			log.Println(err)
		}

		socket.Close()
	}
}
