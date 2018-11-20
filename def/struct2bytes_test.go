package def

import (
    "testing"
    "encoding/json"
)

func TestSubmit2Bytes(t *testing.T) {
    submit := Submit {
        SubmitID: 1,
        ProblemID: 2,
        CodeSource: []byte("Hello, world.c"),
        Language: 2,
    }
    data, _ := submit.StructToBytes()
    var result Submit;
    err := json.Unmarshal(data, &result)
    if err != nil || result.SubmitID != submit.SubmitID {
        t.Error(err)
    }
}

func TestResponse2Bytes(t *testing.T) {
    resp :=  Response {
        ErrCode: 1,
        Msg: []byte("Hello, world"),
    }
    data, _ := resp.StructToBytes()
    var result Response;
    err := json.Unmarshal(data, &result)
    if err != nil || result.ErrCode != resp.ErrCode{
        t.Error(err)
    }
}
