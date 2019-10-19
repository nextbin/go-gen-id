package before_start

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/nextbin/go-gen-id/src/base"
	"github.com/nextbin/go-gen-id/src/domain"
	"github.com/nextbin/go-gen-id/src/domain/enum"
	"github.com/nextbin/go-gen-id/src/server/component/gen"
	log "github.com/sirupsen/logrus"
	"net"
	"regexp"
	"strings"
	"time"
)

func CheckMachineId() (valid bool, err error) {
	if !gen.CheckMachineId(base.MachineId) {
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
	case enum.CheckMachineIdTypeMysql:
		return checkMachineIdByMysql(localIp)
	case enum.CheckMachineIdTypeRedis:
		return checkMachineIdByRedis(localIp)
	case enum.CheckMachineIdTypeRedisSentinel:
		return false, errors.New("Reids sentinel is NOT supported. ")
	case enum.CheckMachineIdTypeMongo:
		return false, errors.New("Mongo DB is NOT supported. ")
	case enum.CheckMachineIdTypeNever:
		return true, err
	}
	return false, err
}

func checkMachineIdByMysql(localIp string) (valid bool, err error) {
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
	machine := new(domain.Machine)
	var rows = db.QueryRow("SELECT * FROM `nb_gen_id_machine` WHERE `ip`=?", localIp)
	err = rows.Scan(&machine.Id, &machine.Ip, &machine.CreateTime)
	if err != nil && err.Error() != "sql: no rows in result set" {
		log.Fatal("SELECT error", err)
	}
	if err == nil {
		log.WithFields(log.Fields{"found": machine.Id, "config": base.MachineId, "ip": localIp}).Info("MachineId compare")
		return machine.Id == base.MachineId, err
	}
	res, err := db.Exec("INSERT INTO `nb_gen_id_machine`(`id`,`ip`,`create_time`) VALUES (?,?,?)", base.MachineId, localIp, time.Now())
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

func checkMachineIdByRedis(localIp string) (valid bool, err error) {
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
		var realRes int32
		for _, digit := range res.([]uint8) {
			realRes = realRes*10 + int32(digit-'0')
			log.WithFields(log.Fields{"found": realRes, "config": base.MachineId, "ip": localIp}).Info("MachineId compare")
		}
		return realRes == base.MachineId, err
	}
	res, err = c.Do("HSETNX", base.RedisCheckMachineIdKey, localIp, base.MachineId)
	if err != nil {
		log.Fatal("HSETNX error", err)
	}
	valid = int(res.(int64)) == 1
	return
}
