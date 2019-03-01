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
	GetAllApis() [][]byte
	GetName(index int) (string, error)
}

type DirDataLoader struct {
	DataLoader
	dir   string
	apis  [][]byte
	names []string
}

func NewDirDataLoader(dir string) (DataLoader, error) {
	dataloader := &DirDataLoader{dir: dir, apis: make([][]byte, 0), names: make([]string, 0)}
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
			d.apis = append(d.apis, data)
			d.names = append(d.names, info.Name())
		}
		return nil
	})
	return
}

func (d *DirDataLoader) GetApi(name string) ([]byte, error) {
	if !strings.HasSuffix(name, ".in") {
		name = name + ".in"
	}
	for i := 0; i < len(d.names); i++ {
		if name == d.names[i] {
			return d.apis[i], nil
		}
	}
	return nil, fmt.Errorf("loader have no api named %s", name)
}

func (d *DirDataLoader) GetAllApis() [][]byte {
	return d.apis
}

func (d *DirDataLoader) GetName(index int) (string, error) {
	if index >= len(d.names) {
		return "", fmt.Errorf("name index out of range, %d", index)
	}
	return d.names[index], nil
}

type FileDataLoader struct {
	DataLoader
	file string
	api  []byte
}

func NewFileDataLoader(file string) (*FileDataLoader, error) {
	loader := &FileDataLoader{file: file}
	return loader, nil
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

func (f *FileDataLoader) GetAllApis() [][]byte {
	res := make([][]byte, 1)
	res[0] = f.api
	return res
}

func (f *FileDataLoader) GetName(index int) (string, error) {
	if index != 0 {
		return "", fmt.Errorf("name index out of range, %d", index)
	}
	return f.file, nil
}
