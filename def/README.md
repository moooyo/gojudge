# 类型定义

本项目中socket之间交互采用小端字节序，交互格式参照

```json
{
    length,
    data
}
```

即先往socket里传递数据体字节数n，接着传递大小为n的数据体

## SubmitStruct

此结构体用于JudgeServer与JudgeCore之间的交互。JudgeServer应该在收到处于docker中的JudgeCore发送的submitID之后回应一个SubmitStruct表示此次评测的信息。

```go

type Submit struct {
	SubmitID int `json:"submitID"`
	ProblemID int `json:"problemId"`
	CodeSource []byte `json:"codeSource"`
	Language int `json:"language"`
}
```

| 字段       | 含义           | 类型   |
| ---------- | -------------- | ------ |
| SubmitID   | 提交号         | int    |
| ProblemID  | 待评测题目编号 | int    |
| CodeSource | 待评测代码     | []byte |
| Language   | 待评测语言     | int    |

## ResponseStruct

此结构体用于JudgeCore返回到JudgeServer的对应的评测结果。利用ErrCode区分评测结果类型，从Msg获取对应信息。

```go
type Response struct {
	ErrCode int `json:"errCode"`
	Msg []byte `json:"msg"`
}
```

| 字段    | 含义   | 类型   |
| ------- | ------ | ------ |
| ErrCode | 错误码 | int    |
| Msg     | 信息   | []byte |

如果ErrCode表示为Judging，则Msg应该含有评测点信息,即Msg应该遵循以下约定

```go
type JudgingMsg struct{
    JudgePoint 		int 	`json:"judgePoint"`
    AllPoint		int		`json:"allPoint"`
    TimeLimit		int		`json:"timelimit"`
    MemoryLimit		int		`json:"memoryLimit"`
}
```

| 字段        | 含义                    | 类型 |
| ----------- | ----------------------- | ---- |
| JudgePoint  | 当前评测点位置          | int  |
| AllPoint    | 总共评测点数量          | int  |
| TimeLimit   | 当前评测点时间限制(ms)  | int  |
| MemoryLimit | 当前评测点内存限制(Mib) | int  |



## 枚举类型

| Language | 评测语言 |
| -------- | -------- |
| 0        | c        |
| 1        | cpp 99   |
| 2        | cpp 11   |
| 3        | cpp 17   |
| 4        | Java     |

| ErrCode | 错误码                 |
| ------- | ---------------------- |
| 0       | JudgingResponseCode    |
| 1       | AcceptResponseCode     |
| 2       | WrongAnwser            |
| 3       | ComplierError          |
| 4       | TimeLimitError         |
| 5       | ComplierTimeLimitError |
| 6       | MemoryLimitError       |
| 7       | OutputLimitError       |
| -1      | OtherError             |

