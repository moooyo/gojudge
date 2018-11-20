package def

import "encoding/json"

type Submit struct {
	SubmitID int `json:"submitID"`
	ProblemID int `json:"problemId"`
	CodeSource []byte `json:"codeSOurce"`
	Language int `json:"language"`

}


func (resp *Submit)StructToBytes() (data []byte,err error){
	data,err =json.Marshal(resp)
	return
}