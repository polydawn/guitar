package commands

import (
	. "fmt"
)

const Version = "0.0.0"

type VersionCmdOpts struct { }

//Version command
func (opts *VersionCmdOpts) Execute(args []string) error {
	Println("guitar version", Version)
	return nil
}
