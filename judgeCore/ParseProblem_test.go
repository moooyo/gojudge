package main

import "testing"
import "../def"
func TestParseProblemFile(t *testing.T) {
	filename:="example.json"
	var problem def.Problem
	err:=ParseProblemFile(filename,&problem)
	if err!=nil{
		t.Errorf("%v",err)
	}
	if len(problem.JudgeList)!=2{
		t.Errorf("Parse JudgeList Error len %d != 2", len(problem.JudgeList))
	}

}