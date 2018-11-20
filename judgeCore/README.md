# JudgeCore

此部分为评测核心部分，承担判评题目的功能

## 题目定义

题目应该将输入和输出文件以及spj(可选)存放到同一目录下，并在同级目录定义problem.json，评测核心依靠problem.jsone的信息进行评测。

```json
#simple problem
{
 	"timeLimit":1000,
    "memorylimit":256,
    "judgelist":[
        {
            "input":"test1.in",
            "output":"output1.out"
        },
        {
            "input":"input2.in",
            "output":"output2.out"
        }
    ],
    "property":0
}
```

上述problem.json定义了一个时限为1s,内存限制为256m的题目，该题目为普通题目，有两组测试数据。