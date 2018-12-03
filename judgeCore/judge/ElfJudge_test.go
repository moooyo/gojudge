package judge

import (
	"fmt"
	"log"
	"net"
	"testing"
)
import "../../def"
import "../../moudle"

var testElf = []struct {
	filename string
	want     int
}{
	{
		"./test/ac",
		def.AcceptCode,
	},
	{
		"./test/wa",
		def.WrongAnwser,
	},
	{
		"./test/tle",
		def.TimeLimitError,
	},
	{
		"./test/ole",
		def.OuputLimitError,
	},
	{
		"./test/re",
		def.RunTimeError,
	},
	{
		"./test/mle",
		def.MemoryLimitError,
	},
}

const addr string = "127.0.0.1:8888"

func TestElfJudge(t *testing.T) {
	var problem def.Problem
	problem.TimeLimit = 1000
	problem.MemoryLimit = 256
	problem.JudgeList = []def.JudgeNode{
		{
			"./test/input.in",
			"./test/output.out",
		},
	}
	syn := make(chan struct{}, 0)
	testData := make(chan int, 0)
	go func() {
		listen, err := net.Listen("tcp", addr)
		syn <- struct{}{}
		if err != nil {
			fmt.Printf("%v", err)
		}
		for {
			rs := <-testData
			coon, _ := listen.Accept()
			func() {
				defer coon.Close()
				for {
					var resp def.Response
					socket := moudle.NewSocket(coon)
					socket.ReadStruct(&resp)
					fmt.Printf("%v\n", resp)
					wantCode := testElf[rs].want
					judgeItem := rs + 1
					if wantCode != resp.ErrCode {
						log.Fatalf("test %d : want %d got %d\n", judgeItem, wantCode, resp.ErrCode)
					}
					if resp.ErrCode != def.AcceptCode {
						break
					}
					if resp.JudgeNode == resp.AllNode {
						break
					}
				}
			}()

		}
	}()
	<-syn
	for i, node := range testElf {
		func() {
			//	fmt.Printf("i=%d\n",i+1)
			testData <- i
			coon, err := net.Dial("tcp", addr)
			if err != nil {
				fmt.Printf("%v", err)
			}
			defer coon.Close()
			ElfJudge(node.filename, &problem, coon)
		}()
	}
}
