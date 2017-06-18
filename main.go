package main

import (
	"github.com/damselem/autohosts/cmd"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
)

func main() {
	app := kingpin.New("autohosts", "")
	dstFile := app.Flag("output", "file to update with hosts").Short('o').String()

	all := app.Command("all", "fetch instance hostnames from AWS and GCP")
	allGcpProjects := all.Flag("projects", "list of GCP project ids to fetch hostnames from").Strings()

	gcp := app.Command("gcp", "fetch instance hostnames from GCP")
	gcpProjects := gcp.Flag("projects", "list of project ids to fetch hostnames from").Strings()

	app.Command("aws", "fetch instance hostnames from AWS")

	var err error
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case "aws":
		err = cmd.RunAWSCommand(*dstFile)
	case "gcp":
		err = cmd.RunGCPCommand(*dstFile, *gcpProjects)
	case "all":
		err = cmd.RunAWSCommand(*dstFile)
		err = cmd.RunGCPCommand(*dstFile, *allGcpProjects)
	}

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
