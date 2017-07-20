package main

import (
	"fmt"
	"os"
	"time"

	"github.com/diogomonica/linked-air/command"
	"github.com/urfave/cli"
)

var Commands = []cli.Command{
	{
		Name:   "companies",
		Usage:  "TABLE FILE",
		Action: command.CmdImportCompanies,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "contacts",
		Usage:  "TABLE FILE",
		Action: command.CmdImportContacts,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "gmail-sync",
		Usage:  "gmail-sync",
		Action: command.CmdGmailSync,
		Flags:  []cli.Flag{cli.DurationFlag{Name: "howlong, s", Value: time.Second * 60}}},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
