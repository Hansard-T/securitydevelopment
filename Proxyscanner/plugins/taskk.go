package plugins

import (
	"fmt"
	"securitydevelopment/Proxyscanner/utils"
	"sync"
)

type ProxyInfo struct {
	Addr     string
	Port     int
	Protocol string
}

type CheckProxyFunc func(ip string, port int, protocol string) (isProxy bool, proxyInfo ProxyInfo, err error)

var (
	httpProxyFunc CheckProxyFunc = CheckHttpProxy
	sockProxyFunc CheckProxyFunc = CheckSockProxy
	Result sync.Map
)

func SaveProxies(isProxy bool, proxyInfo ProxyInfo, err error) error {
	if err == nil && isProxy {
		k := fmt.Sprintf("%v://%v:%v", proxyInfo.Protocol, proxyInfo.Addr, proxyInfo.Port)
		Result.Store(k, true)
	}

	return err
}

func PrintResult() {
	Result.Range(func(key, value interface{}) bool {
		fmt.Printf("%v\n", key)
		return true
	})
}

func CheckProxy(proxies []utils.ProxyAddr){
	var wg sync.WaitGroup
	wg.Add(len(proxies) * (len(HttpProxyProtocol) + len(SockProxyProtocol)))

	for _, addr := range proxies {
		for _, proto := range HttpProxyProtocol {
			go func(ip string, port int, protocol string) {
				defer wg.Done()
				_ = SaveProxies(httpProxyFunc(ip, port, protocol))
			}(addr.IP, addr.Port, proto)
		}

		for proto := range SockProxyProtocol {
			go func(ip string, port int, protocol string) {
				defer wg.Done()
				_ = SaveProxies(sockProxyFunc(ip, port, protocol))
			}(addr.IP, addr.Port, proto)
		}
	}
	wg.Wait()
}

func GenerateTask(IpList string, ScanNum int, ) (err error) {

	proxyAddrList := utils.ReadFile(IpList)
	proxyNum := len(proxyAddrList)

	scanBatch := proxyNum / ScanNum
	for i := 0; i < scanBatch; i++ {
		proxies := proxyAddrList[i*ScanNum : (i+1)*ScanNum]
		CheckProxy(proxies)
	}
	if proxyNum%ScanNum > 0 {
		proxies := proxyAddrList[ScanNum*scanBatch : proxyNum]
		CheckProxy(proxies)
	}
	PrintResult()

	return err
}