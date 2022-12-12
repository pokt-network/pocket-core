package launcher

import (
	"fmt"
	"io"
	"os"
)

func writeBytesToFile(path string, contents []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	defer closeIgnoreError(file)
	if err != nil {
		return err
	}
	_, err = file.Write(contents)
	if err != nil {
		return err
	}
	return nil
}

func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer closeIgnoreError(source)

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer closeIgnoreError(destination)
	_, err = io.Copy(destination, source)
	return err
}

func closeIgnoreError(f *os.File) {
	_ = f.Close()
}
