package influxfs

import (
	"bazil.org/fuse/fs"
	"github.com/influxdata/influxdb-client"
)

type FileSystem struct {
	dir    string
	writer influxdb.Writer
}

func New(dir string, writer influxdb.Writer) *FileSystem {
	return &FileSystem{
		dir:    dir,
		writer: writer,
	}
}

func (g *FileSystem) Root() (fs.Node, error) {
	return &Dir{path: g.dir}, nil
}
