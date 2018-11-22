package judge

import (
	"../../def"
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
	signlParse string="signal: %s"
)
const (
	TimeLimitError="killed"
)

func ElfJudge(problem *def.Problem,conn net.Conn)(err error){

	list:=problem.JudgeList
	for i,node:= range list{
		var resp def.Response
		inputFileName:=node.Input
		outputFileName:="judge_"+strconv.Itoa(i)+".output"
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(problem.TimeLimit)*time.Millisecond)
			defer cancel()
			cmd := exec.CommandContext(ctx, filename)
			cmd.Stdin, err = os.OpenFile(inputFileName, os.O_RDONLY, 0777)
			if err != nil {
				fmt.Printf("open file error:%v", err)
			}
			cmd.Stdout, err = os.OpenFile(outputFileName, os.O_CREATE|os.O_WRONLY, 0777)
			if err != nil {
				fmt.Printf("write file error:%v", err)
			}
			err = cmd.Run()
			var signalErr string
			if err != nil {
				fmt.Sscanf(err.Error(), signlParse, &signalErr)
				fmt.Printf("%s\n%v\n", signalErr, err)
			}
			code,err := judge(outputFileName,node.Output)
			fmt.Printf("%v",err)
			resp.ErrCode=code
			if err!=nil {
				resp.Msg = []byte(err.Error())
			}
			//moudle.StructWrite(conn,&resp)
		}()
	}
	return
}

