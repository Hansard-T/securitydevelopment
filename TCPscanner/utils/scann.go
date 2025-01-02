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
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), 2*time.Second)
	if conn != nil {
		defer conn.Close()
	}
	return ip, port, err
}

func GetIpList(ips string) ([]net.IP, error) {
	addressList, err := iprange.ParseList(ips)
	if err != nil {
		return nil, err
	}
	return addressList.Expand(), nil
}

func GetPortList(selection string) ([]int, error) {
	if selection == "" {
		return nil, nil
	}

	var ports []int
	for _, segment := range strings.Split(selection, ",") {
		segment = strings.TrimSpace(segment)
		if err := processSegment(segment, &ports); err != nil {
			return nil, err
		}
	}

	return ports, nil
}

func processSegment(segment string, ports *[]int) error {
	if strings.Contains(segment, "-") {
		return parsePortRange(segment, ports)
	}
	return parseSinglePort(segment, ports)
}

func parsePortRange(segment string, ports *[]int) error {
	parts := strings.Split(segment, "-")
	if len(parts) != 2 {
		return fmt.Errorf("Invalid port selection segment: '%s'", segment)
	}

	p1, err1 := strconv.Atoi(parts[0])
	p2, err2 := strconv.Atoi(parts[1])

	if err1 != nil {
		return fmt.Errorf("Invalid port number: '%s'", parts[0])
	}
	if err2 != nil {
		return fmt.Errorf("Invalid port number: '%s'", parts[1])
	}
	if p1 > p2 {
		return fmt.Errorf("Invalid port range: %d-%d", p1, p2)
	}

	for i := p1; i <= p2; i++ {
		*ports = append(*ports, i)
	}
	return nil
}

func parseSinglePort(segment string, ports *[]int) error {
	port, err := strconv.Atoi(segment)
	if err != nil {
		return fmt.Errorf("Invalid port number: '%s'", segment)
	}
	*ports = append(*ports, port)
	return nil
}