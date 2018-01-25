package influxfs

import (
	"fmt"
	"os"
	"time"

	"bazil.org/fuse"
	"github.com/influxdata/influxdb-client"
)

func (fs *FileSystem) trace(header *fuse.Header, typ string, f map[string]interface{}) {

	fields := make(map[string]interface{}, len(f)+1)
	for key, value := range f {
		fields[key] = value
	}
	fields["id"] = header.ID

	p := influxdb.Point{
		Name: "fsevents",
		Tags: influxdb.Tags{
			influxdb.Tag{Key: "type", Value: typ},
			influxdb.Tag{Key: "gid", Value: fmt.Sprintf("%d", header.Gid)},
			influxdb.Tag{Key: "node", Value: fmt.Sprintf("%d", header.Node)},
			influxdb.Tag{Key: "pid", Value: fmt.Sprintf("%d", header.Pid)},
			influxdb.Tag{Key: "uid", Value: fmt.Sprintf("%d", header.Uid)},
		},
		Fields: fields,
		Time:   time.Now(),
	}

	if _, err := p.WriteTo(fs.writer); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

}
