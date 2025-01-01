package utils

import (
	"fmt"
	"net"
	"securitydevelopment/vars"
	"sync"
)

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