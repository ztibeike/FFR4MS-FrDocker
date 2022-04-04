package models

import "gitee.com/zengtao321/frdocker/types"

type NetWork struct {
	Id   string `bson:"id"`
	Name string `bson:"name"`
}

type Container struct {
	Container *types.Container `bson:"container"`
	Network   string           `bson:"network"`
}

type ContainerTraffic struct {
	Network string     `bson:"network"`
	IP      string     `bson:"ip"`
	Port    string     `bson:"port"`
	Group   string     `bson:"group"`
	Entry   bool       `bson:"entry"`
	Traffic []*Traffic `bson:"traffic"`
}

// for sort
type ContainerTrafficArray []*ContainerTraffic

func (c ContainerTrafficArray) Len() int {
	return len(c)
}

func (c ContainerTrafficArray) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c ContainerTrafficArray) Less(i, j int) bool {
	return len(c[i].Traffic) > len(c[j].Traffic)
}

type Traffic struct {
	Year   int   `bson:"year" json:"year"`
	Month  int   `bson:"month" json:"month"`
	Day    int   `bson:"day" json:"day"`
	Hour   int   `bson:"hour" json:"hour"`
	Minute int   `bson:"minute" json:"minute"`
	Number int64 `bson:"number" json:"number"`
}
