package judge

import (
	"fmt"
	"log"
	"net"
	"testing"
)
import "../../def"
import "../../moudle"

var testElf=[]struct{
	filename string
	want int
}{
	{
		"./ac",
		def.AcceptCode,
	},
	{
		"./wa",
		def.WrongAnwser,
	},
	{
		"./tle",
		def.TimeLimitError,
	},
	{
		"./ole",
		def.OuputLimitError,
	},
}
const addr  string = "127.0.0.1:8888"
func TestElfJudge(t *testing.T) {
	var wantCode int
	var judgeItem int
	var problem def.Problem
	problem.TimeLimit=1000
	problem.MemoryLimit=256
	problem.JudgeList=[]def.JudgeNode{
		{
			"input.in",
			"output.out",
		},
	}
	syn := make(chan struct {},0)
	go func() {
		listen,err:=net.Listen("tcp",addr)
		syn <- struct{}{}
		if err!=nil{
			fmt.Printf("%v",err)
		}
		for {
			coon, _ := listen.Accept()
			for {
				var resp def.Response
				socket:=moudle.NewSocket(coon)
				socket.ReadStruct(&resp)
				if wantCode!=resp.ErrCode {
					log.Fatalf("test %d : want %d got %d\n",judgeItem,wantCode,resp.ErrCode)
				}
				if resp.ErrCode != def.AcceptCode {
					break
				}
				if resp.JudgeNode == resp.AllNode {
					break
				}
			}

		}
	}()
	<- syn
	for i,node:= range testElf {
		func() {
			fmt.Printf("i=%d\n",i)
			judgeItem=i
			wantCode=node.want
			coon, err := net.Dial("tcp", addr)
			if err!=nil{
				fmt.Printf("%v",err)
			}
			defer coon.Close()
			ElfJudge(node.filename, &problem, coon)
		}()
	}
}
