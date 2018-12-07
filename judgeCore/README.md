# JudgeCore

此部分为评测核心部分，承担判评题目的功能。相关的评测选项，评测设置需要在此部分中进行配置。

## 配置

### docker 

该评测机推荐(在开发过程中依赖于)docker环境运行，以确保进行资源隔离。  

ComplieEnv中构建了docker的编译环境，对于该评测机的运行环境来说，可以通过该Dockerfile进行修改。  

JudgeCore的运行环境从gojudgecomplieenv(ComplieEnv)开始构建，为了进行构建需要对judgeCore部分进行编译并生成二进制文件judgeCore, 并且需要依赖于complie中的config.json文件进行相关参数设置，你可以适当的修改对应的Dockerfile文件或者config.json文件，但除非你明白你正在做什么，你不应该破坏JudgeCore在docker中运行的相关路径信息。  

JudgeCore在docker中的运行环境从/home/gojudge/judgeCore开始，其余部分并不关心。为了获取评测信息，你需要把题目挂载到/home/gojudge/judgeCore/problem下。例如，对于题目编号为1000和1001的两组题目，挂载路径应该为/home/gojudge/judgeCore/problem/1000和/home/gojudge/judgeCore/problem/1001

### 题目定义

本评测机的题目由一个json文件problem定义。

```json
// problem.json
{
  "timeLimit":1000,
  "memorylimit":256,
  "judgelist":[
    {
      "input":"input.in",
      "output":"output.out"
    },
    {
      "input":"input2.in",
      "output":"output2.out"   
    }
  ],
  "property":0
}
```

该文件定义了一个时间限制为1000ms，内存限制为256m,并且含有两组数据。一组为input.in和output.out。另一组为input2.in, output2.out。property字段确定了该题目的评测类型。

| 字段        | 含义         | 类型   |
| ----------- | ------------ | ------ |
| timeLimit   | 时间限制(ms) | number |
| memorylimit | 空间限制(M)  | number |
| judgelist   | 评测数据     | array  |
| property    | 评测类型     | number |

对property的定义如下

| property | 含义     |
| -------- | -------- |
| 0        | 普通题目 |
| 1        | spj      |
| 2        | 交互题   |

题目应该被组织成一个文件夹，该文件夹的名称为题目编号，文件夹内应该包含对应的problem.json和对应的评测数据。此文件夹应该在评测时被挂载到对应的problem路径下。

### 编译环境

注意，此评测环境中Java代码主类名应该为Main

对编译环境的配置主要依靠complieEnv下的Dockerfile和complie下的config文件。除非你知道你在做什么否则不要更改相关文件。  

如果需要对相关配置进行更改，可以在config文件中更改对应的编译命令以开启某些编译选项(如O2优化),但这可能需要你更改complieEnv来修改对应的环境实现。默认提供的编译命令如下

#### GCC

```com
gcc -osubmit submit.c -O2 -std=c99
```

### JAVA

```com
javac Main.java
```

## 使用

为了使用JudgeCore，您需要在配置好相关环境之后以合适的参数拉起Docker实例(主要为对资源的限制和相关文件的挂载),并且需要传入adress、port、submitID三个参数启动JudgeCore程序。这一部分可以参考JudgeServer中的相关配置。

拉起实例后，JudgeCore会首先通过adress和port与JudgeServer建立连接，您应该监听对应的地址和端口，以建立连接。建立连接后，JudgeCore会首先将submitID发送给JudgeServer以确认身份。接着JudgeCore需要接受一个评测请求submit,类型为JSON对象,定义如下

```go
type Submit struct {
	SubmitID   int    `json:"submitID"`
	ProblemID  int    `json:"problemId"`
	CodeSource []byte `json:"codeSource"`
	Language   int    `json:"language"`
}
const (
	CLanguage = iota
	Cpp99Language
	Cpp11Language
	Cpp17Language
	JavaLanguage
)
```

接收到对应的submit后，JudgeCore会在相应的路径上利用Language和CodeSource编译程序，并利用ProblemID获取对应题目信息进行评测。

评测过程中JudgeCore会不断地向JudgeServer发送评测分点信息response,类型为JSON对象。

```go
type Response struct {
	ErrCode   	int    	`json:"errCode"`
	JudgeNode 	int    	`json:"judgeNode"`
	AllNode   	int    	`json:"allNode"`
	TimeCost	int 	`json:"timecost"`
	Msg       	[]byte 	`json:"msg"`
}
const (
	JudgeFinished = iota
	AcceptCode
	WrongAnwser
	ComplierError
	TimeLimitError
	ComlierTimeLimitError
	MemoryLimitError
	OuputLimitError
	RunTimeError
	OtherError = -1
)
```

对应字段和含义为

| 字段      | 含义           |
| --------- | -------------- |
| ErrCode   | 错误码(见定义) |
| JudgeNode | 评测点序号     |
| AllNode   | 总计评测点数量 |
| TimeCost  | 运行时间       |
| Msg       | 错误信息       |

对于某些错误，可以从Msg中拿到对应的错误信息。

评测结束于某一评测点评测结果不为AcceptCode或评测了所有评测点(JudgeNode=AllNode),您应该通过对应的信息判断总评测结果,需要注意的是，有可能因为意外错误导致网络中断或套接字被关闭。JudgeServer中需要进行相应的处理。

分点传输的特性有助于进行详细评测信息的显示。







