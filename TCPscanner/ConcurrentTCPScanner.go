package main

import (
	"fmt"
	"os"
	utils2 "securitydevelopment/TCPscanner/utils"
)

func main() {
	if len(os.Args) == 3{
		ipList := os.Args[1]
		portList := os.Args[2]
		ips, _ := utils2.GetIpList(ipList)
		ports, _ := utils2.GetPortList(portList)
		task, _ := utils2.GenerateTask(ips, ports)
		utils2.AssigningTask(task)
		utils2.PrintResult()
	}else{
		fmt.Println("Parameter error.")
	}
}