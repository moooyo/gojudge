package judge

import (
	"../../def"
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)
const(
	filename string="./submit"
)
func GccJudge(problem *def.Problem,conn net.Conn)(err error){
	list:=problem.JudgeList
	argv := make([]string,3)
	var resp def.Response
	for i,node:= range list{
		stdInputFileName:=node.Input
		stdOutputFileName:=strconv.Itoa(i)+".output"
		argv[0]="<"+stdInputFileName
		argv[1]=">"+stdOutputFileName
		ctx,cancel:=context.WithTimeout(context.Background(),time.Millisecond*time.Duration(500+problem.TimeLimit))
		defer cancel()
		cmd:=exec.CommandContext(ctx,"../submit")

		var out  bytes.Buffer

		cmd.Stderr=&out
		cmd.Stdin=os.OpenFile("../"+stdInputFileName)
		err=cmd.Run()
		fmt.Printf("%v",err)
		if err!=nil{
			resp.ErrCode=def.RunTimeError
			resp.Msg=out.Bytes()
			//moudle.StructWrite(conn,&resp)
			fmt.Printf("%v",resp.Msg)
			return err
		}
		err=judge(i)
	}
	return
}
func judge(i int)(err error){
	//todo
	return nil
}