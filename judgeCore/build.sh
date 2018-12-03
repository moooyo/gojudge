CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' . && docker build . -t gojudgecore
#docker run -it  --net=maymomo --name=judgeServer -p 127.0.0.1:8080:8080 -p 127.0.0.1:8081:8081 gojudge
