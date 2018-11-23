package def

import (
	"encoding/json"
)

type Response struct {
	ErrCode   	int    	`json:"errCode"`
	JudgeNode 	int    	`json:"judgeNode"`
	AllNode   	int    	`json:"allNode"`
	TimeCost	int 	`json:"timecost"`
	Msg       	[]byte 	`json:"msg"`
}

const(
	JudgeFinished=iota
	AcceptCode
	WrongAnwser
	ComplierError
	TimeLimitError
	ComlierTimeLimitError
	MemoryLimitError
	OuputLimitError
	RunTimeError
	OtherError=-1
)

func (resp *Response)StructToBytes() (data []byte,err error){
	data,err =json.Marshal(resp)
	return
}