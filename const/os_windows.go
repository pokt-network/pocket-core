// +build windows

package _const

import "os"

const (
	FILESEPARATOR = "\\"
)

var (
	DATADIR = os.Getenv("APPDATA") + "\\Pocket"
)
