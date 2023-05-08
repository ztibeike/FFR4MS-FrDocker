package utils

import (
	"fmt"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func PathJoin(paths ...string) string {
	var absPath string
	for _, path := range paths {
		absPath += string(os.PathSeparator) + path
	}
	return absPath
}

func GenerateContainerId(ip string, port int) string {
	// return ip + ":" + strconv.Itoa(port) // 由于网关转发时采用随机端口，所以生成id时不能依赖端口
	return ip
}

func GenerateStateId(containerId string, api string, src string, dst string) string {
	return fmt.Sprintf("%s:%s:%s:%s", containerId, api, src, dst)
}

// 从packet中获取src和dst的IP的Port
func GetIPAndPort(packet gopacket.Packet) (string, string, int, int) {
	srcIP := packet.NetworkLayer().NetworkFlow().Src().String()
	dstIP := packet.NetworkLayer().NetworkFlow().Dst().String()
	srcPort := int(packet.TransportLayer().(*layers.TCP).SrcPort)
	dstPort := int(packet.TransportLayer().(*layers.TCP).DstPort)
	return srcIP, dstIP, srcPort, dstPort
}
