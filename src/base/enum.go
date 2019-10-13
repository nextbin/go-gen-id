package base

const (
	CheckMachineIdTypeNever         = iota
	CheckMachineIdTypeMysql         = iota
	CheckMachineIdTypeRedis         = iota
	CheckMachineIdTypeRedisSentinel = iota
	CheckMachineIdTypeMongo         = iota
)
