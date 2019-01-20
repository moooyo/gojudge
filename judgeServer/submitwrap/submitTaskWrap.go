package submitwrap

import (
	"fmt"
	"github.com/ferriciron/gojudge/def"
	"time"
)

type SubmitTaskStatus int

const (
	OK SubmitTaskStatus = iota
	TIMEOUTERROR
	WAITING
	JUDGING
	EXECUTING
	EXECUTERROR
	UNKONW
)

type SubmitTaskWrap struct {
	Status      SubmitTaskStatus
	ExecuteTime time.Duration
	ProcessTime time.Duration
	TimeOut     time.Duration
	ExecCount   int
	Task        *def.Submit
}

func WrapSubmit(task *def.Submit) SubmitTaskWrap {
	return SubmitTaskWrap{
		Status:      WAITING,
		Task:        task,
		ExecCount:   0,
		ExecuteTime: 0,
		ProcessTime: 0,
		TimeOut:     0,
	}
}

func (submitTaskWrap *SubmitTaskWrap) String() string {
	var status string
	switch submitTaskWrap.Status {
	case OK:
		status = "OK"
	case TIMEOUTERROR:
		status = "TIMEOUTERROR"
	case WAITING:
		status = "WAITING"
	case JUDGING:
		status = "JUDGING"
	default:
		status = "unkonw"
	}
	return fmt.Sprintf("Status: %s Task: %s", status, submitTaskWrap.Task.String())
}
