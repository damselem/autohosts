package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/damselem/autohosts/hosts"
)

// RootCmd is the entry point for registering subcommands
var RootCmd = &cobra.Command{
	Use:   "autohosts",
	Short: "Manage hosts for AWS and GCP instances",
	Run: func(cmd *cobra.Command, args []string) {
		awsEntries, err := hosts.AWS()
		if err != nil {
			log.Println(err)
		}

		for _, entry := range awsEntries {
			fmt.Println(entry.String())
		}
	},
}
