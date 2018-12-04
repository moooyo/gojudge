package submitwrap

import (
	"../../def"
	"fmt"
)

type SubmitTaskStatus int

const (
	OK      SubmitTaskStatus = iota
	ERROR
	WAITING
	JUDGING
)

type SubmitTaskWrap struct {
	Status SubmitTaskStatus
	Task   *def.Submit
}

func WrapSubmit(task *def.Submit) SubmitTaskWrap {
	return SubmitTaskWrap{
		Status: WAITING,
		Task:   task,
	}
}

func (submitTaskWrap *SubmitTaskWrap) String() string {
	var status string
	switch submitTaskWrap.Status {
	case OK:
		status = "OK"
	case ERROR:
		status = "ERROR"
	case WAITING:
		status = "WAITING"
	case JUDGING:
		status = "JUDGING"
	default:
		status = "unkonw"
	}
	return fmt.Sprintf("Status: %s Task: %s", status, submitTaskWrap.Task.String())
}
