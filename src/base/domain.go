package base

import "time"

type Machine struct {
	Id         int32
	Ip         string
	CreateTime time.Time
}
