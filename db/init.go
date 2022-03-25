package db

var (
	NetworkMgo   *Mgo
	ContainerMgo *Mgo
	UserMongo    *Mgo
)

func init() {
	NetworkMgo = NewMgo("network")
	ContainerMgo = NewMgo("container")
	UserMongo = NewMgo("user")
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

func GetUserMgo() *Mgo {
	if UserMongo == nil {
		UserMongo = NewMgo("user")
	}
	return UserMongo
}
