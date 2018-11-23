package def

import "fmt"

type Problem struct {
	TimeLimit   int         `timelimit`
	MemoryLimit int         `memorylimit`
	JudgeList   []JudgeNode `judgelist`
	property    int         `property`
}
type JudgeNode struct {
	Input  string `input`
	Output string `output`
}

func (problem *Problem) String() string {
	return fmt.Sprintf("TimeLimit: %dms MemoryLimit: %dm property: %d", problem.TimeLimit, problem.MemoryLimit, problem.property)
}
