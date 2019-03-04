package config

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/pokt-network/pocket-core/crypto"
)

func gidHash() string {
	hashString, err := crypto.NewSHA1Hash()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return *gid + ":" + hashString
}

func checkPrefix(gid string, query string) error {
	if index := strings.IndexByte(query, ':'); index > 0 {
		prefix := query[:index]
		if prefix != gid {
			return errors.New("incorrect gid prefix")
		}
		return nil
	}
	return errors.New("prefix not detected ':' not found")
}

func createGIDFile(filepath string) string{
	f, err := os.Create(filepath)
	if err != nil {
		log.Fatalf(err.Error())
	}
	s := gidHash()
	f.Write([]byte(s))
	return s
}

func GIDSetup() string {
	fp := filepath.FromSlash(*dd + "/gid.dat")
	if _, err := os.Stat(fp); err == nil { // file exists
		b, err := ioutil.ReadFile(fp)
		err2 := checkPrefix(*gid, string(b))
		// unable to read file or prefix is not the same as gid
		if err != nil || err2 != nil {
			err := os.Remove(fp)
			if err != nil {  // unable to delete file
				log.Fatalf(err.Error())
			}
			return createGIDFile(fp)
		}
		return string(b)
	} else if os.IsNotExist(err) {
		return createGIDFile(fp)
	} else {
		log.Fatalf(err.Error())
	}
	return ""
}
