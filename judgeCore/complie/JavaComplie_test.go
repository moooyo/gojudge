package complie

import (
	"../../def"
	"os"
	"testing"
)

var testJava = []struct {
	testFile []byte
	want     bool
}{
	{[]byte(`
public class Main{
	public static void main(String []argv)
	{
		
	}
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
`), false},
	{[]byte(`
	public class test{
	public static void main(String []argv)
	{
		
	}
}
`), false}}

func TestJavaComplie(t *testing.T) {
	var submit = def.Submit{
		1000,
		1000,
		nil,
		0,
	}
	for _, ts := range testJava {
		defer os.Remove("Main.class")
		submit.CodeSource = ts.testFile
		err := JavaComplie(&submit)
		var got = err == nil
		if got != ts.want {
			t.Errorf("test FAILD %s \nerr:%v\n", ts.testFile, err.Error())
		}
	}
}
