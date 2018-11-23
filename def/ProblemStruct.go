package def

type Problem struct{
	TimeLimit 		int 			`timelimit`
	MemoryLimit 	int 			`memorylimit`
	JudgeList		[]JudgeNode		`judgelist`
	property		int 			`property`
}
type JudgeNode struct{
	Input		string		`input`
	Output		string 		`output`
}
