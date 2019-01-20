package judge

import (
	"context"
	"fmt"
	"github.com/ferriciron/gojudge/def"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func ElfJudge(basePath string, judgeFile string, problem *def.Problem, conn net.Conn) (err error) {
	filename := judgeFile
	//	fmt.Println(filename)
	list := problem.JudgeList
	for i, node := range list {
		inputFileName := basePath + "/" + node.Input
		outputFileName := "judge_" + strconv.Itoa(i) + ".output"
		ct := func() bool {
			//clean output
			defer func() {
				os.Remove(outputFileName)
			}()
			var cost int64
			var resp def.Response
			resp.JudgeNode = i + 1
			resp.AllNode = len(problem.JudgeList)
			defer sendResponse(conn, &resp)
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(problem.TimeLimit)*time.Millisecond)
			defer cancel()
			cmd := exec.CommandContext(ctx, filename)
			cmd.Stdin, err = os.OpenFile(inputFileName, os.O_RDONLY, 0777)
			if err != nil {
				return buildResponse(&resp, def.OtherError,
					fmt.Sprintf("open file error:%v", err))
			}

			//ensure outputFile not exit
			os.Remove(outputFileName)
			cmd.Stdout, err = os.OpenFile(outputFileName, os.O_CREATE|os.O_WRONLY, 0777)
			if err != nil {
				return buildResponse(&resp, def.OtherError,
					fmt.Sprintf("write file error:%v", err))
			}
			start := time.Now()
			err = cmd.Run()
			end := time.Now()
			cost = (end.UnixNano() - start.UnixNano()) / (1000 * 1000)
			resp.TimeCost = int(cost)
			if err != nil {
				switch err.Error() {
				default:
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
			code, err := judge(outputFileName, basePath+"/"+node.Output)
			return buildResponse(&resp, code, err.Error())
		}()
		if !ct {
			break
		}
	}
	return
}
