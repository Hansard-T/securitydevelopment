package modules

import (
	"fmt"
	"strings"
	"time"

	"securitydevelopment/Sniffer/WebSpy/logger"
	"securitydevelopment/Sniffer/WebSpy/modules/arpspoof"
	"securitydevelopment/Sniffer/WebSpy/modules/assembly"
	"securitydevelopment/Sniffer/WebSpy/modules/web"
	"securitydevelopment/Sniffer/WebSpy/vars"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	snapshotLen int32 = 10240
	promiscuous bool  = true
	timeout     time.Duration = pcap.BlockForever
	handle      *pcap.Handle
	err         error
	DebugMode  bool
	DeviceName = "eth0"
	filter     = ""
	Mode       = "local"
)

func Start(ctx *cli.Context) error {
	if ctx.IsSet("device") {
		DeviceName = ctx.String("device")
	}

	if ctx.IsSet("mode") {
		Mode = ctx.String("mode")
	}

	if ctx.IsSet("host") {
		vars.HttpHost = ctx.String("host")
	}

	if ctx.IsSet("port") {
		vars.HttpPort = ctx.Int("port")
	}

	if ctx.IsSet("debug") {
		DebugMode = ctx.Bool("debug")
	}
	if DebugMode {
		logger.Log.Logger.Level = logrus.DebugLevel
	}

	if ctx.IsSet("length") {
		snapshotLen = int32(ctx.Int("len"))
	}

	handle, err = pcap.OpenLive(DeviceName, snapshotLen, promiscuous, timeout)
	if err != nil {
		logger.Log.Fatal(err)
	}
	defer handle.Close()

	if ctx.IsSet("filter") {
		filter = ctx.String("filter")
	}

	err = handle.SetBPFFilter(filter)
	if err != nil {
		return err
	}

	go web.Serve(fmt.Sprintf("%v:%v", vars.HttpHost, vars.HttpPort))

	if strings.ToLower(Mode) == "local" {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		assembly.ProcessPackets(packetSource.Packets())
	} else {
		target := ""
		if ctx.IsSet("target") {
			target = ctx.String("target")
		}

		gateway := ""
		if ctx.IsSet("gateway") {
			gateway = ctx.String("gateway")
		}

		if target != "" && gateway != "" {
			logger.Log.Infof("start arpspoof, device: %v, target: %v, gateway:%v", DeviceName, target, gateway)
			go arpspoof.ArpSpoof(DeviceName, handle, target, gateway)

			packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
			assembly.ProcessPackets(packetSource.Packets())
		} else {
			logger.Log.Info("Need to provide target and gateway parameters")
		}
	}

	return err
}