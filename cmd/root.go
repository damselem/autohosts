package cmd

import (
	"fmt"
	"log"

	"github.com/damselem/autohosts/hosts"
	"github.com/spf13/cobra"
)

func GetRootCmd() *cobra.Command {
	// Flag vars
	var dstFile string

	var rootCmd = &cobra.Command{
		Use:   "autohosts",
		Short: "Manage hosts for AWS and GCP instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			awsEntries, err := hosts.AWS()
			if err != nil {
				log.Println(err)
			}

			if dstFile != "" {
				hostFile := hosts.NewHostFile("/etc/hosts")
				return hostFile.Update(awsEntries)
			}

			for _, awsEntry := range awsEntries {
				fmt.Println(awsEntry.String())
			}

			return nil
		},
	}

	rootCmd.Flags().StringVarP(
		&dstFile,
		"output",
		"o",
		"",
		"file to update with new entries",
	)

	return rootCmd
}
