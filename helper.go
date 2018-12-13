package main
import (
	"fmt"
	"os"
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	cp "github.com/otiai10/copy"
	"github.com/mholt/archiver"
)

//go get https://github.com/otiai10/copy
//go get -u github.com/mholt/archiver/cmd/arc

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

type configInfo struct {
	CP string `json:"cp"`
	CSB string `json:"csb"`
	XML string `json:"xml"`
	Output string `json:"output"`
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
	if _, err := os.Stat(d); os.IsNotExist(err) {
		os.MkdirAll(d, os.ModePerm)
	}
}

func Copy(src, dst string) {
	ndst := filepath.Join(dst, filepath.Base(src))
	createDir(ndst)
	cp.Copy(src, ndst)
}

func main() {
	conf, err := loadConf()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = RemoveContents(conf.Output)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = archiver.Archive([]string{conf.CP, conf.CSB, conf.XML}, filepath.Join(conf.Output, "src.zip"))
	if err != nil {
		fmt.Println(err)
		return
	}
}
