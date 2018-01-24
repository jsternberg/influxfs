package main

import (
	"fmt"
	"os"
	"time"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
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

	source := args[0]
	dest := args[1]

	filesystem := influxfs.New(source, bufWriter)
	conn, err := fuse.Mount(
		dest,
		fuse.FSName("influxfs"),
		fuse.LocalVolume(),
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to mount fuse filesystem: %s\n", err)
		return 1
	}
	defer conn.Close()

	err = fs.Serve(conn, filesystem)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to serve fuse filesystem: %s\n", err)
		return 1
	}

	<-conn.Ready

	if err := conn.MountError; err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to connect to the fuse filesystem: %s\n", err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(realMain())
}
