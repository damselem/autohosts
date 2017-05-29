package hosts

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	startComment = "# AUTOHOSTS - START"
	endComment   = "# AUTOHOSTS - END"
)

type file struct {
	path string
}

func NewHostFile(path string) *file {
	return &file{
		path: path,
	}
}

func (hf *file) Update(remoteEntries []Entry) error {
	entries, err := hf.parseHostFile()
	if err != nil {
		return err
	}
	f, err := createTmpHostFile()
	if err != nil {
		return err
	}
	defer f.Close()
	err = writeLocalEntries(f, entries["local"])
	if err != nil {
		return err
	}

	f.WriteString("\n")
	writeManagedEntries(f, remoteEntries)

	fmt.Println(remoteEntries)
	// moveAndApplyPermissions(tmpHostFilePath, "/etc/hosts")
	return err
}

func writeLocalEntries(f *os.File, entries []Entry) error {
	return writeEntries(f, entries)
}

func writeManagedEntries(f *os.File, entries []Entry) {
	f.WriteString(fmt.Sprintf("%s\n", startComment))
	writeEntries(f, entries)
	f.WriteString(fmt.Sprintf("%s\n", endComment))
}

func writeEntries(f *os.File, entries []Entry) error {
	for _, entry := range entries {
		_, err := f.WriteString(entry.String())
		if err != nil {
			return err
		}
	}
	return nil
}

func createTmpHostFile() (*os.File, error) {
	f, err := ioutil.TempFile("", "hosts")
	if err != nil {
		return nil, err
	}
	err = f.Chmod(os.FileMode(0644))
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (hf *file) parseHostFile() (map[string][]Entry, error) {
	file, err := os.Open(hf.path)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	var managedEntries, localEntries []Entry
	blockStarted := false
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if blockStarted {
			if line == endComment {
				blockStarted = false
			} else {
				if line != "" {
					entry := parseLineEntry(line)
					managedEntries = append(managedEntries, entry)
				}
			}
		} else {
			if line == startComment {
				blockStarted = true
			} else {
				if line != "" {
					entry := parseLineEntry(line)
					localEntries = append(localEntries, entry)
				}
			}
		}
	}

	entries := map[string][]Entry{
		"managed": managedEntries,
		"local":   localEntries,
	}

	return entries, nil
}

func parseLineEntry(line string) Entry {
	var aliases []string

	lineParts := strings.Split(line, "#")
	fields := strings.Fields(lineParts[0])
	comment := lineParts[1]

	if len(fields) > 2 {
		aliases = fields[2 : len(fields)-1]
	}
	return Entry{
		Address:  fields[0],
		Hostname: fields[1],
		Aliases:  aliases,
		Comment:  comment,
	}
}
