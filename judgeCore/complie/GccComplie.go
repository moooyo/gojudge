package complie

import (
	"../../def"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)
func GccComplie(submit *def.Submit)(err error){
	filename:="submit.c"
	err=ioutil.WriteFile(filename,submit.CodeSource,os.ModePerm)
	if err!=nil{
		return fmt.Errorf("Write SourceCode to File error")
	}
	cmd:=exec.Command("gcc","submit.c",`-osubmit`,"-O2","-std=c99")
	var out  bytes.Buffer
	cmd.Stderr=&out
	err=cmd.Run()
	if err!=nil {
		var status int
		fmt.Sscanf(err.Error(), "exit status %d", &status)
		return fmt.Errorf(out.String())
	}
	return
}