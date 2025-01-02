package utils

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type ProxyAddr struct {
	IP   string
	Port int
}

func ReadFile (filename string) (sliceProxyAddr []ProxyAddr) {
	proxyFile, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error opening file : %v", err)
	}

	defer proxyFile.Close()

	scanner := bufio.NewScanner(proxyFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		ipports := strings.TrimSpace(scanner.Text())
		if ipports == "" {
			continue
		}
		ipport := strings.Split(ipports, ":")
		ip := ipport[0]
		port, err := strconv.Atoi(ipport[1])
		if err == nil {
			proxyAddr := ProxyAddr{IP: ip, Port: port}
			sliceProxyAddr = append(sliceProxyAddr, proxyAddr)
		}
	}
	return sliceProxyAddr
}