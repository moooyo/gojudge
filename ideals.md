

```json
-----> judeServer
struct JudgeTask {
 	"submitID":
    "problemID":
    "codeSource":
    "language":
} ------> judgeCore
```

```json




Interactive process

judeServer recv judgeTask from web server

judeServer dispatch a task and run a judgeCore in docker with args (submitID, ip, port)

then judgeCore connect to judeServer with submitID

judgeServer send judgeTask info to judgeCore

then judgeServer recv message from judgeCore

    judgeCore                                       judgeServer                   
                    send submitID
                ------------------------>
                      judgeTask
                <------------------------
                     send status
                ------------------------>

```

```json
0. 提交成功
1. 判题队列
2. 编译错误
3. 答案错误
4. 内存超限
5. 时间超限
6. 输出超限
7. 运行错误
8. 判题中
```

```json
{
    "errCode":8
    "msg":   
}

```

```json

{
    input:[input1,input2]
    output:[output1,output2]
}
```



judge

compile

parse

for range{

​	clock()

​	input->output

​	clock()

​	jugde

​	ret

}

ret





Server MQ



Judge1

Judge2

Judge3









{
​    "problemName":
​    "problemID":
​    "timeLimit":
​    "memoryLimit":
​    "property":
}