package frecovery

import (
	"frdocker/db"

	"github.com/google/gopacket/pcap"
)

// var logger = log.New(os.Stderr, "", 0)

var pcapHandler *pcap.Handle

var containerMgo = db.GetContainerMgo()
var networkMgo = db.GetNetworkMgo()
var trafficMgo = db.GetTrafficMgo()
