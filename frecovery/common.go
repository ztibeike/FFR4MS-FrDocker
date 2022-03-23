package frecovery

import (
	"frdocker/db"

	"github.com/google/gopacket/pcap"
)

var pcapHandler *pcap.Handle

var containerMgo = db.GetContainerMgo()
var networkMgo = db.GetNetworkMgo()
