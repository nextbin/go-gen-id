package http_gin

import (
	"github.com/gin-gonic/gin"
	"github.com/nextbin/go-gen-id/src/base"
	"github.com/nextbin/go-gen-id/src/server/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

type MiddlewareType uint32

const (
	MiddlewareLogLatency        MiddlewareType = iota
	MiddlewareLogResponseStatus                = iota
	MiddlewareWhitelist                        = iota
)

func LogLatency() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		t := time.Now()
		// request
		c.Next()
		// before request
		latency := time.Since(t)
		log.WithFields(log.Fields{"fullPath": c.FullPath(), "latency": latency}).Info("log latency")
	}
}

func LogResponseStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		status := c.Writer.Status()
		log.WithFields(log.Fields{"fullPath": c.FullPath(), "status": status}).Info("log status")
	}
}

func Whitelist() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		list := utils.QueryWhitelist()
		permit := false
		for _, v := range list {
			if v.Ip == ip {
				permit = true
				break
			}
		}
		if permit || ip == "::1" || ip == "127.0.0.1" {
			c.Next()
		} else {
			log.WithFields(log.Fields{"fullPath": c.FullPath(), "ip": ip}).Warn("ip is not in whitelist")
			c.JSON(200, gin.H{
				"code": base.CODE_ACCESS_CONTROL_LIST,
				"msg":  "acl whitelist deny",
			})
			c.Abort()
		}
	}
}
