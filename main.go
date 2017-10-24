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
	allAwsEmr := all.Flag("aws-with-emr", "include AWS EMR instances").Bool()
	allAwsAutoscale := all.Flag("aws-with-autoscale", "include AWS autoscaled instances").Bool()

	gcp := app.Command("gcp", "fetch instance hostnames from GCP")
	gcpProjects := gcp.Flag("project", "project id to fetch hostnames from").Strings()

	aws := app.Command("aws", "fetch instance hostnames from AWS")
	awsEmr := aws.Flag("with-emr", "include EMR instances").Bool()
	awsAutoscale := aws.Flag("with-autoscale", "include autoscaled instances").Bool()

	var err error
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case "aws":
		err = cmd.RunAWSCommand(*dstFile, *awsEmr, *awsAutoscale)
	case "gcp":
		err = cmd.RunGCPCommand(*dstFile, *gcpProjects)
	case "all":
		err = cmd.RunAllCommand(*dstFile, *allGcpProjects, *allAwsEmr, *allAwsAutoscale)
	}

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
