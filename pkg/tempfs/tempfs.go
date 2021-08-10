package tempfs

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type TempFS interface {
	Open(name string) (fs.File, error)
	ReadDir(name string) ([]fs.DirEntry, error)
}

type NopTempFS struct{}

func (n *NopTempFS) Open(name string) (fs.File, error) {
	return os.Create("NopTempFS")
}

func (n *NopTempFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return []fs.DirEntry{}, nil
}

type UserTempFS struct {
	path string
	fs   fs.FS
}

func NewUserTempFS(name string) (*UserTempFS, error) {
	if !fs.ValidPath(name) {
		return nil, errors.New("invalid path")
	}
	if _, err := os.Open(name); err != nil {
		return nil, errors.New("dir path not exists")
	}
	name, err := filepath.Abs(name)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(name, string(filepath.Separator)) {
		name = strings.TrimSuffix(name, string(filepath.Separator)) 
	}
	name = strings.TrimSuffix(name, "template")
	ufs := &UserTempFS{
		path: name,
		fs:   os.DirFS(name),
	}
	return ufs, nil
}

func (ufs *UserTempFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, errors.New("invalid file path")
	}
	return ufs.fs.Open(name)
}

func (ufs *UserTempFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(ufs.fs, name)
}
