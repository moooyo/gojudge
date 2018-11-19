package def

type Submit struct {
	SubmitID int `json:"submitID"`
	ProblemID int `json:"problemId"`
	CodeSource []byte `json:"codeSOurce"`
	Language int `json:"language"`
}

