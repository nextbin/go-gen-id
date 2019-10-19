package gen

import (
	"fmt"
	"github.com/nextbin/go-gen-id/src/base"
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

func TestParseId(t *testing.T) {
	var id int64 = 2407998756487168
	time := (id&(timePureMask<<timeOffset))>>timeOffset + base.ZeroTimestamp
	machineId := id & (machineIdPureMask << machineIdOffset) >> machineIdOffset
	seqNum := id & (seqNumPureMask << seqNumOffset) >> seqNumOffset
	fmt.Println(map[string]int64{"time": time, "machineId": machineId, "seqNum": seqNum})
}
