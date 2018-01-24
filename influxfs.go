package main

import (
	"errors"
	"fmt"
	"os"

	"bazil.org/fuse"
)

func realMain() int {
	conn, err := fuse.Mount("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to mount fuse filesystem: %s\n", err)
		return 1
	}
	defer conn.Close()

	<-conn.Ready
	for {
		req, err := conn.ReadRequest()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not read request: %s\n", err)
			continue
		}
		req.RespondError(errors.New("unimplemented"))
	}
	return 0
}

func main() {
	os.Exit(realMain())
}
