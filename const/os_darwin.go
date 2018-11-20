// +build darwin

package _const

/*
This file is for OS specific data directory configuration
 */
const (
	FILESEPARATOR = "/"
)

var (
	DATADIR = os.Getenv("HOME") + "/Library/Pocket"
)
