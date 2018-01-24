package main

import (
	"fmt"
	"os"

	"bazil.org/fuse"
	"github.com/jsternberg/influxfs/influxfs"
)

func realMain() int {
	conn, err := fuse.Mount("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to mount fuse filesystem: %s\n", err)
		return 1
	}
	defer conn.Close()

	<-conn.Ready
	fs := influxfs.New()
	if err := fs.Serve(conn); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(realMain())
}
