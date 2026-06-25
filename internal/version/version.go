package version

import "fmt"

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func Short() string {
	return fmt.Sprintf("tux-letter %s", Version)
}

func String() string {
	return fmt.Sprintf("tux-letter %s\ncommit: %s\ndate: %s", Version, Commit, Date)
}
