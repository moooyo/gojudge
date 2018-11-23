package complie

import (
	"os"
	"testing"
)
import "../../def"

var submit def.Submit
var test= []struct {
	testFile []byte
	want	bool
}{
	{[]byte(`#include<stdio.h>
			int main()
			{
				printf("hello gojudge!\n");
				return 0;
			}
`), true},
{
	[]byte(`#include<iostream>
					usingnamespace std;
					int main()
					{
						cout<<"hello gojudge\n"<<endl;
						return 0;
					}
`),false}}
func TestGccComplie(t *testing.T) {
	var submit=def.Submit{
		1000,
		1000,
		nil,
		0,
	}
	for _,ts := range test{
		submit.CodeSource=ts.testFile
		err:=GccComplie(&submit)
		os.Remove("submit")
		var got = err==nil
		if got!=ts.want{
			t.Errorf("test FAILD %s ",ts.testFile)
		}
	}
}
