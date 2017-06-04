package hosts

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
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

func (hf *file) Update(managedEntries []Entry) error {
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
	writeManagedEntries(f, managedEntries)

	err = moveAndApplyPermissions(f.Name(), hf.path)
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
		_, err := f.WriteString(fmt.Sprintf("%s\n", entry.String()))
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
			} else if line != "" {
				entry := parseLineEntry(line)
				managedEntries = append(managedEntries, entry)
			}
		} else {
			if line == startComment {
				blockStarted = true
			} else if line != "" {
				entry := parseLineEntry(line)
				localEntries = append(localEntries, entry)
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
	var address, hostname, comment string
	var aliases []string

	line = strings.TrimSpace(line)

	lineParts := strings.SplitN(line, "#", 2)
	entryFields := strings.Fields(lineParts[0])

	if len(lineParts) == 2 {
		comment = fmt.Sprintf("#%s", lineParts[1])
	}

	if len(entryFields) > 0 {
		address = entryFields[0]
	}
	if len(entryFields) > 1 {
		hostname = entryFields[1]
	}
	if len(entryFields) > 2 {
		aliases = entryFields[2:len(entryFields)]
	}

	return Entry{
		Address:  address,
		Hostname: hostname,
		Aliases:  aliases,
		Comment:  comment,
	}
}

func moveAndApplyPermissions(tmpHostFilePath, dstHostFilePath string) error {
	fi, err := os.Stat(dstHostFilePath)
	if err != nil {
		return errors.New(fmt.Sprintf("%s could not be updated", dstHostFilePath))
	}

	uid := fi.Sys().(*syscall.Stat_t).Uid
	gid := fi.Sys().(*syscall.Stat_t).Gid

	runInteractively("sudo", "-s", "mv", tmpHostFilePath, dstHostFilePath)
	return runInteractively("sudo", "-s", "chown", fmt.Sprintf("%d:%d", uid, gid), dstHostFilePath)
}

func runInteractively(command string, arguments ...string) error {
	cmd := exec.Command(command, arguments...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
