package influxfs

import (
	"errors"
	"fmt"
	"os"

	"bazil.org/fuse"
)

type FileSystem struct {
}

func New() *FileSystem {
	return &FileSystem{}
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
