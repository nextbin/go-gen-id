package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nextbin/go-id-gen/src/base"
	"github.com/nextbin/go-id-gen/src/gen"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true)
	file, err := os.OpenFile(base.LogFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("init log failed", err)
	}
	log.SetOutput(file)
}

func main() {
	valid, err := gen.CheckMachineId()
	if err != nil || !valid {
		log.Fatal("error before startup check, exit now! ", err)
	}
	gin.SetMode(gin.ReleaseMode) //需要 GIN-debug 日志时把该行注释掉，恢复 debug 模式
	r := gin.Default()
	r.GET("/genId", func(c *gin.Context) {
		var id, code = gen.GenId()
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
	log.WithField("port", base.HttpPort).Info("starting server now...")
	err = r.Run("0.0.0.0:" + strconv.Itoa(base.HttpPort))
	if err != nil {
		log.Fatal(err)
	}
}
