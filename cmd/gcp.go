package cmd

import (
	"fmt"

	"github.com/damselem/autohosts/hosts"
)

func RunGCPCommand(dstFile string, projects []string) error {
	if dstFile != "" {
		fmt.Printf("Updating %s with new hostnames from GCP...\n", dstFile)
	}

	entries, err := hosts.GCP(projects)
	if err != nil {
		return err
	}

	if dstFile != "" {
		hostFile := hosts.NewHostFile(dstFile)
		return hostFile.Update(entries)
	}

	for _, entry := range entries {
		fmt.Println(entry.String())
	}

	return nil
}
