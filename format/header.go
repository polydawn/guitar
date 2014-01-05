package format

import (
	"archive/tar"
	. "fmt"
	"strconv"
	"time"
)

//A subset of the tar package's header. We don't need every field.
//Also added some annotations to some commonly-zero fields, reducing noise.
//	See: http://golang.org/pkg/archive/tar/#Header
type Header struct {
	Name       string                        // name of header file entry
	Type       string                        // type of header entry
	Mode       int64                         // permission and mode bits
	ModTime    time.Time `json:",omitempty"` // modified time
	Uid        int       `json:",omitempty"` // user id of owner
	Gid        int       `json:",omitempty"` // group id of owner
	Linkname   string    `json:",omitempty"` // target name of link
	Devmajor   int64     `json:",omitempty"` // major number of character or block device
	Devminor   int64     `json:",omitempty"` // minor number of character or block device
}


//Exports a tar header into one of our own.
func Export(hdr *tar.Header) (*Header, error) {
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
		case tar.TypeReg, tar.TypeRegA: typeLetter = "F" // regular file (LOSSILY DROP TypeRegA. Who cares? Not me!)

		default: // unknown filetype, bad news bears
			return nil, Errorf("WAT: Unexpected TypeFlag " + string(hdr.Typeflag))
	}

	//Copy header values
	converted := &Header{
		Name: hdr.Name,
		Type: typeLetter,
		Mode: int64(mode), //cast octal-formatted int to int64
		ModTime: hdr.ModTime.UTC(),
		Uid: hdr.Uid,
		Gid: hdr.Gid,
		Linkname: hdr.Linkname,
		Devmajor: hdr.Devmajor,
		Devminor: hdr.Devminor,
	}

	return converted, nil
}

//Imports our custom header into a tar header.
func Import(hdr *Header) (*tar.Header, error) {

	//Convert octal file mode back to base 10
	mode, err := strconv.ParseInt(strconv.Itoa(int(hdr.Mode)), 8, 64)
	if err != nil {
		return nil, Errorf("Error converting " + string(hdr.Mode) + " to integer: " + err.Error())
	}

	//Convert human-friendly character back to tar's integer typeflag.
	var flag byte = tar.TypeReg
	switch hdr.Type {
		case "H": flag = tar.TypeLink    // hard link
		case "C": flag = tar.TypeChar    // character device node
		case "B": flag = tar.TypeBlock   // block device node
		case "P": flag = tar.TypeFifo    // fifo node
		case "S": flag = tar.TypeSymlink // symbolic link
		case "D": flag = tar.TypeDir     // directory
		case "F": flag = tar.TypeReg     // regular file

		default: // unknown filetype, bad news bears
			return nil, Errorf("WAT: Unexpected type " + hdr.Type)
	}

	//Copy header values
	converted := &tar.Header {
		Name: hdr.Name,
		Typeflag: flag,
		Mode: mode,
		ModTime: hdr.ModTime,
		Uid: hdr.Uid,
		Gid: hdr.Gid,
		Linkname: hdr.Linkname,
		Devmajor: hdr.Devmajor,
		Devminor: hdr.Devminor,
	}

	return converted, nil
}
