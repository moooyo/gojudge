package main

import (
	"../def"
	"../moudle"
	"log"
)

func main() {

	for i := 0; i < 100; i++ {
		conn, err := moudle.Dial("127.0.0.1:8080")
		if err != nil {
			log.Fatal(err)
		}
		submit := def.Submit{
			SubmitID:   i,
			ProblemID:  1000,
			CodeSource: []byte("#include <stdio.h> \n int main(void) { printf(\"1 2 3 4 5\"); return 0;}"),
			Language:   def.CLanguage,
		}
		encoder := moudle.NewEnCoder()
		encoder.AppendStruct(&submit)
		encoder.Send(conn)

		if err != nil {
			log.Println(err)
		}
		conn.Close()
	}
}
