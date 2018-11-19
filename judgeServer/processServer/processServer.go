package processServer

import "../defs"

type ProcessServer struct {
	Tasks map[int] defs.JudgeTaskWrap
}

func RunProcessServer(processPort string, dispatchChannel chan defs.JudgeTaskWrap) {
}
