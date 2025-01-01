package utils

import (
	"fmt"
	"github.com/malfunkt/iprange"
	"net"
	"strconv"
	"strings"
	"time"
)

func Connect(ip string, port int) (string, int, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ip, port), 2*time.Second)

	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()

	return ip, port, err
}

func GetIpList(ips string) ([]net.IP, error) {
	addressList, err := iprange.ParseList(ips)
	if err != nil {
		return nil, err
	}
	list := addressList.Expand()
	return list, err
}

func GetPortList(selection string) ([]int, error) {
	ports := []int{}
	if selection == "" {
		return ports, nil
	}
	ranges := strings.Split(selection, ",")
	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if strings.Contains(r, "-") {
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("Invalid port selection segment: '%s'", r)
			}
			p1, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("Invalid port number: '%s'", parts[0])
			}
			p2, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("Invalid port number: '%s'", parts[1])
			}
			if p1 > p2 {
				return nil, fmt.Errorf("Invalid port range: %d-%d", parts[0], parts[1])
			}
			for i := p1; i <= p2; i++ {
				ports = append(ports, i)
			}
		}else{
			if port, err := strconv.Atoi(r); err != nil{
				return nil, fmt.Errorf("Invalid port number: '%s'", r)
			}else{
				ports = append(ports, port)
			}
		}
	}
	return ports, nil
}