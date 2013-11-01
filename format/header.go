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
	Type       string                        // type of header entry
	Mode       int64                         // permission and mode bits
	ModTime    time.Time                     // modified time
	Uid        int       `json:",omitempty"` // user id of owner
	Gid        int       `json:",omitempty"` // group id of owner
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

	//The tar's typeflag is based on integers.
	//Because of this, it's not immediately obvious to a human was sort of entry it is ("typeflag 54?")
	//We convert this to an upper-case rune for readability. F: file, D: dir, S: symlink, etc.
	typeLetter := "Z"
	switch hdr.Typeflag {
		case tar.TypeLink:              typeLetter = "H" // hard link
		case tar.TypeChar:              typeLetter = "C" // character device node
		case tar.TypeBlock:             typeLetter = "B" // block device node
		case tar.TypeFifo:              typeLetter = "P" // fifo node
		case tar.TypeSymlink:           typeLetter = "S" // symbolic link
		case tar.TypeDir:               typeLetter = "D" // directory
		case tar.TypeReg, tar.TypeRegA: typeLetter = "F" //regular file

		default: // unknown filetype, bad news bears
			return nil, Errorf("WAT: Unexpected TypeFlag " + string(hdr.Typeflag))
	}

	//Copy header values
	converted := &Header{
		Name: hdr.Name,
		Mode: int64(mode), //cast octal-formatted int to int64
		Uid: hdr.Uid,
		Gid: hdr.Gid,
		ModTime: hdr.ModTime,
		Type: typeLetter,
		Linkname: hdr.Linkname,
		Devmajor: hdr.Devmajor,
		Devminor: hdr.Devminor,
	}

	return converted, nil
}

func HeaderImport(hdr *Header) *tar.Header {
	return nil
}
