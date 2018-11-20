package def

import (
	"encoding/json"
)

type Response struct {
	ErrCode int `json:"errCode"`
	Msg []byte `json:"msg"`
}

func (resp *Response)StructToBytes() (data []byte,err error){
	data,err =json.Marshal(resp)
	return
}