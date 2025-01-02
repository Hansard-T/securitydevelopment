package main

import "securitydevelopment/Proxyscanner/plugins"

func main() {
	plugins.GenerateTask("iplist.txt", 100)
}