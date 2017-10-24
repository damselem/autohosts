package cmd

import (
	"fmt"
	"sync"

	"github.com/damselem/autohosts/hosts"
)

// RunAllCommand prints all hosts form AWS and GCP
func RunAllCommand(dstFile string, projects []string, withAwsEmr, withAwsAutoscaled bool) error {
	var wg sync.WaitGroup

	entriesCh := make(chan []hosts.Entry)
	errsCh := make(chan error)

	if dstFile != "" {
		fmt.Printf("Updating %s with new hostnames from AWS and GCP...\n", dstFile)
	}

	wg.Add(1)
	go func() {
		entries, err := hosts.AWS(withAwsEmr, withAwsAutoscaled)
		if err != nil {
			errsCh <- err
			return
		}

		entriesCh <- entries
	}()

	wg.Add(1)
	go func(projects []string) {
		entries, err := hosts.GCP(projects)
		if err != nil {
			errsCh <- err
			return
		}

		entriesCh <- entries
	}(projects)

	var entries []hosts.Entry
	go func() {
		for provEntries := range entriesCh {
			entries = append(entries, provEntries...)
			wg.Done()
		}
	}()

	var instErrs []error
	go func() {
		for provErr := range errsCh {
			instErrs = append(instErrs, provErr)
			wg.Done()
		}
	}()

	wg.Wait()

	if len(instErrs) > 0 {
		return instErrs[0]
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
