// +build windows

package _const

import "os"

/*
This file is for OS specific data directory configuration
 */
 var(
 	DATADIR = os.Getenv("APPDATA")+"\\"+"Pocket"
 )
