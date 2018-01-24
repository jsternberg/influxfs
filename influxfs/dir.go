package influxfs

import (
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

type Dir struct {
}

func (dir *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 1
	a.Mode = os.ModeDir | 0555
	return nil
}

func (dir *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	return nil, fuse.ENOENT
}

func (dir *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {

	files := []string{"first", "second", "third"}

	var de []fuse.Dirent
	for _, file := range files {
		de = append(de, fuse.Dirent{
			Inode: 2,
			Name:  file,
			Type:  fuse.DT_File,
		})
	}
	return de, nil
}

func (dir *Dir) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	de := &File{}
	return de, de, nil
}
