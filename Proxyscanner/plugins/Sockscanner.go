package plugins

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"h12.io/socks"
)

var (
	SockProxyProtocol = map[string]int{"SOCKS4": socks.SOCKS4, "SOCKS4A": socks.SOCKS4A, "SOCKS5": socks.SOCKS5}
)

func CheckSockProxy(ip string, port int, protocol string) (isProxy bool, proxyInfo ProxyInfo, err error) {
	proxyInfo.Addr = ip
	proxyInfo.Port = port
	proxyInfo.Protocol = protocol

	proxy := fmt.Sprintf("%v:%v", ip, port)
	dialSocksProxy := socks.DialSocksProxy(SockProxyProtocol[protocol], proxy)
	tr := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialSocksProxy(network, addr)
		},
	}

	httpClient := &http.Client{Transport: tr, Timeout: 10 * time.Second}

	resp, err := httpClient.Get(WebUrl)

	if err != nil {
		return false, proxyInfo, err
	}

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return false, proxyInfo, err
		}
		if strings.Contains(string(body), "163网易免费邮-你的专业电子邮局") {
			isProxy = true
		}
	}
	return isProxy, proxyInfo, err
}