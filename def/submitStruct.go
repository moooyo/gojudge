package def

import "encoding/json"

type Submit struct {
	SubmitID int `json:"submitID"`
	ProblemID int `json:"problemId"`
	CodeSource []byte `json:"codeSource"`
	Language int `json:"language"`

}
const (
	CLanguage=iota
	Cpp99Language
	Cpp11Language
	Cpp17Language
	JavaLanguage
)

func (resp *Submit)StructToBytes() (data []byte,err error){
	data,err =json.Marshal(resp)
	return
}