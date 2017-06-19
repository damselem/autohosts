package main

import (
	"log"
	"os"
	"sync"

	"github.com/damselem/autohosts/cmd"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
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
		var wg sync.WaitGroup
		errsCh := make(chan error)

		wg.Add(1)
		go func(dstFile string) {
			if err := cmd.RunAWSCommand(dstFile); err != nil {
				errsCh <- err
				return
			}

			wg.Done()
		}(*dstFile)

		wg.Add(1)
		go func(dstFile string, projects []string) {
			if err := cmd.RunGCPCommand(dstFile, projects); err != nil {
				errsCh <- err
				return
			}

			wg.Done()
		}(*dstFile, *allGcpProjects)

		go func() {
			for cloudErr := range errsCh {
				err = cloudErr
				wg.Done()
			}
		}()

		wg.Wait()
	}

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
