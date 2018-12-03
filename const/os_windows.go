// +build windows

package _const

import "os"

// "os_windows.go" is for OS specific constants.

const (
	FILESEPARATOR = "\\"								// os specific file separator
)

var (
	DATADIR = os.Getenv("APPDATA") + "\\Pocket"		// os specific data directory.
)
