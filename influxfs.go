package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"syscall"

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
	querier := client.Querier()
	if err := querier.Execute("CREATE DATABASE fuse"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not create database: %s\n", err)
		return 1
	}

	writer := client.Writer()
	writer.Database = "fuse"
	bufWriter := influxdb.NewTimedWriter(influxdb.NewBufferedWriter(writer), time.Second)
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

	<-conn.Ready

	if err := conn.MountError; err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to connect to the fuse filesystem: %s\n", err)
		return 1
	}

	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt)

	done := make(chan struct{})
	defer close(done)
	go func() {
		for {
			select {
			case <-signalCh:
				syscall.Unmount(dest, 0)
			case <-done:
				return
			}
		}
	}()

	err = fs.Serve(conn, filesystem)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to serve fuse filesystem: %s\n", err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(realMain())
}
