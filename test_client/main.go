package main

import (
	"../def"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"unsafe"
)

func main() {
	client := http.Client{}
	for i := 0; i < 1000; i++ {
		submit := def.Submit{
			SubmitID:   i,
			ProblemID:  1000,
			CodeSource: []byte("#include <stdio.h> \n int main(void) { printf(\"1 2 3 4 5\"); return 0;}"),
			Language:   def.CLanguage,
		}
		data, _ := submit.StructToBytes()
		reader := bytes.NewReader(data)
		url := "http://127.0.0.1:8080/submit_task"
		request, err := http.NewRequest("POST", url, reader)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		request.Header.Set("Content-Type", "application/json;charset=UTF-8")
		resp, err := client.Do(request)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		str := (*string)(unsafe.Pointer(&respBytes))
		fmt.Println(*str)
		resp.Body.Close()
	}
}
