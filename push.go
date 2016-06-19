package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	gapi "github.com/retzkek/go-grafana-api"
)

var pushCmd = &Command{
	Name:    "push",
	Usage:   "[OPTIONS] [DASHBOARD...]",
	Summary: "Read dashboards from file and push to Grafana.",
	Help: `The push command reads dashboards from file and pushes them to Grafana.
If no dashboards are specified, push all dashboards in the specified path.
Specify dashboards by slug, e.g. 'db/foo' or just 'foo'.

Since only database-stored dashboards can be saved through the Grafana API,
only dashboards in the 'db' sub-directory are pushed.`,
}

var (
	pushPath = pushCmd.Flag.String("path", ".",
		"path to read files from (default is current working dir)")
	pushOverwrite = pushCmd.Flag.Bool("overwrite", false,
		"overwrite existing dashboards")
)

func pushFunc(client *gapi.Client, cmd *Command, args []string) error {
	dirname := filepath.Join(*pushPath, "db")
	df, err := os.Open(dirname)
	defer df.Close()
	if err != nil {
		log.WithField("path", dirname).Error(err)
		return fmt.Errorf("error opening dashboard directory")
	}
	// get dashboards from args or directory list
	var dashboards []string
	if len(args) == 0 {
		dashboards, err = df.Readdirnames(0)
		if err != nil {
			log.WithField("path", dirname).Error(err)
			return fmt.Errorf("error getting list of dashboards")
		}
	} else {
		dashboards = make([]string, len(args))
		for i, s := range args {
			dashboards[i] = s + ".json"
		}
	}
	// read and push dashboards
	for _, d := range dashboards {
		filename := filepath.Join(df.Name(), d)
		log.WithField("filename", filename).Info("saving dashboard")
		model, err := readDashboard(filename)
		if err != nil {
			log.WithField("file", filename).Error(err)
			return fmt.Errorf("error loading dashboard from file")
		}
		resp, err := client.SaveDashboard(model, *pushOverwrite)
		// grafana returns 404 if dashboard we're trying to send includes an id,
		// but no dashboard exists with that db. Try sending dashboard with nil
		// id (i.e. create new).
		if err != nil {
			if err.Error() == "404 Not Found" {
				log.Warning("Grafana returned 404. Trying to create new dashboard.")
				model["id"] = nil
				resp, err = client.SaveDashboard(model, *pushOverwrite)
			}
		}
		if err != nil {
			log.WithField("file", filename).Error(err)
			return fmt.Errorf("error pushing dashboard to Grafana")
		}
		if resp != nil {
			log.WithFields(log.Fields{
				"slug":    resp.Slug,
				"status":  resp.Status,
				"version": resp.Version,
			}).Info("dashboard saved")
		}
	}
	return nil
}

// readDashboard reads a JSON dashboard from file and unmarshals it
func readDashboard(filename string) (map[string]interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dat, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var v map[string]interface{}
	if err = json.Unmarshal(dat, &v); err != nil {
		return nil, err
	}
	return v, nil
}

func init() {
	pushCmd.Function = pushFunc
}
