package db

var (
	NetworkMgo   *Mgo
	ContainerMgo *Mgo
	UserMongo    *Mgo
	TrafficMgo   *Mgo
	ErrorLogMgo  *Mgo
)

func init() {
	NetworkMgo = NewMgo("network")
	ContainerMgo = NewMgo("container")
	UserMongo = NewMgo("user")
	TrafficMgo = NewMgo("traffic")
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

func GetTrafficMgo() *Mgo {
	if TrafficMgo == nil {
		TrafficMgo = NewMgo("traffic")
	}
	return TrafficMgo
}

func GetErrorLogMgo() *Mgo {
	if ErrorLogMgo == nil {
		ErrorLogMgo = NewMgo("errorlog")
		ErrorLogMgo.CreateIndex("network")
		ErrorLogMgo.CreateIndex("id")
	}
	return ErrorLogMgo
}
