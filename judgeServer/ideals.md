timer
goroutine pool

```json                         judge  
web server -------------->dispatch goroutine-------> exec("docker run ... /bin/bash/judgeCore args")

judgeCore -----------------------> judgeServer

go process(judgeCore)
```
