package main

import (
	"../def"
)
type JudgeTaskWrapStatus int


const (
	OK 	JudgeTaskWrapStatus = 0
	ERROR JudgeTaskWrapStatus = 1
)

type JudgeTaskWrap struct {
	Status    JudgeTaskWrapStatus
	Task *def.Submit
}
