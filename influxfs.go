package main

import (
	"fmt"
	"os"
	"time"

	"bazil.org/fuse"
	"github.com/influxdata/influxdb-client"
	"github.com/jsternberg/influxfs/influxfs"
	flag "github.com/spf13/pflag"
)

func realMain() int {
	influxdbUrl := flag.StringP("host", "H", "http://localhost:8086", "Host to report statistics too")
	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Error: Requires exactly two arguments\n")
		return 1
	}

	client, err := influxdb.NewClient(*influxdbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid InfluxDB URL: %s\n", err)
		return 1
	}
	bufWriter := influxdb.NewTimedWriter(influxdb.NewBufferedWriter(client.Writer()), time.Second)
	defer bufWriter.Flush()

	fs := influxfs.New(args[0], bufWriter)
	conn, err := fuse.Mount(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to mount fuse filesystem: %s\n", err)
		return 1
	}
	defer conn.Close()

	<-conn.Ready
	if err := fs.Serve(conn); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(realMain())
}
