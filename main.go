package main

import "github.com/dhiemaz/bank-api/cmd"

func main() {
	command := cmd.NewCommand()
	command.Run()
}
