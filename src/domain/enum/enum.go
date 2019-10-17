package enum

type CheckMachineIdTypeEnum int32

const (
	CheckMachineIdTypeNever         CheckMachineIdTypeEnum = iota
	CheckMachineIdTypeMysql                                = iota
	CheckMachineIdTypeRedis                                = iota
	CheckMachineIdTypeRedisSentinel                        = iota
	CheckMachineIdTypeMongo                                = iota
)

type ServerFlagEnum int32

const (
	ServerFlagHttpPure ServerFlagEnum = 1 << iota
	ServerFlagRpcPure                 = 1 << iota
	ServerFlagHttpGin                 = 1 << iota
	ServerFlagRpcGrpc                 = 1 << iota
)
