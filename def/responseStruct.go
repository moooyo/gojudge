package def

import (
	"encoding/json"
)

type Response struct {
	ErrCode int `json:"errCode"`
	Msg []byte `json:"msg"`
}
const(
	JudgingResponseCode=iota
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