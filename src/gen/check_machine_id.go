package gen

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/nextbin/go-id-gen/src/base"
	log "github.com/sirupsen/logrus"
	"net"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func CheckMachineId() (valid bool, err error) {
	valid = true
	if base.MachineId < 0 || base.MachineId > (1<<MachineIdBit) {
		valid = false
		log.WithField("MachineId", base.MachineId).Warn("MachineId error")
		return
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil || len(addrs) == 0 {
		log.Warn("addrs error. ", err)
		return false, err
	}
	var localIp = ""
	for _, addr := range addrs {
		if !strings.HasPrefix(addr.String(), "127.0.0.1") {
			matched, err := regexp.MatchString("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}", addr.String())
			if err != nil {
				log.Warn("regexp error")
				return false, err
			}
			if matched {
				localIp = addr.String()
				break
			}
		}
	}
	if localIp == "" {
		return false, errors.New("cannot get local IP")
	}
	log.WithFields(log.Fields{"localIp": localIp, "checkType": base.CheckMachineIdType}).Info("check MachineId")
	switch base.CheckMachineIdType {
	case base.CheckMachineIdTypeMysql:
		return CheckMachineIdByMysql(localIp)
	case base.CheckMachineIdTypeRedis:
		return CheckMachineIdByRedis(localIp)
	case base.CheckMachineIdTypeRedisSentinel:
		return false, errors.New("Reids sentinel is NOT supported. ")
	case base.CheckMachineIdTypeMongo:
		return false, errors.New("Mongo DB is NOT supported. ")
	case base.CheckMachineIdTypeNever:
		return
	}
	return
}

func CheckMachineIdByMysql(localIp string) (valid bool, err error) {
	db, err := sql.Open("mysql", base.MysqlDataSourceNaming)
	if err != nil {
		log.Warn("Connect to mysql server error", err)
		return false, err
	}
	defer func() {
		err = db.Close()
		if err != nil {
			log.Warn(err)
		}
	}()
	machine := new(base.Machine)
	var rows = db.QueryRow("SELECT * FROM `nb_id_gen_machine` WHERE `ip`=?", localIp)
	err = rows.Scan(&machine.Id, &machine.Ip, &machine.CreateTime)
	if err != nil && err.Error() != "sql: no rows in result set" {
		log.Warn("SELECT error", err)
		return
	}
	if err == nil {
		return machine.Id == base.MachineId, err
	}
	res, err := db.Exec("INSERT INTO `nb_id_gen_machine`(`id`,`ip`,`create_time`) VALUES (?,?,?)", base.MachineId, localIp, time.Now())
	if err != nil {
		log.Warn("INSERT error", err)
		return
	}
	insertId, err := res.LastInsertId()
	if err != nil {
		log.Warn("LastInsertId error", err)
		return
	}
	if int32(insertId) == base.MachineId {
		return true, err
	}
	return
}

func CheckMachineIdByRedis(localIp string) (valid bool, err error) {
	var c redis.Conn
	c, err = redis.Dial("tcp", base.RedisAddr)
	if err != nil {
		log.Warn("Connect to redis server error", err)
		return
	}
	defer func() {
		err = c.Close()
		if err != nil {
			log.Warn(err)
		}
	}()
	res, err := c.Do("HGET", base.RedisCheckMachineIdKey, localIp)
	if err != nil {
		log.Warn("HGET error", err)
		return
	}
	if res != nil {
		var realRes int32 = 0
		switch v := res.(type) {
		case []uint8:
			for _, digit := range v {
				realRes = realRes*10 + int32(digit-'0')
			}
			log.WithFields(log.Fields{"found": realRes, "config": base.MachineId, "ip": localIp}).Info("MachineId compare")
			return realRes == base.MachineId, err
		}
		log.Warnf("Find redis type: %s", reflect.TypeOf(res))
		return res == base.MachineId, err
	}
	res, err = c.Do("HSETNX", base.RedisCheckMachineIdKey, localIp, base.MachineId)
	if err != nil {
		log.Warn("HSETNX error", err)
		return
	}
	switch v := res.(type) {
	case int64:
		valid = int(v) == 1
		return
	}
	log.Warnf("Find redis type: %s", reflect.TypeOf(res))
	return
}
