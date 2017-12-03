package xparser

import (
	"errors"
	"os"
	"time"
)

type FileOutputStream struct {
	outputFilename string
	fout           *os.File
}

func (t *FileOutputStream) Read(b []byte) (int, error) {
	return 0, errors.New("FileOutputStream don't support Read method")
}

func (t *FileOutputStream) Write(b []byte) (int, error) {
	return t.fout.Write(b)
}

func (t *FileOutputStream) Open() (err error) {
	if t.fout, err = os.Create(t.outputFilename); err != nil {
		return
	}
	return
}

func (t *FileOutputStream) Close() (err error) {
	if t.fout == nil {
		return
	}
	return t.fout.Close()
}

func (t *FileOutputStream) SetTimeout(d time.Duration) {
}

func (t *FileOutputStream) Flush() error {
	return nil
}

func NewFileOutputStream(out string) *FileOutputStream {
	return &FileOutputStream{
		outputFilename: out,
	}
}
