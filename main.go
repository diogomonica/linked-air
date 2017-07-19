package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Author = "diogomonica"
	app.Email = "diogo.monica@gmail.com"
	app.Usage = "linked-air import TABLE FILE"

	app.Commands = Commands
	app.CommandNotFound = CommandNotFound

	app.Run(os.Args)
}
