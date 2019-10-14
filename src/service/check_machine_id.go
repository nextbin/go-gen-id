package service

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/nextbin/go-gen-id/src/base"
	log "github.com/sirupsen/logrus"
	"net"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func CheckMachineId() (valid bool, err error) {
	if base.MachineId < 0 || base.MachineId > (1<<MachineIdBit) {
		log.WithField("MachineId", base.MachineId).Fatal("MachineId error")
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil || len(addrs) == 0 {
		log.Fatal("addrs error. ", err)
	}
	var localIp = ""
	for _, addr := range addrs {
		if !strings.HasPrefix(addr.String(), "127.0.0.1") {
			matched, err := regexp.MatchString("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}", addr.String())
			if err != nil {
				log.Fatal("regexp error")
			}
			if matched {
				localIp = addr.String()
				break
			}
		}
	}
	if localIp == "" {
		log.Fatal("cannot get local IP")
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
		return true, err
	}
	return false, err
}

func CheckMachineIdByMysql(localIp string) (valid bool, err error) {
	db, err := sql.Open("mysql", base.MysqlDataSourceNaming)
	if err != nil {
		log.Fatal("Connect to mysql server error", err)
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
		log.Fatal("SELECT error", err)
	}
	if err == nil {
		log.WithFields(log.Fields{"found": machine.Id, "config": base.MachineId, "ip": localIp}).Info("MachineId compare")
		return machine.Id == base.MachineId, err
	}
	res, err := db.Exec("INSERT INTO `nb_id_gen_machine`(`id`,`ip`,`create_time`) VALUES (?,?,?)", base.MachineId, localIp, time.Now())
	if err != nil {
		log.Fatal("INSERT error", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		log.Fatal("LastInsertId error", err)
	}
	valid = affected == 1
	return
}

func CheckMachineIdByRedis(localIp string) (valid bool, err error) {
	var c redis.Conn
	c, err = redis.Dial("tcp", base.RedisAddr)
	if err != nil {
		log.Fatal("Connect to redis server error", err)
	}
	defer func() {
		err = c.Close()
		if err != nil {
			log.Warn(err)
		}
	}()
	res, err := c.Do("HGET", base.RedisCheckMachineIdKey, localIp)
	if err != nil {
		log.Fatal("HGET error", err)
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
		log.Fatalf("Find redis type: %s", reflect.TypeOf(res))
		return res == base.MachineId, err
	}
	res, err = c.Do("HSETNX", base.RedisCheckMachineIdKey, localIp, base.MachineId)
	if err != nil {
		log.Fatal("HSETNX error", err)
	}
	switch v := res.(type) {
	case int64:
		valid = int(v) == 1
		return
	}
	log.Fatalf("Find redis type: %s", reflect.TypeOf(res))
	return
}
