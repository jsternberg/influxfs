package influxfs

import (
	"os"

	"path/filepath"

	"io/ioutil"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/influxdata/influxdb-client"
	"golang.org/x/net/context"
)

type Dir struct {
	path   string
	writer influxdb.Writer
}

func (dir *Dir) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	return dir, nil
}

func (dir *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	st, err := os.Stat(dir.path)
	if err != nil {
		return err
	}
	a.Mode = st.Mode()
	a.Size = uint64(st.Size())
	return nil
}

func (dir *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {

	path := filepath.Join(dir.path, name)
	st, err := os.Stat(path)
	if err != nil {
		return nil, fuse.ENOENT
	}

	if st.Mode()&os.ModeDir != 0 {
		return &Dir{path: path, writer: dir.writer}, nil
	}
	return &File{name: path}, nil
}

func (dir *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	dirs, err := ioutil.ReadDir(dir.path)
	if err != nil {
		return nil, err
	}

	de := make([]fuse.Dirent, len(dirs))
	for i, dir := range dirs {
		de[i] = fuse.Dirent{
			Name: dir.Name(),
			Type: fuse.DT_File,
		}
		if dir.Mode()&os.ModeDir != 0 {
			de[i].Type = fuse.DT_Dir
		}
	}
	return de, nil
}

func (dir *Dir) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	de := &File{}
	return de, de, nil
}

func (dir *Dir) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	trace(dir.writer, req.Hdr(), "mkdir", map[string]interface{}{"name": req.Name})
	path := filepath.Join(dir.path, req.Name)
	if err := os.Mkdir(path, req.Mode); err != nil {
		return nil, err
	}
	return &Dir{path: path, writer: dir.writer}, nil
}
