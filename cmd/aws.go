package cmd

import (
	"fmt"

	"github.com/damselem/autohosts/hosts"
)

// RunAWSCommand prints hosts from AWS
func RunAWSCommand(dstFile string, withEmr, withAutoscale bool) error {
	if dstFile != "" {
		fmt.Printf("Updating %s with new hostnames from AWS...\n", dstFile)
	}

	entries, err := hosts.AWS(withEmr, withAutoscale)
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
