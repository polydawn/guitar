package commands

import (
	. "fmt"
	"polydawn.net/guitar/stream"
)

type ExportCmdOpts struct {
}

//Transforms a container
func (opts *ExportCmdOpts) Execute(args []string) error {
	_ = stream.ExportToFilesystem //import for compile checking
	return Errorf("Not implemented yet")
}
