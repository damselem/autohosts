package cmd

import (
	"fmt"

	"github.com/damselem/autohosts/hosts"
)

func RunAWSCommand(dstFile string) error {
	if dstFile != "" {
		fmt.Printf("Updating %s with new hostnames from AWS...\n", dstFile)
	}

	entries, err := hosts.AWS()
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
