package domain

import "time"

type Machine struct {
	Id         int32
	Ip         string
	CreateTime time.Time
}

type Whitelist struct {
	Id         int32
	Ip         string
	CreateTime time.Time
	Status     int32
}
