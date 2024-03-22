// Code generated for package data by go-bindata DO NOT EDIT. (@generated)
// sources:
// db/migrations/20240310083405_init.sql
// db/migrations/20240312074226_add_role.sql
// db/migrations/20240315173713_add_tokens.sql
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

var _dbMigrations20240310083405_initSql = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x95\xc1\x8e\xda\x30\x10\x86\xef\x3c\xc5\x28\x27\x90\xba\x52\xa5\xde\xda\x53\x20\xc3\xca\x2a\x18\x9a\x38\x12\x9c\x2a\x13\x5b\xd4\x12\xc4\xa9\x33\x59\x76\xfb\xf4\x55\xbc\x4b\x5a\xc0\x6c\x56\x5c\x76\xaf\x99\xff\x1f\x8f\xe7\x1b\x4f\xee\xee\x60\x6f\xb6\x4e\x92\xfe\xda\x54\x83\x49\x8a\xb1\x40\xc0\x95\x40\x9e\xb1\x05\x07\x36\x05\xbe\x10\x80\x2b\x96\x89\x0c\xa2\xc2\x90\x7e\xa4\xe8\xdb\xe0\xa8\x14\xf1\x78\x86\x67\xaa\xa6\xd6\xae\x86\xe1\x00\x20\x32\x2a\x82\x31\xbb\x67\x5c\xc0\x3d\x72\x4c\x63\x81\x09\x8c\xd7\x90\xe0\x34\xce\x67\x02\xe2\x0c\x58\x82\x5c\x30\xb1\x86\x65\xca\xe6\x71\xba\x86\xef\xb8\xfe\xd4\x7a\x0b\xa7\x25\x69\x15\x53\x04\x82\xcd\x31\x13\xf1\x7c\xe9\x8f\xe1\xf9\x6c\xd6\x65\x98\xe4\x69\x8a\x5c\xfc\xec\x24\xc3\x2f\x23\x6f\x2f\xe5\x5e\x47\x20\x70\x25\x3a\x93\xff\xae\xf7\xd2\xec\x22\x98\xb0\x93\x10\xe4\x9c\xfd\xc8\xd1\x2b\xa4\x52\x4e\xd7\x75\xc8\x5c\x18\x7a\x0a\x7d\xaf\x49\x92\xbe\x48\xea\x43\x7f\x4c\x15\x72\x6c\x8c\xa3\x5f\x89\x77\x25\x6d\x23\x4f\x82\x3b\x49\x86\x1a\xd5\xc6\x16\x79\xdb\xe0\x65\x8a\x13\xe6\x81\x9c\xea\x6c\xb9\x7d\x93\xb0\x92\x75\x7d\xb0\x4e\x05\x6b\xb7\x8d\x2b\xae\x14\xaf\xf4\x4e\x9f\x33\x18\x8c\x5e\xc7\x5f\x39\xab\x9a\x82\xde\x7f\x02\x0a\x49\x7a\x6b\xdd\x53\xf8\x6a\x5a\x96\xa1\x6e\x54\xce\x14\xbd\xed\xfc\xdd\xc8\x92\xfc\x28\xb4\xf7\xba\xa8\xe8\xb3\x17\x39\x49\xa6\xdc\xf6\xa5\xba\x36\xa6\x0f\xba\x54\xd6\x85\x22\xb7\x30\xb1\x4e\x7d\x84\x37\xd9\x6e\x06\xf6\xef\xf8\xce\x98\xe2\x14\x53\xe4\x13\x7c\x59\x1e\x43\xa3\x46\x2f\x34\xfc\x2c\xf5\x78\x8e\x13\xd7\xd9\x94\xa9\x0b\xdb\x94\xf4\x4a\xf3\xcf\x60\x85\x89\x3e\xbf\x8f\x66\x43\x96\xe4\xae\x8f\x24\xc9\xc7\x5e\xc9\x5b\xf2\xdc\xc2\xd7\xe9\x07\xa3\x0f\xef\x0f\xf8\xb9\x0e\x1d\x9c\xdb\x1b\x59\x1e\x9f\xd1\x05\x97\x8d\x55\xc1\x5d\x7c\xb5\x7f\xff\xfd\xea\x94\x3d\x94\x83\xbf\x01\x00\x00\xff\xff\x83\x7c\xe4\x65\xfb\x06\x00\x00"

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

	info := bindataFileInfo{name: "db/migrations/20240310083405_init.sql", size: 1787, mode: os.FileMode(438), modTime: time.Unix(1710836594, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _dbMigrations20240312074226_add_roleSql = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd2\xd5\x55\xc8\xcd\x4c\x2f\x4a\x2c\x49\xb5\x2a\x2d\xe0\x72\xf4\x09\x71\x0d\x52\x08\x71\x74\xf2\x71\x55\x28\x2d\x4e\x2d\x2a\x56\x70\x74\x71\x51\x70\xf6\xf7\x09\xf5\xf5\x53\x50\x2a\xca\xcf\x49\x55\x52\xf0\xf4\x0b\x51\xf0\xf3\x0f\x51\xf0\x0b\xf5\xf1\x51\x70\x71\x75\x73\x0c\xf5\x09\x51\x30\xb0\xe6\xe2\x42\x32\x2a\x25\xbf\x3c\x8f\x0b\x10\x00\x00\xff\xff\xcd\x09\xbb\x1d\x5b\x00\x00\x00"

func dbMigrations20240312074226_add_roleSqlBytes() ([]byte, error) {
	return bindataRead(
		_dbMigrations20240312074226_add_roleSql,
		"db/migrations/20240312074226_add_role.sql",
	)
}

func dbMigrations20240312074226_add_roleSql() (*asset, error) {
	bytes, err := dbMigrations20240312074226_add_roleSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "db/migrations/20240312074226_add_role.sql", size: 91, mode: os.FileMode(438), modTime: time.Unix(1710518157, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _dbMigrations20240315173713_add_tokensSql = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd2\xd5\x55\xc8\xcd\x4c\x2f\x4a\x2c\x49\xb5\x2a\x2d\xe0\x72\x0e\x72\x75\x0c\x71\x55\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\xf0\xf3\x0f\x51\x70\x8d\xf0\x0c\x0e\x09\x56\x28\xc9\xcf\x4e\xcd\x2b\x56\xd0\xe0\x52\x50\x50\xca\x4c\x51\x52\x70\xf2\x74\xf7\xf4\x0b\x51\x08\x08\xf2\xf4\x75\x0c\x8a\x54\xf0\x76\x8d\xd4\x01\x49\x81\x95\x29\x29\x84\xb8\x46\x84\x80\x75\xfb\x85\xfa\xf8\x70\x69\x5a\x73\x71\x21\x59\x94\x92\x5f\x9e\xc7\x05\x08\x00\x00\xff\xff\xcc\xb1\x9d\x99\x79\x00\x00\x00"

func dbMigrations20240315173713_add_tokensSqlBytes() ([]byte, error) {
	return bindataRead(
		_dbMigrations20240315173713_add_tokensSql,
		"db/migrations/20240315173713_add_tokens.sql",
	)
}

func dbMigrations20240315173713_add_tokensSql() (*asset, error) {
	bytes, err := dbMigrations20240315173713_add_tokensSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "db/migrations/20240315173713_add_tokens.sql", size: 121, mode: os.FileMode(438), modTime: time.Unix(1711352444, 0)}
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
	"db/migrations/20240310083405_init.sql":       dbMigrations20240310083405_initSql,
	"db/migrations/20240312074226_add_role.sql":   dbMigrations20240312074226_add_roleSql,
	"db/migrations/20240315173713_add_tokens.sql": dbMigrations20240315173713_add_tokensSql,
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
			"20240310083405_init.sql":       &bintree{dbMigrations20240310083405_initSql, map[string]*bintree{}},
			"20240312074226_add_role.sql":   &bintree{dbMigrations20240312074226_add_roleSql, map[string]*bintree{}},
			"20240315173713_add_tokens.sql": &bintree{dbMigrations20240315173713_add_tokensSql, map[string]*bintree{}},
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
