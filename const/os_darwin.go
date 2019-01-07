// +build darwin

package _const

import "os"

// "os_darwin.go" is for OS specific constants.

const (
	FILESEPARATOR = "/" // os specific file separator
)

var (
	DATADIR = os.Getenv("HOME") + "/Library/Pocket" // os specific data directory.
)
