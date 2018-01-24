package influxfs

import (
	"errors"
	"fmt"
	"os"

	"bazil.org/fuse"
	"github.com/influxdata/influxdb-client"
)

type FileSystem struct {
	dir    string
	writer *influxdb.TimedWriter
}

func New(dir string) *FileSystem {
	return &FileSystem{
		dir: dir,
	}
}

func (fs *FileSystem) Serve(conn *fuse.Conn) error {
	for {
		req, err := conn.ReadRequest()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not read request: %s\n", err)
			continue
		}
		req.RespondError(errors.New("unimplemented"))
	}
	return nil
}
