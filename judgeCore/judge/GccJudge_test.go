package judge

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)
import "../../def"
func TestGccJudge(t *testing.T) {
	var problem def.Problem
	data,_:=ioutil.ReadFile("../test.json")
	json.Unmarshal(data,&problem)
	GccJudge(&problem,nil)

}
