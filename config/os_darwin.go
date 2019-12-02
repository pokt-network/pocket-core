// +build darwin

package config

import "os"

const (
	FILESEPARATOR = "/"
)

var (
	DATADIR = os.Getenv("HOME") + "/Library/Pocket"
)
