package main


type ProcessServer struct {
	Tasks map[int] JudgeTaskWrap
}

func RunProcessServer(processPort string, dispatchChannel chan JudgeTaskWrap) {
}
