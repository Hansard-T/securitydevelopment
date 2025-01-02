package plugins

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	HttpProxyProtocol = []string{"http", "https"}
	WebUrl            = "https://mail.163.com/"
)

func CheckHttpProxy(ip string, port int, protocol string) (isProxy bool, proxyInfo ProxyInfo, err error) {
	proxyInfo.Addr = ip
	proxyInfo.Port = port
	proxyInfo.Protocol = protocol

	rawProxyUrl := fmt.Sprintf("%v://%v:%v", protocol, ip, port)
	proxyUrl, err := url.Parse(rawProxyUrl)

	if err != nil {
		return false, proxyInfo, err
	}

	Transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	client := &http.Client{Transport: Transport, Timeout: 10 * time.Second}
	resp, err := client.Get(WebUrl)

	if err != nil {
		return false, proxyInfo, err
	}

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return false, proxyInfo, err
		}

		if strings.Contains(string(body), "<title>163网易免费邮-你的专业电子邮局") {
			isProxy = true
		}
	}

	return isProxy, proxyInfo, err
}