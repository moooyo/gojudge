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
const JAVABASE = 2 //java timelimit must mul JAVABASE
func JavaJudge(JudgeFile string,problem *def.Problem,conn net.Conn)(err error){
	filename:=JudgeFile
	judgelist:=problem.JudgeList
	for i,node:=range judgelist{
		inputFilename:=node.Input
		outputFileName := "judge_" + strconv.Itoa(i) + ".output"
		ct := func() bool {
			defer os.Remove(outputFileName)
			var resp def.Response
			resp.JudgeNode = i + 1
			resp.AllNode = len(problem.JudgeList)
			defer sendResponse(conn, &resp)
			ctx,cancel:=context.WithTimeout(context.Background(),time.Duration(problem.TimeLimit*JAVABASE)*time.Millisecond)
			defer cancel()

			// Command: java Main
			// It means that the Judged Java SourceCode
			// Must have public class Main. You can change it in there
			// and if you changed you must change JavaComplie.go
			cmd:=exec.CommandContext(ctx,"java",filename)
			//fmt.Print(cmd.Args)
			cmd.Stdin,err=os.OpenFile(inputFilename,os.O_RDONLY,0777)
			if err != nil {
				return buildResponse(&resp, def.OtherError,
					fmt.Sprintf("open file error:%v", err))
			}

			//make output clean
			os.Remove(outputFileName)
			cmd.Stdout,err=os.OpenFile(outputFileName,os.O_CREATE|os.O_WRONLY,0777)
			if err != nil {
				return buildResponse(&resp, def.OtherError,
					fmt.Sprintf("write file error:%v", err))
			}
			var errStream bytes.Buffer
			cmd.Stderr=&errStream
			startTime:=time.Now()
			err=cmd.Run()
			endTime:=time.Now()
			cost:=startTime.UnixNano()-endTime.UnixNano()/(1000 * 1000)
			resp.TimeCost=int(cost)
			if err != nil {
				switch err.Error() {
				default:
					fmt.Print(errStream.String())
					return buildResponse(&resp, def.OtherError, fmt.Sprintf("%v", err))
				case SIGKILL:
					if (float64(problem.TimeLimit-int(cost)) / float64(problem.TimeLimit)) < eps {
						return buildResponse(&resp, def.TimeLimitError, "")
					} else {
						return buildResponse(&resp, def.MemoryLimitError, "")
					}
				case RuntimeError:
					return buildResponse(&resp, def.RunTimeError, RuntimeError)
				}
			}
			code,err:=judge(outputFileName,node.Output)
			return buildResponse(&resp,code,err.Error())
		}()
		if !ct{
			break
		}
	}
	return 
}
