package hosts

import (
	"fmt"
	"strings"
)

// Entry represents an entry in /etc/hosts
type Entry struct {
	Address  string
	Hostname string
	Aliases  []string
	Comment  string
}

func (e *Entry) String() string {
	aliases := strings.Join(e.Aliases, " ")
	return fmt.Sprintf("%-15s %s %s %s", e.Address, e.Hostname, aliases, e.Comment)
}
