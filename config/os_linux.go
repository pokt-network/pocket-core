// +build linux

package config

import "os"

const (
	FILESEPARATOR = "/"
)

var (
	DATADIR = os.Getenv("HOME") + "/.pocket"
)
