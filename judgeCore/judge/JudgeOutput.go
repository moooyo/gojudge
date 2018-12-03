package judge

import (
	"bytes"
	"fmt"
	"io"
	"os"
)
import "../../def"

func judge(output, stdOutput string) (errno int, err error) {
	outputinfo, err := os.Stat(output)
	if err != nil {
		fmt.Printf("%v not found!\n", output)
	}
	stdOutputinfo, err := os.Stat(stdOutput)
	if err != nil {
		fmt.Printf("%v not found\n", stdOutput)
	}
	if outputinfo.Size() > stdOutputinfo.Size()+1 {
		return def.OuputLimitError, fmt.Errorf("outputSize larger than stdOutput")
	}
	out, _ := os.OpenFile(output, os.O_RDONLY, 0777)
	defer out.Close()
	out2, err := os.OpenFile(stdOutput, os.O_RDONLY, 0777)
	defer out2.Close()
	r1 := make([]byte, 1024)
	r2 := make([]byte, 1024)
	for {
		size1, err1 := out.Read(r1)
		size2, err2 := out2.Read(r2)
		if (err1 == io.EOF || err2 == io.EOF) && err1 != err2 {
			return def.WrongAnwser, fmt.Errorf("WrongAnwser")
		}
		if size1 != size2 {
			//	fmt.Printf("%d!=%d",size1,size2)
			return def.WrongAnwser, fmt.Errorf("WrongAnwser")
		}
		if bytes.Equal(r1, r2) != true {
			return def.WrongAnwser, fmt.Errorf("WrongAnwser")
		}
		if err1 == io.EOF || err2 == io.EOF {
			return def.AcceptCode, fmt.Errorf("")
		}
	}
	return def.AcceptCode, nil
}
