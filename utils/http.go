package utils

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"gitee.com/zengtao321/frdocker/constants"
	"gitee.com/zengtao321/frdocker/types"
	"gitee.com/zengtao321/frdocker/utils/logger"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

type HttpStreamFactory struct{}

// httpStream will handle the actual decoding of http requests.
type HttpStream struct {
	net, transport gopacket.Flow
	r              tcpreader.ReaderStream
}

func (h *HttpStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	hstream := &HttpStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
	}
	go hstream.Run() // Important... we must guarantee that data from the reader stream is read.

	// ReaderStream implements tcpassembly.Stream, so we can return a pointer to it.
	return &hstream.r
}

func (h *HttpStream) Run() {
	buf := bufio.NewReader(&h.r)
	for {
		req, err := http.ReadRequest(buf)
		if err == io.EOF {
			// We must read until we see an EOF... very important!
			return
		} else if err != nil {
			logger.Errorln(nil, "Error reading stream", h.net, h.transport, ":", err)
		} else {
			bodyBytes := tcpreader.DiscardBytesToEOF(req.Body)
			req.Body.Close()
			logger.Errorln(nil, "Received request from stream", h.net, h.transport, ":", req, "with", bodyBytes, "bytes in request body")
		}
	}
}

func GetHttpType(payload []byte) (string, error) {
	data := string(payload)
	idx := strings.Index(data, "HTTP/1.1")
	if idx == -1 {
		return "", errors.New("payload is not http type")
	} else if idx == 0 {
		return "RESPONSE", nil
	} else {
		return "REQUEST", nil
	}
}

func GetTraceId(payload []byte) string {
	data := string(payload)
	idx := strings.Index(data, "tid")
	if idx == -1 {
		return ""
	}
	dataRune := []rune(data)
	var traceId string
	for i := idx; i < len(dataRune); i++ {
		if string(dataRune[i]) == "\r" || string(dataRune[i]) == "\n" {
			break
		}
		traceId += string(dataRune[i])
	}
	return strings.Split(traceId, ": ")[1]
}

func GetURL(payload []byte) string {
	data := string(payload)
	start := strings.Index(data, "/")
	end := strings.Index(data, "HTTP")
	if start == -1 || end == -1 || start >= end {
		return ""
	}
	return data[start : end-1]
}

func GetHttpInfo(packet gopacket.Packet, tcpLayer *layers.TCP) (*types.HttpInfo, error) {
	srcIP := packet.NetworkLayer().NetworkFlow().Src().String()
	dstIP := packet.NetworkLayer().NetworkFlow().Dst().String()
	srcPort := strconv.Itoa(int(tcpLayer.SrcPort))
	dstPort := strconv.Itoa(int(tcpLayer.DstPort))
	payload := tcpLayer.Payload
	httpInfo := &types.HttpInfo{
		SrcIP:   srcIP,
		SrcPort: srcPort,
		DstIP:   dstIP,
		DstPort: dstPort,
	}
	httpType, err := GetHttpType(payload)
	if err != nil {
		return httpInfo, err
	}
	httpInfo.Type = httpType
	httpInfo.TraceId = GetTraceId(payload)
	httpInfo.URL = GetURL(payload)
	httpInfo.Internal = constants.IPAllMSMap.Has(srcIP)
	httpInfo.Timestamp = packet.Metadata().Timestamp
	return httpInfo, nil
}
