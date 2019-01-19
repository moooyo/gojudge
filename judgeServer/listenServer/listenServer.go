package listenServer

import (
	"github.com/ferriciron/gojudge/def"
	"github.com/ferriciron/gojudge/judgeServer/submitwrap"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type ListenServerConfig struct {
	ListenAddr string `json:"listenAddr"`
}

type ListenServer struct {
	addr              string
	server            *gin.Engine
	dispatcherChannel chan<- submitwrap.SubmitTaskWrap
}

func NewListenServer(config ListenServerConfig, dispatcherChannel chan<- submitwrap.SubmitTaskWrap) *ListenServer {
	gin.SetMode(gin.ReleaseMode)
	return &ListenServer{
		addr:              config.ListenAddr,
		dispatcherChannel: dispatcherChannel,
		server:            gin.Default(),
	}
}

func (listenServer *ListenServer) Run() {
	listenServer.server.POST("/submit_task", func(c *gin.Context) {
		var submit def.Submit
		c.BindJSON(&submit)
		log.Println(&submit)
		listenServer.dispatcherChannel <- submitwrap.WrapSubmit(&submit)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	listenServer.server.Run(listenServer.addr)
}
