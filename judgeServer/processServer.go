package main


type ProcessServer struct {
	Tasks map[int]SubmitTaskWrap
}

func RunProcessServer(processPort string, dispatchChannel chan SubmitTaskWrap) {
}
