package complie

import (
	"../../def"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

func GccComplie(submit *def.Submit) (err error) {
	err = ParseConfig()
	if err != nil {
		fmt.Print(err)
		panic("parse config error "+err.Error())
	}
	filename := "submit.c"
	err = ioutil.WriteFile(filename, submit.CodeSource, os.ModePerm)
	defer os.Remove(filename)

	if err != nil {
		return fmt.Errorf("Write SourceCode to File error")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(Config["gcc"].TimeLimit))
	defer cancel()
	cmd := exec.CommandContext(ctx, "gcc")
	cmd.Args = Config["gcc"].Argv
	var out bytes.Buffer
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		switch err.Error() {
		default:
			return fmt.Errorf(out.String())
		case "signal: killed":
			return fmt.Errorf("complie timeout")
		}
	}
	return
}
