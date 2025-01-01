package main

import (
	"fmt"
	"github.com/malfunkt/iprange"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func Connect(ip string, port int) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ip, port), 2 * time.Second)
	defer  func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()
	return conn, err
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

func main() {
	if len(os.Args) == 3{
		ipList := os.Args[1]
		portList := os.Args[2]
		ips, _ := GetIpList(ipList)
		ports, _ := GetPortList(portList)
		for _, ip := range ips {
			for _, port := range ports {
				_, err := Connect(ip.String(), port)
				if err != nil {
					continue
				}
				fmt.Printf("ip: %s, port: %d is open\n", ip.String(), port)
			}
		}
	}else{
		fmt.Println("Parameter error.")
	}
}