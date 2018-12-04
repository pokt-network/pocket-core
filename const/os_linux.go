// +build linux

package _const

import "os"

// "os_linux.go" is for OS specific constants.

const (
	FILESEPARATOR = "/" // os specific file separator
)

var (
	DATADIR = os.Getenv("HOME") + "/.pocket" // os specific data directory.
)
