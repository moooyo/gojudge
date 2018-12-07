package moudle

import (
	"../def"
	"os"
	"reflect"
	"testing"
)

func TestCode(t *testing.T) {
	file, err := os.Create("test_code")
	if err != nil {
		t.Fatal(err)
	}
	coder := NewDECoder(1024 * 6)

	for i := 0; i < 10; i++ {
		submit := def.Submit{
			SubmitID:   i,
			ProblemID:  i,
			CodeSource: make([]byte, 4096),
			Language:   i,
		}
		err = coder.AppendStruct(&submit)
		coder.AppendInt(i)
		if err != nil {
			t.Fatal(err)
			file.Close()
			os.Remove(file.Name())
		}
	}
	err = coder.Send(file)
	if err != nil {
		file.Close()
		os.Remove(file.Name())
		t.Fatal(err)
	}

	file.Close()

	file, err = os.Open("test_code")

	defer file.Close()
	defer os.Remove(file.Name())
	for i := 0; i < 10; i++ {
		submit := def.Submit{
			SubmitID:   i,
			ProblemID:  i,
			CodeSource: make([]byte, 4096),
			Language:   i,
		}
		var result def.Submit

		err = coder.ReadStruct(file, &result)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(result, submit) {
			t.Fail()
		}
		value, err := coder.ReadInt(file)
		if err != nil {
			t.Fatal(err)
		}
		if value != i {
			t.Fail()
		}
	}
}
