package http_gin

import (
	"github.com/gin-gonic/gin"
	"github.com/nextbin/go-gen-id/src/base"
	"github.com/nextbin/go-gen-id/src/server/component/gen"
	log "github.com/sirupsen/logrus"
	"strconv"
)

const (
	HttpGinServerDefaultMode = gin.ReleaseMode
)

type HttpGinServer struct {
	Port            int
	Mode            string
	MiddlewareTypes []MiddlewareType
}

func genIdHandler(c *gin.Context) {
	var id, code, err = gen.GenId()
	if err != nil {
		log.Warn(err)
	}
	if c.Query("pure") == "1" {
		c.String(200, strconv.FormatInt(id, 10))
	} else {
		c.JSON(200, gin.H{
			"code":    code,
			"data":    id,
			"message": "",
		})
	}
}

func (s *HttpGinServer) Run() {
	gin.SetMode(s.Mode) //需要 GIN-debug 日志时把该行注释掉，恢复 debug 模式
	r := gin.Default()
	for _, v := range s.MiddlewareTypes {
		switch v {
		case MiddlewareLogLatency:
			r.Use(LogLatency())
		case MiddlewareLogResponseStatus:
			r.Use(LogResponseStatus())
		case MiddlewareWhitelist:
			r.Use(Whitelist())
		}
	}
	r.GET("/genId", genIdHandler)
	log.WithField("port", base.HttpGinPort).Info("starting server now...")
	err := r.Run("0.0.0.0:" + strconv.Itoa(s.Port))
	if err != nil {
		log.Fatal(err)
	}
}
