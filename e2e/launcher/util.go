package launcher

import (
	"fmt"
	"io"
	"os"
)

const (
	filePermissions = 0666 // chmod a+rwx,u-x,g-x,o-x,ug-s,-t:  (U)ser / owner can read, can write and can't execute. (G)roup can read, can write and can't execute. (O)thers can read, can write and can't execute.

)

func writeBytesToFile(path string, contents []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, filePermissions)
	if err != nil {
		return err
	}
	defer closeIgnoreError(file)

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
