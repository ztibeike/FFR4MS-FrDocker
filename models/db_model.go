package models

import "frdocker/types"

type NetWork struct {
	Id   string `bson:"id"`
	Name string `bson:"name"`
}

type Container struct {
	Container *types.Container `bson:"container"`
	NetworkId string           `bson:"networkId"`
}
