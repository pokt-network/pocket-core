// +build darwin

package _const

import "os"

const (
	FILESEPARATOR = "/"
)

var (
	DATADIR = os.Getenv("HOME") + "/Library/Pocket"
)
