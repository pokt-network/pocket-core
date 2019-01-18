// +build linux

package _const

import "os"

const (
	FILESEPARATOR = "/"
)

var (
	DATADIR = os.Getenv("HOME") + "/.pocket"
)
