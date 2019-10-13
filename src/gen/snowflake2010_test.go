package gen

import (
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T) {
	fmt.Println(GenId())
	fmt.Println(GenId())
	fmt.Println(GenId())
	time.Sleep(time.Second * 1)
	fmt.Println(GenId())
	fmt.Println(GenId())
	fmt.Println(GenId())
}
