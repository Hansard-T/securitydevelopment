package main

import (
	"fmt"
	"os"
	"securitydevelopment/utils"
)

func main() {
	if len(os.Args) == 3{
		ipList := os.Args[1]
		portList := os.Args[2]
		ips, _ := utils.GetIpList(ipList)
		ports, _ := utils.GetPortList(portList)
		task, _ := utils.GenerateTask(ips, ports)
		utils.AssigningTask(task)
		utils.PrintResult()
	}else{
		fmt.Println("Parameter error.")
	}
}