package frecovery

import (
	"frdocker/db"

	"github.com/google/gopacket/pcap"
)

var handler *pcap.Handle

var containerMgo = db.GetContainerMgo()
var networkMgo = db.GetNetworkMgo()
