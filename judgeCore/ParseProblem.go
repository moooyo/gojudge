package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

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

func ParseProblemFile(filename string,problem *Problem)(err error){
	file,err:=ioutil.ReadFile(filename)
	if err!=nil{
		return
	}
	err=json.Unmarshal(file,problem)
	if err!=nil{
		return
	}
	if len(problem.JudgeList)==0{
		return fmt.Errorf("empty problem file")
	}
	return nil
}
