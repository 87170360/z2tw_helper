package main
import (
	"fmt"
	"os"
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"github.com/mholt/archiver"
	cp "github.com/otiai10/copy"
)

//go get -u github.com/mholt/archiver/cmd/archiver

func RemoveContents(dir string) error {
    d, err := os.Open(dir)
    if err != nil {
        return err
    }
    defer d.Close()
    names, err := d.Readdirnames(-1)
    if err != nil {
        return err
    }
    for _, name := range names {
        err = os.RemoveAll(filepath.Join(dir, name))
        if err != nil {
            return err
        }
    }
    return nil
}

func secCopy(src, dst string) {
	nsrc := filepath.Join(src, "output", filepath.Base(dst))
	//fmt.Println("src", nsrc)
	//fmt.Println("dst", dst)
	err := cp.Copy(nsrc, dst)
	if err != nil {
		fmt.Println(err)
	}
}

type configInfo struct {
	CP string `json:"cp"`
	CSB string `json:"csb"`
	XML string `json:"xml"`
	Output string `json:"output"`
	Input string `json:"input"`
	File1 string `json:file1`
	File2 string `json:file2`
}

func loadConf() (*configInfo, error) {
	data, err := ioutil.ReadFile("conf.json")
    if err != nil {
        return nil, err
    }

	ci := &configInfo{}
    err = json.Unmarshal(data, &ci)
    if err != nil {
        return nil, err
    }
    return ci, nil
}

func createDir(d string) {
	if isFileAndDirExist(d){
		os.MkdirAll(d, os.ModePerm)
	}
}

func isFileAndDirExist(d string) bool {
	if _, err := os.Stat(d); os.IsNotExist(err) {
		return false
	}
	return true
}

func main() {
	conf, err := loadConf()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = RemoveContents(conf.Input)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = archiver.Archive([]string{conf.CP, conf.CSB, conf.XML}, filepath.Join(conf.Input, conf.File1))
	if err != nil {
		fmt.Println(err)
		return
	}

	outputFile := filepath.Join(conf.Output, conf.File2)
	if isFileAndDirExist(outputFile) {
		err = archiver.Unarchive(outputFile, conf.Output)
		if err != nil {
			fmt.Println(err)
			return
		}

		secCopy(conf.Output, conf.CP)
		secCopy(conf.Output, conf.CSB)
		secCopy(conf.Output, conf.XML)
	}

	err = RemoveContents(conf.Output)
	if err != nil {
		fmt.Println(err)
		return
	}
}
