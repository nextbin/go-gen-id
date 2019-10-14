package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nextbin/go-gen-id/src/base"
	"github.com/nextbin/go-gen-id/src/service"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func runHttpGin() {
	gin.SetMode(gin.ReleaseMode) //需要 GIN-debug 日志时把该行注释掉，恢复 debug 模式
	r := gin.Default()
	r.GET("/genId", func(c *gin.Context) {
		var id, code, err = service.GenId()
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
	})
	log.WithField("port", base.HttpGinPort).Info("starting server now...")
	err := r.Run("0.0.0.0:" + strconv.Itoa(base.HttpGinPort))
	if err != nil {
		log.Fatal(err)
	}
}
