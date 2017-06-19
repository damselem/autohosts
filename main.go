package main

import (
	"log"
	"os"

	"github.com/damselem/autohosts/cmd"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("autohosts", "")
	dstFile := app.Flag("output", "file to update with hosts").Short('o').String()

	all := app.Command("all", "fetch instance hostnames from AWS and GCP")
	allGcpProjects := all.Flag("gcp-project", "GCP project ID to fetch hostnames from").Strings()

	gcp := app.Command("gcp", "fetch instance hostnames from GCP")
	gcpProjects := gcp.Flag("project", "project id to fetch hostnames from").Strings()

	app.Command("aws", "fetch instance hostnames from AWS")

	var err error
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case "aws":
		err = cmd.RunAWSCommand(*dstFile)
	case "gcp":
		err = cmd.RunGCPCommand(*dstFile, *gcpProjects)
	case "all":
		err = cmd.RunAllCommand(*dstFile, *allGcpProjects)
	}

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
