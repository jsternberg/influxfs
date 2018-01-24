package influxfs

import (
	"errors"
	"fmt"
	"os"
	"time"

	"bazil.org/fuse"
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

func (fs *FileSystem) Serve(conn *fuse.Conn) error {
	for {
		req, err := conn.ReadRequest()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not read request: %s\n", err)
			continue
		}

		header := req.Hdr()
		p := influxdb.Point{
			Name: "fsevents",
			Tags: influxdb.Tags{
				influxdb.Tag{Key: "gid", Value: fmt.Sprintf("%d", header.Gid)},
				influxdb.Tag{Key: "node", Value: fmt.Sprintf("%d", header.Node)},
				influxdb.Tag{Key: "pid", Value: fmt.Sprintf("%d", header.Pid)},
				influxdb.Tag{Key: "uid", Value: fmt.Sprintf("%d", header.Uid)},
			},

			Fields: map[string]interface{}{
				"id": header.ID,
			},
			Time: time.Now(),
		}
		if _, err := p.WriteTo(fs.writer); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		req.RespondError(errors.New("unimplemented"))
	}
	return nil
}
