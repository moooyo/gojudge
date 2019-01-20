package def

import (
	"encoding/json"
	"fmt"
)

type Submit struct {
	SubmitID   int    `json:"submitID"`
	ProblemID  int    `json:"problemID"`
	CodeSource []byte `json:"codeSource"`
	Language   int    `json:"language"`
}

const (
	CLanguage     = 1
	Cpp99Language = 2
	Cpp11Language = 3
	Cpp17Language = 4
	JavaLanguage  = 5
)

func (submit *Submit) StructToBytes() (data []byte, err error) {
	data, err = json.Marshal(submit)
	return
}

func (submit *Submit) String() string {
	var language string
	switch submit.Language {
	case CLanguage:
		language = "c"
	case Cpp99Language:
		language = "c++99"
	case Cpp11Language:
		language = "c++11"
	case Cpp17Language:
		language = "c++17"
	case JavaLanguage:
		language = "java"
	default:
		language = "unkonw"
	}
	return fmt.Sprintf("submitID: %d problemId: %d language: %s", submit.SubmitID, submit.ProblemID, language)
}
