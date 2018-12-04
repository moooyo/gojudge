package judge

import (
	"fmt"
	"log"
	"net"
	"testing"
)
import "../../def"
import "../../moudle"



var testJava = []struct {
	filename string
	want     int
}{
	{
		"Ac",
		def.AcceptCode,
	},
	{
		"Wa",
		def.WrongAnwser,
	},
	{
		"Tle",
		def.TimeLimitError,
	},
}

func TestJavaJudge(t *testing.T) {
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
		go func() {
			<-syn
			listen.Close()
		}()
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
					//fmt.Printf("%v\n", resp)
					wantCode := testElf[rs].want
					judgeItem := rs + 1
					if wantCode != resp.ErrCode {
						fmt.Printf("%s\n",string(resp.Msg))
						log.Fatalf("Java test %d : want %d got %d\n", judgeItem, wantCode, resp.ErrCode)
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
	defer func() {
		syn <- struct{}{}
	}()
	for i, node := range testJava {
		func() {
			//	fmt.Printf("i=%d\n",i+1)
			testData <- i
			coon, err := net.Dial("tcp", addr)
			if err != nil {
				fmt.Printf("%v", err)
			}
			defer coon.Close()
			JavaJudge(node.filename,&problem, coon)
		}()
	}
}
