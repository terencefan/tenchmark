package xdispatcher

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type DataLoader interface {
	Load() error
	GetApi(name string) ([]byte, error)
	GetAllApis() ([][]byte, error)
}

type DirDataLoader struct {
	DataLoader
	dir  string
	apis map[string][]byte
}

func NewDirDataLoader(dir string) (DataLoader, error) {
	dataloader := &DirDataLoader{dir: dir, apis: make(map[string][]byte)}
	return dataloader, nil
}

func (d *DirDataLoader) Load() (err error) {
	filepath.Walk(d.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(info.Name(), ".in") {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			d.apis[info.Name()] = data
		}
		return nil
	})
	return
}

func (d *DirDataLoader) GetApi(name string) ([]byte, error) {
	if !strings.HasSuffix(name, ".in") {
		name = name + ".in"
	}
	if api, ok := d.apis[name]; ok {
		return api, nil
	} else {
		return nil, fmt.Errorf("loader have no api named %s", name)
	}
}

func (d *DirDataLoader) GetAllApis() ([][]byte, error) {
	res := make([][]byte, len(d.apis))
	idx := 0
	for _, api := range d.apis {
		res[idx] = api
		idx++
	}
	return res, nil
}

type FileDataLoader struct {
	DataLoader
	file string
	api  []byte
}

func (f *FileDataLoader) Load() (err error) {
	f.api, err = ioutil.ReadFile(f.file)
	return err
}

func (f *FileDataLoader) GetApi(name string) ([]byte, error) {
	if !strings.HasSuffix(name, ".in") {
		name = name + ".in"
	}

	if strings.HasSuffix(f.file, name) {
		return f.api, nil
	} else {
		return nil, fmt.Errorf("file have no api named %s", name)
	}
}
