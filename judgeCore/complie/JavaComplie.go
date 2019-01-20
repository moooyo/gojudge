package complie

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ferriciron/gojudge/def"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

func JavaComplie(submit *def.Submit) (err error) {
	err = ParseConfig()
	if err != nil {
		panic(err)
	}
	filename := "Main.java"
	err = ioutil.WriteFile(filename, submit.CodeSource, os.ModePerm)
	defer os.Remove(filename)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(Config["java"].TimeLimit))
	defer cancel()
	cmd := exec.CommandContext(ctx, "javac")
	cmd.Args = Config["java"].Argv
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
