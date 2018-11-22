package judge

import (
	"../../def"
	"../../moudle"
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
		inputFileName:=node.Input
		outputFileName:="judge_"+strconv.Itoa(i)+".output"
		ct:=func() bool{
			var resp def.Response
			defer sendResponse(conn,&resp)
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(problem.TimeLimit)*time.Millisecond)
			defer cancel()
			cmd := exec.CommandContext(ctx, filename)
			cmd.Stdin, err = os.OpenFile(inputFileName, os.O_RDONLY, 0777)
			if err != nil {
				return buildResponse(&resp,def.OtherError,
					fmt.Sprintf("open file error:%v", err))
			}
			cmd.Stdout, err = os.OpenFile(outputFileName, os.O_CREATE|os.O_WRONLY, 0777)
			if err != nil {
				return buildResponse(&resp,def.OtherError,
					fmt.Sprintf("write file error:%v", err))
			}
			err = cmd.Run()
			var signalErr string
			if err != nil {
				fmt.Sscanf(err.Error(), signlParse, &signalErr)
				fmt.Printf("%s\n%v\n", signalErr, err)
				switch signalErr {
				default:
					return buildResponse(&resp,def.OtherError,fmt.Sprintf("%v",err))
				case TimeLimitError:
					return buildResponse(&resp,def.TimeLimitError,"")
				}
			}
			code,err := judge(outputFileName,node.Output)
			return buildResponse(&resp,code,err.Error())
		}()
		if !ct{
			break
		}
	}
	return
}
func buildResponse(resp *def.Response,code int ,msg string) bool{
	resp.ErrCode=code
	resp.Msg=[]byte(msg)
	return resp.ErrCode==def.AcceptCode
}
func sendResponse(conn net.Conn,resp *def.Response){
	moudle.StructWrite(conn,resp)
}
