package main

import (
	"flag"
	"os"

	log "github.com/Sirupsen/logrus"
	gapi "github.com/retzkek/go-grafana-api"
)

var (
	verbose = flag.Bool("v", false,
		"turn on verbose output")
	url = flag.String("url", "http://play.grafana.org",
		"Grafana base URL (or set GRAFANA_URL)")
	key = flag.String("key", "",
		"Grafana API key (or set GRAFANA_API_KEY)")
)

var commands = []*Command{
	getCmd,
	helpCmd,
}

func findCommand(cmdName string) *Command {
	for _, c := range commands {
		if cmdName == c.Name {
			return c
		}
	}
	return nil
}

func main() {
	flag.Parse()
	args := flag.Args()

	// check environment variables
	if *url == "" {
		if env := os.Getenv("GRAFANA_URL"); env != "" {
			flag.Set("url", env)
		}
	}
	if *key == "" {
		if env := os.Getenv("GRAFANA_API_KEY"); env != "" {
			flag.Set("key", env)
		}
	}

	// setup log
	if *verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// setup client
	client, err := gapi.New(*key, *url)
	if err != nil {
		log.Fatal(err)
	}

	// if no command, print help
	if len(args) == 0 {
		helpFunc(client, nil, args)
		os.Exit(1)
	}

	// find and run command
	cmdName, args := args[0], args[1:]
	if cmd := findCommand(cmdName); cmd != nil {
		cmd.Run(client, args)
	} else {
		log.Fatal("Unknown command. 'help' for usage.")
	}
}
