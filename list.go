package main

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/retzkek/grafanactl/gapi"
)

var listCmd = &Command{
	Name:    "list",
	Usage:   "[OPTIONS]",
	Summary: "List dashboards.",
	Help:    `The list command lists dashboard names and meta information.`,
}

var (
	format = listCmd.Flag.String("format", "short",
		"list format: short, long, json")
)

func listFunc(client *gapi.Client, cmd *Command, args []string) error {
	dl, err := client.ListDashboards()
	if err != nil {
		log.Error(err)
		return fmt.Errorf("error getting dashboard list")
	}
	switch *format {
	case "short":
		for _, db := range *dl {
			fmt.Println(db.URI)
		}
	case "long":
		fmt.Printf("ID     URI                                      TITLE\n")
		for _, db := range *dl {
			fmt.Printf("%-6d %-40s %-40s\n", db.Id, db.URI, db.Title)
		}
	case "json":
		b, err := json.MarshalIndent(dl, "", "\t")
		if err != nil {
			log.Error(err)
			return fmt.Errorf("error marshalling dashboard list to JSON")
		}
		fmt.Printf("%s", b)
	default:
		return fmt.Errorf("unknown format %s", *format)
	}
	return nil
}

func init() {
	listCmd.Function = listFunc
}
