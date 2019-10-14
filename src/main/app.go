package main

import (
	"github.com/nextbin/go-id-gen/src/base"
	"github.com/nextbin/go-id-gen/src/service"
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
	valid, err := service.CheckMachineId()
	if err != nil || !valid {
		log.Fatal("error before startup check, exit now! ", err)
	}
	if base.ServerFlag <= 0 {
		log.Fatal("ServerFlag <= 0")
	}
	if base.ServerFlag&base.ServerFlagHttpGin > 0 {
		go runHttpGin()
	}
	if base.ServerFlag&base.ServerFlagRpcGrpc > 0 {
		go runRpcGrpc()
	}
	for {
		time.Sleep(time.Hour)
	}
}
