package test

import (
	"../def"
	"fmt"
	"net"
	"../moudle"
)
func main(){
	var submit def.Submit
	submit.Language=def.CLanguage
	submit.CodeSource=[]byte(`
		#include<stdio.h>
		int main()
		{
			printf("1 2 3 4 5");
			return 0;
		}
`)
	submit.ProblemID=1000
	submit.SubmitID=1000
	conn, err :=net.Dial("tcp","127.0.0.1:8080")
	if err!=nil{
		panic("open net error")
	}
	socket:=moudle.NewSocket(conn)
	socket.WriteStruct(&submit)
	var resp def.Response
	socket.ReadStruct(&resp)
	fmt.Print(resp)
	fmt.Print(string(resp.Msg))
	return
}
