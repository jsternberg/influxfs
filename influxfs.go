package main

import (
	"fmt"
	"os"

	"bazil.org/fuse"
	"github.com/jsternberg/influxfs/influxfs"
	flag "github.com/spf13/pflag"
)

func realMain() int {
	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Error: Requires exactly two arguments\n")
		return 1
	}

	fs := influxfs.New(args[0])
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
