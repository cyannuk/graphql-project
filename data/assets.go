package data

import (
	"bytes"
	"io/fs"
	"os"
	"path"
	"strings"

	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
)

type assetFS assetOperator

type assetFileFS struct {
	reader     *bytes.Reader
	name       string
	dirEntries []os.DirEntry
	offset     int
}

type assetDirEntry struct {
	info os.FileInfo
}

func (f *assetFileFS) Read(b []byte) (int, error) {
	return f.reader.Read(b)
}

func (f *assetFS) Open(name string) (fs.File, error) {
	if len(name) > 0 && (name[0] == '/' || name[0] == '\\') {
		name = name[1:]
	}
	if content, err := Asset(name); err == nil {
		return &assetFileFS{name: name, reader: bytes.NewReader(content)}, nil
	}
	if children, err := AssetDir(name); err == nil {
		dirEntries := make([]os.DirEntry, len(children))
		for i, child := range children {
			if info, err := AssetInfo(path.Join(name, child)); err == nil {
				dirEntries[i] = &assetDirEntry{info}
			} else {
				dirEntries[i] = &assetDirEntry{newDirFileInfo(child)}
			}
		}
		return &assetFileFS{name: name, dirEntries: dirEntries}, nil
	} else {
		if strings.Contains(err.Error(), "not found") {
			return nil, os.ErrNotExist
		}
		return nil, err
	}
}

func (f *assetFileFS) Name() string {
	return f.name
}

func (f *assetFileFS) Stat() (fs.FileInfo, error) {
	if len(f.dirEntries) != 0 {
		return newDirFileInfo(f.name), nil
	}
	return AssetInfo(f.name)
}

func (f *assetFileFS) Close() error {
	return nil
}

func (d *assetDirEntry) Name() string {
	name := d.info.Name()
	if d.info.IsDir() {
		return name
	} else {
		return path.Base(name)
	}
}

func (d *assetDirEntry) IsDir() bool {
	return d.info.IsDir()
}

func (d *assetDirEntry) Type() fs.FileMode {
	return d.info.Mode().Type()
}

func (d *assetDirEntry) Info() (fs.FileInfo, error) {
	return d.info, nil
}

func (f *assetFileFS) ReadDir(n int) ([]os.DirEntry, error) {
	if len(f.dirEntries) == 0 {
		return nil, os.ErrNotExist
	}
	if n <= 0 {
		return f.dirEntries, nil
	}
	if f.offset+n > len(f.dirEntries) {
		n = len(f.dirEntries) - f.offset
	}
	offset := f.offset
	f.offset += n
	return f.dirEntries[offset : offset+n], nil
}

func AssetFS() fs.FS {
	return &assetFS{}
}
