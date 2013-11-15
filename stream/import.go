package stream

// He was a bear, all along.

import (
	"bufio"
	"encoding/json"
	. "fmt"
	"io"
	"os"
	"path/filepath"
	"github.com/dotcloud/tar" // Dotcloud's fork of the core tar library
	"polydawn.net/guitar/format"
)

//Given a writer to a tar stream, import from the filesystem.
func ImportFromFilesystem(w io.WriteCloser, basePath string) error {
	defer w.Close()

	//Open a tar writer
	out := tar.NewWriter(w)

	//Get a sane path.
	tempDir, err := filepath.Abs(filepath.Clean(basePath))
	if (err != nil) { return err }
	base, err := filepath.EvalSymlinks(tempDir)
	if (err != nil) { return err }

	//Open the metadata file
	metaFilename := filepath.Join(base, ".guitar")
	metafile, err := os.Open(metaFilename)
	defer metafile.Close()
	if (err != nil) { return err }
	buffer := bufio.NewReader(metafile)

	//Cannot hard link to files that don't exist yet.
	//We buffer them here, and apply them last.
	hardLinks := make([]*tar.Header, 0)

	//Read each line
	Println("Importing files...")
	for {
		//Get bytes until next newline
		line, err := buffer.ReadBytes('\n')
		if err == io.EOF {
			break //Last line has been read; done importing
		} else if (err != nil) {
			return err
		}

		//Decode that line's JSON
		var hdr *format.Header
		err = json.Unmarshal(line, &hdr)
		if (err != nil) { return err }

		//Convert to tar header format, write to the stream
		header, err := format.Import(hdr)
		if (err != nil) { return err }

		//Take action on various file types
		switch header.Typeflag {
			case tar.TypeLink: //hard link
				hardLinks = append(hardLinks, header) //save for later

			case tar.TypeReg, tar.TypeRegA: //regular file
				filename := filepath.Join(base, header.Name)

				//Open file, get size
				file, err := os.Open(filename)
				if (err != nil) { return err }
				info, err := file.Stat()
				if (err != nil) { return err }

				//Set header size
				header.Size = info.Size()

				//Write the header
				err = out.WriteHeader(header)
				if (err != nil) { return err }

				//Write the file
				_, err = io.Copy(out, file)
				if err != nil {
					return Errorf("Could not write file " + filename + ": " + err.Error())
				}
				file.Close()
			default: //everything
				//Write the header
				err = out.WriteHeader(header)
				if (err != nil) { return err }
		}
	}

	//Write all the hardlinks
	for _, hdr := range hardLinks {
		err = out.WriteHeader(hdr)
		if (err != nil) { return err }
	}

	return nil
}
