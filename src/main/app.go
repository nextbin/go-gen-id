package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nextbin/go-gen-id/src/base"
	"github.com/nextbin/go-gen-id/src/domain/enum"
	"github.com/nextbin/go-gen-id/src/server/component/before_start"
	"github.com/nextbin/go-gen-id/src/server/http/http_gin"
	"github.com/nextbin/go-gen-id/src/server/rpc/rpc_grpc"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
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
	valid, err := before_start.CheckMachineId()
	if err != nil || !valid {
		log.Fatal("error before startup check, exit now! ", err)
	}
	if base.ServerFlag <= 0 {
		log.Fatal("ServerFlag <= 0")
	}
	if base.ServerFlag&enum.ServerFlagHttpGin > 0 {
		middlewareTypes := []http_gin.MiddlewareType{http_gin.MiddlewareLogLatency, http_gin.MiddlewareLogResponseStatus, http_gin.MiddlewareWhitelist}
		//mode := http_gin.HttpGinServerDefaultMode
		mode := gin.DebugMode
		s := http_gin.HttpGinServer{base.HttpGinPort, mode, middlewareTypes}
		go s.Run()
	}
	if base.ServerFlag&enum.ServerFlagRpcGrpc > 0 {
		s := rpc_grpc.RpcGrpcServer{Port: base.RpcGrpcPort}
		go s.Run()
	}
	for {
		time.Sleep(time.Hour)
	}
}
