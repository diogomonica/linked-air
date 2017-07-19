package main

import (
	"fmt"
	"os"

	"github.com/diogomonica/linked-air/command"
	"github.com/urfave/cli"
)

var Commands = []cli.Command{
	{
		Name:   "companies",
		Usage:  "",
		Action: command.CmdImportCompanies,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "contacts",
		Usage:  "",
		Action: command.CmdImportContacts,
		Flags:  []cli.Flag{},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
