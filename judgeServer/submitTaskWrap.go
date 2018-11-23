package main

import (
	"../def"
)

type SubmitTaskStatus int

const (
	OK      SubmitTaskStatus = 0
	ERROR   SubmitTaskStatus = 1
	WAITING SubmitTaskStatus = 2
	JUDGING SubmitTaskStatus = 3
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
