package main



type ProcessServerConfig struct {
    Port    int
}

type ProcessServer struct {
	Tasks map[int]SubmitTaskWrap
}

func RunProcessServer(processPort string, dispatchChannel chan SubmitTaskWrap) {
}
