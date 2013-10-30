package format

import (
	. "fmt"
	"strconv"
	"time"
	"github.com/dotcloud/tar" // Dotcloud's fork of the core tar library
)

//A subset of the tar package's header. We don't need every field.
//Also added some annotations to some commonly-zero fields, reducing noise.
//	See: http://golang.org/pkg/archive/tar/#Header
type Header struct {
	Name       string                        // name of header file entry
	Mode       int64                         // permission and mode bits
	Uid        int       `json:",omitempty"` // user id of owner
	Gid        int       `json:",omitempty"` // group id of owner
	ModTime    time.Time                     // modified time
	Typeflag   byte                          // type of header entry
	Linkname   string    `json:",omitempty"` // target name of link
	Devmajor   int64     `json:",omitempty"` // major number of character or block device
	Devminor   int64     `json:",omitempty"` // minor number of character or block device
}


//Exports a tar header into one of our own.
//Changes the file mode to octal for human-readability and
func HeaderExport(hdr *tar.Header) (*Header, error) {
	//Convert integer file mode to octal, because it's 100x more useful that way.
	//Definitely the best way to convert an integer's base EVAR, more string ops desired
	mode, err := strconv.Atoi(strconv.FormatInt(hdr.Mode, 8))
	if err != nil {
		return nil, Errorf("Error converting " + string(hdr.Mode) + " to octal: " + err.Error())
	}

	//Copy header values
	converted := &Header{
		Name: hdr.Name,
		Mode: int64(mode), //cast octal-formatted int to int64
		Uid: hdr.Uid,
		Gid: hdr.Gid,
		ModTime: hdr.ModTime,
		Typeflag: hdr.Typeflag,
		Linkname: hdr.Linkname,
		Devmajor: hdr.Devmajor,
		Devminor: hdr.Devminor,
	}

	return converted, nil
}

func HeaderImport(hdr *Header) *tar.Header {
	return nil
}
