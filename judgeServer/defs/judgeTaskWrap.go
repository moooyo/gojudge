package defs


type JudgeTaskWrapStatus int

const (
	OK 	JudgeTaskWrapStatus = 0
	ERROR JudgeTaskWrapStatus = 1
)

type JudgeTask int
type JudgeTaskWrap struct {
	Status    JudgeTaskWrapStatus
	Task *JudgeTask
}
