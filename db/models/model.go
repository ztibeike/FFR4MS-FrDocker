package models

import "frdocker/types"

type NetWork struct {
	Id   int    `bson:"id"`
	Name string `bson:"name"`
}

type Container struct {
	Container *types.Container `bson:"container"`
	NetWorkId int              `bson:"netWorkId"`
}
