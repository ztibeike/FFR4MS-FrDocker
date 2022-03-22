package db

var (
	NetworkMgo   *Mgo
	ContainerMgo *Mgo
)

func init() {
	NetworkMgo = NewMgo("network")
	ContainerMgo = NewMgo("container")
}

func GetNetworkMgo() *Mgo {
	if NetworkMgo == nil {
		NetworkMgo = NewMgo("network")
	}
	return NetworkMgo
}

func GetContainerMgo() *Mgo {
	if ContainerMgo == nil {
		ContainerMgo = NewMgo("container")
	}
	return ContainerMgo
}
