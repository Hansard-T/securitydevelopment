package main

import (
	"os"
	"runtime"
	"securitydevelopment/Sniffer/WebSpy/cmd"

	"github.com/urfave/cli"


)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{cmd.Start}
	app.Flags = append(app.Flags, cmd.Start.Flags...)
	_ = app.Run(os.Args)
}