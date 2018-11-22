package main

import "testing"

func TestParseProblemFile(t *testing.T) {
	filename:="test.json"
	var problem Problem
	err:=ParseProblemFile(filename,&problem)
	if err!=nil{
		t.Errorf("%v",err)
	}
	if len(problem.JudgeList)!=2{
		t.Errorf("Parse JudgeList Error len %d != 2", len(problem.JudgeList))
	}

}