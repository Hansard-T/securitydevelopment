package utils

import (
	"fmt"
	"net"
	"securitydevelopment/TCPscanner/vars"
	"sync"
)

func GenerateTask(ipList []net.IP, ports []int) ([]map[string]int, int) {
	tasks := make([]map[string]int, 0, len(ipList)*len(ports))
	for _, ip := range ipList {
		for _, port := range ports {
			tasks = append(tasks, map[string]int{ip.String(): port})
		}
	}
	return tasks, len(tasks)
}

func AssigningTask(tasks []map[string]int) {
	batchSize := 10
	for i := 0; i < len(tasks); i += batchSize {
		end := i + batchSize
		if end > len(tasks) {
			end = len(tasks)
		}
		RunTask(tasks[i:end])
	}
}

func RunTask(tasks []map[string]int) {
	var wg sync.WaitGroup
	wg.Add(len(tasks))

	for _, task := range tasks {
		for ip, port := range task {
			go func(ip string, port int) {
				defer wg.Done()
				_ = SaveResult(Connect(ip, port))
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
	return nil
}

func PrintResult() {
	vars.Result.Range(func(key, value interface{}) bool {
		fmt.Printf("IP: %v\nPorts: %v\n", key, value)
		return true
	})
}