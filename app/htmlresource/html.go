package htmlresource

import (
	"embed"
	"errors"
	"io/fs"
	"path"
	"path/filepath"
	"strings"
)

//go:embed vuefiles/assets
var Static embed.FS

//go:embed  vuefiles/index.html
var Html []byte

type Resource struct {
	fs   embed.FS
	path string
}

func NewResource() *Resource {
	return &Resource{
		fs:   Static,
		path: "vuefiles",
	}
}

func (r *Resource) Open(name string) (fs.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	fullName := filepath.Join(r.path, filepath.FromSlash(path.Clean("/assets/"+name)))
	file, err := r.fs.Open(fullName)

	return file, err
}
