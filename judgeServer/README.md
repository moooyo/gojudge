### 设计

此模块为gojudge的判题任务调度系统。



分为三个子模块

ListenServer：接收web前端的判题请求, 将请求提交给Dispatcher

Dispatcher： 根据规则，启动docker，运行JudgeCore， 将请求交给ProcessServer

ProcessServer： 与 JudgeCore进行交互，监控判题情况并进行处理



关键数据结构



```go


type SubmitTaskStatus int

const (
    OK    		SubmitTaskStatus = 0
    ERROR 		SubmitTaskStatus = 1
    WAITING 	        SubmitTaskStatus = 2
    JUDGING  	        SubmitTaskStatus = 3
)

type SubmitTaskWrap struct {
    Status  SubmitTaskStatus
    Task    *Submit
}
```



SubmitTaskStatus

| Status  | Mean             |
| ------- | ---------------- |
| OK      | 判题结束         |
| ERROR   | 判题过程出错     |
| WAITING | 未派发的判题请求 |
| JUDGING | 正在判题         |



#### ListenServer

监听web前端发来的判题请求，将请求包装成SubmitTaskWrap，置Status为WAITING，通过channel传给Dispatcher

```go
for {
    conn, err := listener.Accept() 
    go func(conn) {
        submit = readSubmit()
        submitTaskWrap = wrapSubmit(submit)
        channel2Dispatcher <- submitTaskWrap
        conn.Close()
    }(conn)
       
}
```



#### Dispatcher

1. 监听ListenServer发来的submit，将submit存至队列或直接派发

2. 根据规则派发给ProcessServer，等待

3. 监听ProcessServer发回的判题请求，进行下一次任务派发


```go
for {
    submitTaskWrap := <- channel2Dispatcher
    if submitTaskWrap.Status == WAITING {
        // submitTaskWrap from ListenServer
        if can not direct dispatch {
        	store submitTaskWrap in queue and continue wait
    	} 
    } else {
        // submitTaskWrap from ProcessServer
        if submitTaskWrap.Status == OK {
            do something
        } else submitTaskWrap.Status == ERROR {
            do something
        }
        update dispatcher info like memory and docker count
    }
    if can dispatch {
        get submitTaskWrap from queue
        channel2ProcessServer <- submitTaskWrap // or just add submitTaskWrap to ProcessServer's hashMap
        update dispatcher info and set the args
        run docker args
    }
}
```



#### ProcessServer

1. 监听JudgeCore
2. 与judgeCore交互
3. 更新SubmitTaskWrap，返回给Dispatcher

```go
for {
    conn, err = listener.Accpet()
    go func(conn) {
        submitId := readSubmitID()
        submitTaskWrap := getSubmitTaskWrapByID(submitID)
        {
            do something
        }
        update submitTaskWrap.Status
        channel2Dispatcher <- submitTaskWrap
        conn.Close()
    }(conn)
}
```

