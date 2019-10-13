package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nextbin/go-id-gen/src/base"
	"github.com/nextbin/go-id-gen/src/gen"
	"strconv"
)

func main() {
	valid, err := gen.CheckMachineId()
	if err != nil || !valid {
		fmt.Println("error before startup check, exit now!", err)
		return
	}
	//gin.SetMode(gin.ReleaseMode)
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
	fmt.Println("starting server now...")
	err = r.Run("0.0.0.0:" + strconv.Itoa(base.HttpPort))
	if err != nil {
		fmt.Println(err)
	}
}
