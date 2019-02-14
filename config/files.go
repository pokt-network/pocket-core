package config

import (
	"io"
	"os"
	"path/filepath"
)

func WriteFixtures() error {
	datadirectory := GlobalConfig().DD
	cFile, err := filepath.Abs("./config/fixtures/chains.json")
	if err != nil {
		return err
	}
	pFile, err := filepath.Abs("./config/fixtures/peers.json")
	if err != nil {
		return err
	}
	dwl, err := filepath.Abs("./config/fixtures/developer_whitelist.json")
	if err != nil {
		return err
	}
	swl, err := filepath.Abs("./config/fixtures/service_whitelist.json")
	if err != nil {
		return err
	}
	err = copyFile(cFile, datadirectory+"/chains.json")
	if err != nil {
		return err
	}
	err = copyFile(pFile, datadirectory+"/peers.json")
	if err != nil {
		return err
	}
	copyFile(dwl, datadirectory+"/developer_whitelist.json")
	if err != nil {
		return err
	}
	copyFile(swl, datadirectory+"/service_whitelist.json")
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
