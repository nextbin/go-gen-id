package base

import (
	"github.com/nextbin/go-gen-id/src/domain/enum"
)

const (
	/**
	【重要】机器ID，部署ID生成服务时需要修改该值
	*/
	MachineId int32 = 0
	/**
	HTTP-Gin 服务端口
	*/
	HttpGinPort = 11001
	/**
	RPC-Grpc 服务端口
	*/
	RpcGrpcPort = 12001
	/**
	服务器类型
	*/
	ServerFlag = enum.ServerFlagHttpGin | enum.ServerFlagRpcGrpc
	/**
	为节省原始snowflake算法的时间戳可用范围，将时间戳减去项目第一次发布的时间戳
	2019-10-13 00:00:00 +0800
	*/
	ZeroTimestamp int64 = 1570896000000
	/**
	检查机器ID是否冲突
	支持的检查方式：MySQL、Redis
	未实现的检查方式：Redis-Sentinel、Mongo
	*/
	CheckMachineIdType enum.CheckMachineIdTypeEnum = enum.CheckMachineIdTypeNever
	/**
	日志文件
	*/
	LogFilePath = "go-gen-id.log"
)
const (
	/**
	Redis const var
	*/
	RedisAddr              = "127.0.0.1:6379"
	RedisCheckMachineIdKey = "nb:go-gen-id:machine-id"
)
const (
	/**
	MySQL const var
	*/
	//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	MysqlDataSourceNaming string = "root:root@tcp(localhost:3306)/go_gen_id?parseTime=true"
)
