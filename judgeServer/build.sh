docker rm  judgeServer
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' . && docker build . -t gojudge
docker run -it  --net=gojudge --name=judgeServer -p 127.0.0.1:7070:8080  gojudge
