package judge

import (
	"fmt"
	"gopl.io/ch13/equal"
	"testing"
)
import "../../def"
var testJudge =[]struct{
	output string
	stdOutput string
	wantCode int
	wantErr error
}{
	{
		"./test/ac.out",
		"./test/output.out",
		def.AcceptCode,
		fmt.Errorf(""),
	},
	{
		"./test/wa.out",
		"./test/output.out",
		def.WrongAnwser,
		fmt.Errorf("WrongAnwser"),
	},
	{
		"./test/ole.out",
		"./test/output.out",
		def.OuputLimitError,
		fmt.Errorf("outputSize larger than stdOutput"),
	},
}
func TestJudge(t *testing.T) {
	for _,node:=range testJudge{
		code,err:=judge(node.output,node.stdOutput)
		if code!=node.wantCode||equal.Equal(err,node.wantErr)!=true{
			t.Errorf("got %d want %d",code,node.wantCode)
		}
	}
}
