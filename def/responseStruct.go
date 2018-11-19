package def

type response struct {
	ErrCode int `json:"errCode"`
	Msg []byte `json:"msg"`
}
