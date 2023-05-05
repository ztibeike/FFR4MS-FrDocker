package frecovery

import (
	"errors"

	"gitee.com/zengtao321/frdocker/frecovery/entity"
	"gitee.com/zengtao321/frdocker/utils"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// 检查packet是否合理
func (app *FrecoveryApp) checkPacketValid(packet gopacket.Packet) bool {
	// 检查packet有效性
	if packet == nil || packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
		return false
	}
	tcp, ok := packet.TransportLayer().(*layers.TCP)
	if !ok || len(tcp.Payload) < 16 {
		return false
	}
	// 检查packet是否是微服务系统内部通信消息, 微服务间消息通信路径为: service->gateway->service
	srcIP, dstIP, srcPort, dstPort := utils.GetIPAndPort(packet)
	srcType := app.getContainerType(utils.GenerateIdFromIPAndPort(srcIP, srcPort))
	dstType := app.getContainerType(utils.GenerateIdFromIPAndPort(dstIP, dstPort))
	if srcType == entity.CTN_INVALID || dstType == entity.CTN_INVALID {
		return false
	}
	return true
}

// 获取容器类型
func (app *FrecoveryApp) getContainerType(id string) entity.ContainerType {
	if _, ok := app.Containers[id]; !ok {
		return entity.CTN_INVALID
	}
	serviceName := app.Containers[id].ServiceName
	if _, ok := app.Services[serviceName]; ok {
		return entity.CTN_SERVICE
	}
	if _, ok := app.Gateways[serviceName]; ok {
		return entity.CTN_GATEWAY
	}
	return entity.CTN_INVALID
}

func (app *FrecoveryApp) setHttpRole(httpInfo *entity.HttpInfo) error {
	srcId := utils.GenerateIdFromIPAndPort(httpInfo.Src.IP, httpInfo.Src.Port)
	dstId := utils.GenerateIdFromIPAndPort(httpInfo.Dst.IP, httpInfo.Dst.Port)
	srcType := app.getContainerType(srcId)
	dstType := app.getContainerType(dstId)
	if srcType == entity.CTN_INVALID || dstType == entity.CTN_INVALID {
		return errors.New("invalid httpInfo")
	}
	httpInfo.Src.Type = srcType
	httpInfo.Dst.Type = dstType
	httpInfo.Src.Name = app.Containers[srcId].ServiceName
	httpInfo.Dst.Name = app.Containers[dstId].ServiceName
	return nil
}
