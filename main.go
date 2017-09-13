package main

import (
	"flag"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/retzkek/grafanactl/gapi"
)

const (
	DEFAULT_URL  = "http://play.grafana.org"
	DEFAULT_PATH = "."
)

// build-time vars
var (
	VERSION = "0.1.5"
	REF     = "scratch"
	BUILD   = ""
)

// run-time flags
var (
	verbose = flag.Bool("v", false,
		"turn on verbose output")
	url = flag.String("url", DEFAULT_URL,
		"Grafana base URL (or set GRAFANA_URL)")
	key = flag.String("key", "",
		"Grafana API key (or set GRAFANA_API_KEY)")
	path = flag.String("path", DEFAULT_PATH,
		"path to local dashboard repository (or set GRAFANA_PATH)")
	headers = flag.String("headers", "",
		"Comma-separated list of extra headers to pass, e.g. \"X-User:foo,X-Grafana-Org-Id:1\" (or set GRAFANA_HEADERS)")
)

var commands = []*Command{
	getCmd,
	helpCmd,
	listCmd,
	pushCmd,
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

	// if no command, print help
	if len(args) == 0 {
		helpFunc(nil, nil, args)
		os.Exit(1)
	}
	// check environment variables
	getenv()
	// setup log
	if *verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	// setup client
	client, err := gapi.New(*key, *headers, *url)
	if err != nil {
		log.Fatal(err)
	}
	// find and run command
	cmdName, args := args[0], args[1:]
	if cmd := findCommand(cmdName); cmd != nil {
		cmd.Run(client, args)
	} else {
		log.Fatal("Unknown command. 'help' for usage.")
	}
}

func getenv() {
	if *url == DEFAULT_URL {
		if env := os.Getenv("GRAFANA_URL"); env != "" {
			flag.Set("url", env)
		}
	}
	if *key == "" {
		if env := os.Getenv("GRAFANA_API_KEY"); env != "" {
			flag.Set("key", env)
		}
	}
	if *headers == "" {
		if env := os.Getenv("GRAFANA_HEADERS"); env != "" {
			flag.Set("headers", env)
		}
	}
	if *path == DEFAULT_PATH {
		if env := os.Getenv("GRAFANA_PATH"); env != "" {
			flag.Set("path", env)
		}
	}
}
