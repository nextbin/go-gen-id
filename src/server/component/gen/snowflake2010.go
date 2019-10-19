package gen

import (
	log "github.com/sirupsen/logrus"
	"time"
)
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
	timeBit         uint = 41
	machineIdBit    uint = 10
	seqNumBit       uint = 12
	timeOffset      uint = machineIdBit + seqNumBit
	machineIdOffset uint = seqNumBit
	seqNumOffset    uint = 0
	//为避免类型转换导致性能下降，machineIdPureMask、seqNumPureMask、machineIdValue 也设置为 int64
	timePureMask      int64 = (1 << timeBit) - 1
	machineIdPureMask int64 = (1 << machineIdBit) - 1
	seqNumPureMask    int64 = (1 << seqNumBit) - 1
	machineIdValue    int64 = int64(base.MachineId) << machineIdOffset
)

var lastTime int64 = 0
var seqNum int64 = 0
var seqNumLock sync.Mutex
var zeroTime time.Time

func init() {
	zeroTime = time.Unix(base.ZeroTimestamp/1e3, base.ZeroTimestamp%1e3*1e6)
}

func CheckMachineId(machineId int32) (valid bool) {
	//machineIdPureMask 实际上使用 int32，这里仅在启动时调用，只发生 1 次类型转换
	if base.MachineId < 0 || base.MachineId > int32(machineIdPureMask) {
		return false
	}
	return true
}

func GenId() (id int64, code int32, err error) {
	var nowTime = time.Since(zeroTime).Milliseconds()
	if nowTime >= timePureMask {
		//时间戳超过 41 bit
		return -1, base.CODE_TIME_OVERFLOW, nil
	}
	//未设置过lasTime时，假定机器时钟正常，无偏移量
	var timeDiff int64
	timeDiff = nowTime - lastTime
	if timeDiff < -5000 {
		//时钟差距超过5秒，认为时钟不同步
		log.WithField("timeDiff", timeDiff).Warn("timeDiff is too large")
	}
	seqNumLock.Lock()
	if timeDiff < 0 {
		//时钟延迟
		log.WithFields(log.Fields{"timeDiff": timeDiff, "now": nowTime, "last": lastTime}).Warn("time clock delay")
		for timeDiff < 0 {
			nowTime = time.Since(zeroTime).Milliseconds()
			timeDiff = nowTime - lastTime
		}
	}
	if timeDiff == 0 {
		seqNum++
		for seqNum > seqNumPureMask {
			return -1, base.CODE_SEQ_NUM_OVERFLOW, nil
		}
	}
	if timeDiff > 0 {
		lastTime = nowTime
		seqNum = 0
	}
	id = int64(nowTime<<timeOffset) | machineIdValue | seqNum
	seqNumLock.Unlock()
	return id, base.CODE_SUCCESS, nil
}
