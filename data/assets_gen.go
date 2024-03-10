// Code generated for package data by go-bindata DO NOT EDIT. (@generated)
// sources:
// db/migrations/20240310083405_init.sql
package data

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data, name string) ([]byte, error) {
	gz, err := gzip.NewReader(strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}


type assetFile struct {
	*bytes.Reader
	name            string
	childInfos      []os.FileInfo
	childInfoOffset int
}

type assetOperator struct{}

// Open implement http.FileSystem interface
func (f *assetOperator) Open(name string) (http.File, error) {
	var err error
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}
	content, err := Asset(name)
	if err == nil {
		return &assetFile{name: name, Reader: bytes.NewReader(content)}, nil
	}
	children, err := AssetDir(name)
	if err == nil {
		childInfos := make([]os.FileInfo, 0, len(children))
		for _, child := range children {
			childPath := filepath.Join(name, child)
			info, errInfo := AssetInfo(filepath.Join(name, child))
			if errInfo == nil {
				childInfos = append(childInfos, info)
			} else {
				childInfos = append(childInfos, newDirFileInfo(childPath))
			}
		}
		return &assetFile{name: name, childInfos: childInfos}, nil
	} else {
		// If the error is not found, return an error that will
		// result in a 404 error. Otherwise the server returns
		// a 500 error for files not found.
		if strings.Contains(err.Error(), "not found") {
			return nil, os.ErrNotExist
		}
		return nil, err
	}
}

// Close no need do anything
func (f *assetFile) Close() error {
	return nil
}

// Readdir read dir's children file info
func (f *assetFile) Readdir(count int) ([]os.FileInfo, error) {
	if len(f.childInfos) == 0 {
		return nil, os.ErrNotExist
	}
	if count <= 0 {
		return f.childInfos, nil
	}
	if f.childInfoOffset+count > len(f.childInfos) {
		count = len(f.childInfos) - f.childInfoOffset
	}
	offset := f.childInfoOffset
	f.childInfoOffset += count
	return f.childInfos[offset : offset+count], nil
}

// Stat read file info from asset item
func (f *assetFile) Stat() (os.FileInfo, error) {
	if len(f.childInfos) != 0 {
		return newDirFileInfo(f.name), nil
	}
	return AssetInfo(f.name)
}

// newDirFileInfo return default dir file info
func newDirFileInfo(name string) os.FileInfo {
	return &bindataFileInfo{
		name:    name,
		size:    0,
		mode:    os.FileMode(2147484068), // equal os.FileMode(0644)|os.ModeDir
		modTime: time.Time{}}
}

// AssetFile return a http.FileSystem instance that data backend by asset
func AssetFile() http.FileSystem {
	return &assetOperator{}
}

var _dbMigrations20240310083405_initSql = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xc4\x55\xc1\x8e\xda\x30\x10\xbd\xe7\x2b\x46\x7b\x02\xa9\x2b\x55\xea\xad\x3d\x85\x64\x58\x59\x05\x43\x1d\x47\x82\x13\x32\xb1\x45\x2d\x41\x9c\x3a\x93\x65\xe9\xd7\x57\x66\x17\xba\x90\x94\xb4\xbd\x70\xf5\xbc\x37\x1e\xbf\xf7\x46\x7e\x7c\x84\x9d\xdd\x78\x45\xe6\x73\x53\x45\x89\xc0\x58\x22\xe0\x42\x22\xcf\xd8\x8c\x03\x1b\x03\x9f\x49\xc0\x05\xcb\x64\x06\x0f\x85\x25\xf3\x42\x0f\x5f\xa2\x13\x52\xc6\xa3\x09\x5e\xa1\x9a\xda\xf8\x1a\x06\x11\x80\xd5\x30\x62\x4f\x8c\x4b\x78\x42\x8e\x22\x96\x98\xc2\x68\x09\x29\x8e\xe3\x7c\x22\x21\xce\x80\xa5\xc8\x25\x93\x4b\x98\x0b\x36\x8d\xc5\x12\xbe\xe2\xf2\x43\x04\x50\x78\xa3\xc8\xe8\x95\x22\x90\x6c\x8a\x99\x8c\xa7\xf3\xe3\x1d\x3c\x9f\x4c\xce\x0d\x92\x5c\x08\xe4\x72\x75\x86\x0c\x3e\x0d\x03\xbb\x54\x3b\x03\x12\x17\xf2\x4c\x09\xa7\x66\xa7\xec\x16\x12\x76\x51\x80\x9c\xb3\x6f\x39\x86\xba\xd2\xda\x9b\xba\x6e\x13\x0b\x4b\x87\xf6\x69\x4d\x8a\xcc\x75\xbb\x50\xf8\x69\xab\x36\x7a\x6d\x3d\x7d\x5f\xe9\x40\x49\x83\x70\xef\x6b\x5b\x45\x96\x1a\x6d\x20\x9d\xe5\x41\xce\xb9\xc0\x84\x1d\xe5\xbf\x40\xb9\x72\xf3\x17\xb0\x4a\xd5\xf5\xde\x79\xdd\x31\xb0\x6b\x7c\xd1\x9a\x38\x1a\xde\x76\xb3\xf2\x4e\x37\x05\xdd\xd7\xd0\x42\x91\xd9\x38\x7f\xe8\x92\xdb\xa8\xb2\xfd\xd6\xca\xdb\xa2\x47\xa8\x1f\x8d\x2a\x29\x38\x1b\x9e\xd3\x9a\xe4\x63\x80\x78\x45\xb6\xdc\xdc\x6e\xd3\x9d\xb5\x67\x53\x6a\xe7\xe1\x9f\x94\x76\x5e\xdf\x7b\x71\xc2\xea\xae\x7e\x5f\x7e\xe6\x09\x1c\xa3\x40\x9e\xe0\xdb\x76\x0f\xac\x1e\xbe\xca\x7c\x0c\x47\x0f\xe5\x14\xa1\x13\x4b\xdb\xba\x70\x4d\x49\x7f\x16\xf6\xd2\x86\x4e\xa7\x8e\x89\x6e\xd6\xe4\x48\x6d\x6f\x3b\x44\xea\xa5\x07\x70\xbb\x47\x9f\x6f\xde\x3c\x5b\xb3\xbf\xaf\x71\xaf\x33\x18\xdf\xb5\x09\xff\x63\xd1\x5b\xf2\xaf\x05\x5f\x3b\x7d\xe8\x08\xf5\xbb\x4f\x44\xbb\x7d\x19\xfd\x0a\x00\x00\xff\xff\x53\xd1\x10\x50\x55\x06\x00\x00"

func dbMigrations20240310083405_initSqlBytes() ([]byte, error) {
	return bindataRead(
		_dbMigrations20240310083405_initSql,
		"db/migrations/20240310083405_init.sql",
	)
}

func dbMigrations20240310083405_initSql() (*asset, error) {
	bytes, err := dbMigrations20240310083405_initSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "db/migrations/20240310083405_init.sql", size: 1621, mode: os.FileMode(438), modTime: time.Unix(1710104348, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"db/migrations/20240310083405_init.sql": dbMigrations20240310083405_initSql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"db": &bintree{nil, map[string]*bintree{
		"migrations": &bintree{nil, map[string]*bintree{
			"20240310083405_init.sql": &bintree{dbMigrations20240310083405_initSql, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
