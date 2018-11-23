package main

import (
	"../def"
)
type SubmitTaskStatus int


const (
	OK    SubmitTaskStatus = 0
	ERROR SubmitTaskStatus = 1
)

type SubmitTaskWrap struct {
	Status SubmitTaskStatus
	Task   *def.Submit
}
