package config

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pokt-network/pocket-core/const"
)

func WriteFixtures() error {
	datadirectory := GlobalConfig().DD + _const.FILESEPARATOR
	cFile, err := filepath.Abs("config" + _const.FILESEPARATOR + "fixtures" + _const.FILESEPARATOR + "chains.json")
	if err != nil {
		return err
	}
	pFile, err := filepath.Abs("config" + _const.FILESEPARATOR + "fixtures" + _const.FILESEPARATOR + "peers.json")
	if err != nil {
		return err
	}
	err = copyFile(cFile, datadirectory+"chains.json")
	if err != nil {
		return err
	}
	err = copyFile(pFile, datadirectory+"peers.json")
	if err != nil {
		return err
	}
	return nil
}

func copyFile(src, dst string) error {
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		in, err := os.Open(src)
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, in)
		if err != nil {
			return err
		}
		return out.Close()
	}
	return nil
}

func filePaths() {
	if *cFile == _const.CHAINSFILENAME {
		*cFile = *dd + _const.FILESEPARATOR + "chains.json"
	}
}
