package service

import "time"
import "sync"
import "github.com/nextbin/go-gen-id/src/base"

/**
https://github.com/twitter-archive/snowflake/tree/snowflake-2010#solution

Solution
Thrift Server written in Scala
id is composed of:
time - 41 bits (millisecond precision w/ a custom epoch gives us 69 years)
configured machine id - 10 bits - gives us up to 1024 machines
sequence number - 12 bits - rolls over every 4096 per machine (with protection to avoid rollover in the same ms)
System Clock Dependency
You should use NTP to keep your system clock accurate. Snowflake protects from non-monotonic clocks, i.e. clocks that run backwards. If your clock is running fast and NTP tells it to repeat a few milliseconds, snowflake will refuse to generate ids until a time that is after the last time we generated an id. Even better, run in a mode where ntp won't move the clock backwards. See http://wiki.dovecot.org/TimeMovedBackwards#Time_synchronization for tips on how to do this.

在此基础上稍作优化，把时间戳减去项目创建时间，使得gen-id可以用多49年
*/

const (
	TimeBit         uint = 41
	MachineIdBit    uint = 10
	SeqNumBit       uint = 12
	timeOffset      uint = MachineIdBit + SeqNumBit
	machineIdOffset uint = SeqNumBit
)

var lastTime int64 = 0
var seqNum int64 = 0
var seqNumLock sync.Mutex

func GenId() (id int64, code int32, err error) {
	//if base.MachineId < 0 || base.MachineId >= (1<<MachineIdBit) {
	//	return base.CODE_MACHINE_ID_OVERFLOW
	//}
	var now = time.Now()
	var nowTime = now.Unix()*1000 + int64(now.Nanosecond()/1e6) - base.ZeroTime
	if nowTime >= (1 << TimeBit) {
		//时间戳超过 41 bit
		return -1, base.CODE_TIME_OVERFLOW, nil
	}
	seqNumLock.Lock()
	if nowTime != lastTime {
		seqNum = 0
		lastTime = nowTime
	} else {
		seqNum++
	}
	var curSeqNum = seqNum
	seqNumLock.Unlock()
	if seqNum > (1 << SeqNumBit) {
		return -1, base.CODE_SEQ_NUM_OVERFLOW, nil
	}
	return int64(nowTime<<timeOffset) + int64(base.MachineId<<machineIdOffset) + curSeqNum, base.CODE_SUCCESS, nil
}
