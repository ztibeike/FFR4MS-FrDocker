package entity

import (
	"errors"
	"strings"
	"time"

	"gitee.com/zengtao321/frdocker/config"
	"gitee.com/zengtao321/frdocker/utils"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type HttpType int

const (

	// http类型: 无效
	HTTP_INVALID HttpType = iota

	// http类型: 请求
	HTTP_REQUEST

	// http类型: 响应
	HTTP_RESPONSE
)

type HttpInfo struct {
	Type      HttpType
	URL       string
	Src       HttpRole
	Dst       HttpRole
	TraceId   string
	Timestamp time.Time
}

type HttpRole struct {
	IP   string
	Port int
	Type ContainerType
	Name string
}

func NewHttpInfo(packet gopacket.Packet) (*HttpInfo, error) {
	// 检查packet有效性
	if packet == nil || packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
		return nil, errors.New("invalid packet")
	}
	tcp, ok := packet.TransportLayer().(*layers.TCP)
	if !ok || len(tcp.Payload) < 16 {
		return nil, errors.New("invalid packet")
	}
	srcIP, dstIP, srcPort, dstPort := utils.GetIPAndPort(packet)
	httpInfo := &HttpInfo{
		Src: HttpRole{
			IP:   srcIP,
			Port: srcPort,
		},
		Dst: HttpRole{
			IP:   dstIP,
			Port: dstPort,
		},
		Timestamp: packet.Metadata().Timestamp,
	}
	err1 := httpInfo.setTraceId(tcp.Payload)
	err2 := httpInfo.setHttpType(tcp.Payload)
	var err3 error
	if err2 != nil && httpInfo.Type == HTTP_REQUEST {
		err3 = httpInfo.setURL(tcp.Payload)
	}
	if err1 != nil || err2 != nil || err3 != nil {
		return nil, errors.New("invalid packet")
	}
	return httpInfo, nil
}

func (httpInfo *HttpInfo) setTraceId(payload []byte) error {
	data := string(payload)
	idx := strings.Index(data, config.TRACE_ID_HEADER)
	if idx == -1 {
		return errors.New("invalid payload")
	}
	dataRune := []rune(data)
	var traceId string
	for i := idx; i < len(dataRune); i++ {
		if string(dataRune[i]) == "\r" || string(dataRune[i]) == "\n" {
			break
		}
		traceId += string(dataRune[i])
	}
	httpInfo.TraceId = strings.Split(traceId, ": ")[1]
	return nil
}

func (httpInfo *HttpInfo) setHttpType(payload []byte) error {
	data := string(payload)
	idx := strings.Index(data, "HTTP/1.1")
	if idx == -1 {
		return errors.New("invalid payload")
	} else if idx == 0 {
		httpInfo.Type = HTTP_RESPONSE
	} else {
		httpInfo.Type = HTTP_REQUEST
	}
	return nil
}

func (httpInfo *HttpInfo) setURL(payload []byte) error {
	data := string(payload)
	start := strings.Index(data, "/")
	end := strings.Index(data, "HTTP")
	if start == -1 || end == -1 || start >= end {
		return errors.New("invalid payload")
	}
	httpInfo.URL = data[start : end-1]
	return nil
}
