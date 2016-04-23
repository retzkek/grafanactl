package main

import (
	"flag"

	log "github.com/Sirupsen/logrus"
	gapi "github.com/retzkek/go-grafana-api"
)

// Command represents a sub-command.
// Inspired by sub-command implementation in Kyle Lemons' rx
// package (http://github.com/kylelemons/rx)
type Command struct {
	Name    string
	Usage   string
	Summary string
	Help    string

	Flag     flag.FlagSet
	Function func(*gapi.Client, *Command, []string) error
}

// Run sets up the environment then executes the command.
func (c *Command) Run(client *gapi.Client, args []string) {
	c.Flag.Parse(args)

	if err := c.Function(client, c, c.Flag.Args()); err != nil {
		log.Fatal(err)
	}
}
