package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"time"
)

var (
	device = "en0"
	snapshotLength int32 = 1024
	promiscuous = false
	timeout = 30 * time.Second
)

func Protocal_Analysis(packet gopacket.Packet) {
	alllayers := packet.Layers()

	for _, layer := range alllayers {
		fmt.Printf("Layer : %v\n", layer.LayerType())
	}
	fmt.Println("----------------------------------------------------------------")

	ethernet := packet.Layer(layers.LayerTypeEthernet)
	if ethernet != nil {
		ethernetPacket, _ := ethernet.(*layers.Ethernet)
		fmt.Printf("Ethernet Type: %v, Source MAC: %v, Destination MAC: %v\n", ethernetPacket.EthernetType, ethernetPacket.SrcMAC, ethernetPacket.DstMAC)
		fmt.Println("----------------------------------------------------------------")
	}

	ip := packet.Layer(layers.LayerTypeIPv4)
	if ip != nil {
		ipv4, _ := ip.(*layers.IPv4)
		fmt.Printf("proto: %v, from: %v, to: %v\n", ipv4.Protocol, ipv4.SrcIP, ipv4.DstIP)
		fmt.Println("----------------------------------------------------------------")
	}

	tcp := packet.Layer(layers.LayerTypeTCP)
	if tcp != nil {
		tcp1, _ := tcp.(*layers.TCP)
		fmt.Printf("source port: %v, dest Port: %v\n", tcp1.SrcPort, tcp1.DstPort)
		fmt.Println("----------------------------------------------------------------")
	}

	udp := packet.Layer(layers.LayerTypeUDP)
	if udp != nil {
		udp1 := udp.(*layers.UDP)
		fmt.Printf("src port: %v, dst port: %v\n", udp1.SrcPort, udp1.DstPort)
		fmt.Println("----------------------------------------------------------------")
	}

	app := packet.ApplicationLayer()
	if app != nil {
		fmt.Printf("application payload: %v\n", string(app.Payload()))
	}

	err := packet.ErrorLayer()
	if err != nil {
		fmt.Printf("decode packet err: %v\n", err)
	}
}

func main() {
	handle, err := pcap.OpenLive(device, snapshotLength, promiscuous, timeout)

	if err != nil {
		fmt.Println(err)
	}

	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		Protocal_Analysis(packet)
	}
}