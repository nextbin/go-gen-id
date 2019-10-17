package service

import (
	"fmt"
	"github.com/nextbin/go-gen-id/src/server/component/gen"
	"testing"
	"time"
)

func Test(t *testing.T) {
	fmt.Println(gen.GenId())
	fmt.Println(gen.GenId())
	fmt.Println(gen.GenId())
	time.Sleep(time.Second * 1)
	fmt.Println(gen.GenId())
	fmt.Println(gen.GenId())
	fmt.Println(gen.GenId())
}
