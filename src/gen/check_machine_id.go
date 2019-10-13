package gen

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/nextbin/go-id-gen/src/base"
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
		fmt.Println("MachineId error")
		return
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil || len(addrs) == 0 {
		fmt.Println("addrs error")
		return false, err
	}
	var localIpBuffer = bytes.Buffer{}
	for _, addr := range addrs {
		if !strings.HasPrefix(addr.String(), "127.0.0.1") {
			var matched bool
			matched, err = regexp.MatchString("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}", addr.String())
			if matched {
				localIpBuffer.WriteString(addr.String())
				localIpBuffer.WriteString(",")
			}
		}
	}
	var localIp = localIpBuffer.String()[:localIpBuffer.Len()-1]
	fmt.Printf("Local IP: %s, check machine id type: %d\n", localIp, base.CheckMachineIdType)
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
		fmt.Println("Connect to mysql server error", err)
		return false, err
	}
	defer func() {
		err = db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	machine := new(base.Machine)
	var rows = db.QueryRow("SELECT * FROM `nb_id_gen_machine` WHERE `ip`=?", localIp)
	err = rows.Scan(&machine.Id, &machine.Ip, &machine.CreateTime)
	if err != nil && err.Error() != "sql: no rows in result set" {
		fmt.Println("SELECT error", err)
		return
	}
	if err == nil {
		return machine.Id == base.MachineId, err
	}
	res, err := db.Exec("INSERT INTO `nb_id_gen_machine`(`id`,`ip`,`create_time`) VALUES (?,?,?)", base.MachineId, localIp, time.Now())
	if err != nil {
		fmt.Println("INSERT error", err)
		return
	}
	insertId, err := res.LastInsertId()
	if err != nil {
		fmt.Println("LastInsertId error", err)
		return
	}
	if int32(insertId) == base.MachineId {
		return true, err
	}
	fmt.Println(insertId)
	fmt.Println(base.MachineId)
	return
}

func CheckMachineIdByRedis(localIp string) (valid bool, err error) {
	var c redis.Conn
	c, err = redis.Dial("tcp", base.RedisAddr)
	if err != nil {
		fmt.Println("Connect to redis server error", err)
		return
	}
	defer func() {
		err = c.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	res, err := c.Do("HGET", base.RedisCheckMachineIdKey, localIp)
	if err != nil {
		fmt.Println("HGET error", err)
		return
	}
	if res != nil {
		var realRes int32 = 0
		switch v := res.(type) {
		case []uint8:
			for _, digit := range v {
				realRes = realRes*10 + int32(digit-'0')
			}
			fmt.Printf("MachineId found in redis: %d, config: %d, IP: %s\n", realRes, base.MachineId, localIp)
			return realRes == base.MachineId, err
		}
		fmt.Printf("Find redis type: %s\n", reflect.TypeOf(res))
		return res == base.MachineId, err
	}
	res, err = c.Do("HSETNX", base.RedisCheckMachineIdKey, localIp, base.MachineId)
	if err != nil {
		fmt.Println("HSETNX error", err)
		return
	}
	switch v := res.(type) {
	case int64:
		valid = int(v) == 1
		return
	}
	fmt.Printf("Find redis type: %s\n", reflect.TypeOf(res))
	return
}
