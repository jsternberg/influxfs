package influxfs

import (
	"fmt"

	"bazil.org/fuse"
	"golang.org/x/net/context"
)

type File struct {
	name string
}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 2
	a.Mode = 0444
	a.Size = uint64(100)
	return nil
}

func (f *File) Getxattr(ctx context.Context, req *fuse.GetxattrRequest, resp *fuse.GetxattrResponse) error {
	return nil
}

func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
	return []byte("Hello its me, thanks for opening"), nil
}

func (f *File) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	fmt.Println("Not really writing, thank you")

	return nil
}

func (f *File) Flush(ctx context.Context, req *fuse.FlushRequest) error {
	fmt.Println("I was supposing to flush but i didnt")
	return nil
}
