package assembly

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"securitydevelopment/Sniffer/WebSpy/logger"
	"securitydevelopment/Sniffer/WebSpy/models"
	"securitydevelopment/Sniffer/WebSpy/vars"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

type httpStreamFactory struct{}

type httpStream struct {
	net, transport gopacket.Flow
	r              tcpreader.ReaderStream
}

// 创建一个新的HTTP流
func (h *httpStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	hStream := &httpStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
	}
	go hStream.run() // 启动处理 goroutine
	return &hStream.r
}

// 处理HTTP请求
func (h *httpStream) run() {
	buf := bufio.NewReader(&h.r)
	for {
		if req, err := http.ReadRequest(buf); err != nil {
			if err == io.EOF {
				return
			}
			continue
		} else {
			handleHTTPRequest(req, h)
		}
	}
}

// 处理HTTP请求
func handleHTTPRequest(req *http.Request, h *httpStream) {
	defer req.Body.Close()

	clientIp, dstIp := SplitNet2Ips(h.net)
	srcPort, dstPort := Transport2Ports(h.transport)

	httpReq, err := models.NewHttpReq(req, clientIp, dstIp, dstPort)
	if err != nil {
		logger.Log.Error("创建 HTTP 请求模型错误: ", err)
		return
	}

	logger.Log.Infof("httpReq: %v", httpReq)

	go func() {
		reqInfo := fmt.Sprintf("%v:%v -> %v(%v:%v), %v, %v, %v, %v",
			httpReq.Client, srcPort, httpReq.Host, httpReq.Ip,
			httpReq.Port, httpReq.Method, httpReq.URL,
			httpReq.Header, httpReq.ReqParameters)
		logger.Log.Infof("reqInfo: %v", reqInfo)
		SendHTML(reqInfo)
	}()
}

// 从网络流中提取客户端和主机IP
func SplitNet2Ips(net gopacket.Flow) (client, host string) {
	ips := strings.Split(net.String(), "->")
	if len(ips) == 2 {
		client, host = ips[0], ips[1]
	}
	return client, host
}

// Transport2Ports 从传输流中提取源端口和目的端口
func Transport2Ports(transport gopacket.Flow) (src, dst string) {
	ports := strings.Split(transport.String(), "->")
	if len(ports) == 2 {
		src, dst = ports[0], ports[1]
	}
	return src, dst
}

// 处理接收到的数据包
func ProcessPackets(packets chan gopacket.Packet) {
	streamFactory := &httpStreamFactory{}
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case packet := <-packets:
			if packet == nil {
				return
			}
			if isTCPPacket(packet) {
				tcp := packet.TransportLayer().(*layers.TCP)
				assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)
			}

		case <-ticker.C:
			assembler.FlushOlderThan(time.Now().Add(-20 * time.Second))
		}
	}
}

// 检查数据包是否为TCP
func isTCPPacket(packet gopacket.Packet) bool {
	return packet.NetworkLayer() != nil && packet.TransportLayer() != nil && packet.TransportLayer().LayerType() == layers.LayerTypeTCP
}

// 将请求信息输入共享数据变量
func SendHTML(reqInfo string) {
	vars.Data.Put(reqInfo)
}