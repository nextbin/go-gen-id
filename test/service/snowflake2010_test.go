package service

import (
	"fmt"
	"github.com/nextbin/go-id-gen/src/service"
	"testing"
	"time"
)

func Test(t *testing.T) {
	fmt.Println(service.GenId())
	fmt.Println(service.GenId())
	fmt.Println(service.GenId())
	time.Sleep(time.Second * 1)
	fmt.Println(service.GenId())
	fmt.Println(service.GenId())
	fmt.Println(service.GenId())
}
