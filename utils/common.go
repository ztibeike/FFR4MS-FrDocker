package utils

import (
	"frdocker/db"
	"log"
	"os"
)

var logger = log.New(os.Stderr, "", 0)

var containerMgo = db.GetContainerMgo()
var networkMgo = db.GetNetworkMgo()
