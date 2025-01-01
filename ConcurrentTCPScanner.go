package main

import (
	"fmt"
	"github.com/malfunkt/iprange"
	"net"
	"os"
	"securitydevelopment/vars"
	"strconv"
	"strings"
	"sync"
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

func GenerateTask(ipList []net.IP, ports[]int) ([]map[string]int, int) {
	tasks := make([]map[string]int, 0)
	for _, ip := range ipList{
		for _, port := range ports{
			ipPort := map[string]int{ip.String():port}
			tasks = append(tasks, ipPort)
		}
	}
	return tasks, len(tasks)
}

func AssigningTask(tasks []map[string]int){
	BatchNum := len(tasks) / 10
	for i := 0; i < BatchNum; i++ {
		curTask := tasks[10 * i : 10 * (i+1)]
		RunTask(curTask)
	}

	if len(tasks) % BatchNum > 0 {
		lastTask := tasks[10*BatchNum:]
		RunTask(lastTask)
	}
}

func RunTask(tasks []map[string]int){
	var wg sync.WaitGroup
	wg.Add(len(tasks))

	for _, task := range tasks{
		for ip, port := range task{
			go func(string, int) {
				err := SaveResult(Connect(ip, port))
				_ = err
				wg.Done()
			}(ip, port)
		}
	}
	wg.Wait()
}

func SaveResult(ip string, port int, err error) error {
	if err != nil {
		return err
	}

	v, ok := vars.Result.Load(ip)
	if ok {
		ports, ok1 := v.([]int)
		if ok1 {
			ports = append(ports, port)
			vars.Result.Store(ip, ports)
		}
	} else {
		ports := make([]int, 0)
		ports = append(ports, port)
		vars.Result.Store(ip, ports)
	}
	return err
}

func PrintResult() {
	vars.Result.Range(func(key, value interface{}) bool {
		fmt.Printf("ip:%v\n", key)
		fmt.Printf("ports: %v\n", value)
		return true
	})
}

func main() {
	if len(os.Args) == 3{
		ipList := os.Args[1]
		portList := os.Args[2]
		ips, _ := GetIpList(ipList)
		ports, _ := GetPortList(portList)
		task, _ := GenerateTask(ips, ports)
		AssigningTask(task)
		PrintResult()
	}else{
		fmt.Println("Parameter error.")
	}
}