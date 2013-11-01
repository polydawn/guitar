package stream

//Given a tar, export the files and metadata regarding the file's properites.
//This is because git does not store everything we need, so we store that separately in a metadata file.

import (
	"encoding/json"
	. "fmt"
	"io"
	"os"
	"path"
	"github.com/dotcloud/tar" // Dotcloud's fork of the core tar library
	"polydawn.net/guitar/format"
)

//Given a reader to a tar stream, export the contents to the filesystem.
func ExportToFilesystem(r io.Reader, fsPath string) error {
	//A set of headers. These are cached then sorted before writing as metadata.
	//This ensures the same filesystem will always the same metadata, because tar archives do not guarantee ordering.
	headers := make([]*format.Header, 0)

	//A closure that exports each file
	forEachFile := func(stream *tar.Reader, hdr *tar.Header) error {

		//Convert the tar header into a human-friendly format, then cache for later
		export, err := format.Export(hdr)
		if err != nil { return err }
		headers = append(headers, export)

		//Write the file
		switch hdr.Typeflag {

			//No action taken for these types.
			case tar.TypeLink:              // hard link
			case tar.TypeChar:              // character device node
			case tar.TypeBlock:             // block device node
			case tar.TypeFifo:              // fifo node

			// symbolic link
			case tar.TypeSymlink:
				//TODO (see issue #2)

			// directory
			case tar.TypeDir:
				os.MkdirAll(path.Join(fsPath, hdr.Name), os.ModeDir)

			//regular file
			case tar.TypeReg, tar.TypeRegA:
				filename := path.Join(fsPath, hdr.Name)
				folder   := path.Join(fsPath, path.Dir(hdr.Name))

				//Create any folders
				err := os.MkdirAll(path.Join(folder), os.ModeDir)
				if err != nil {
					return Errorf("Could not create folder " + folder + ": " + err.Error())
				}

				//Create file
				file, err := os.Create(filename)
				if err != nil {
					return Errorf("Could not create file " + filename + ": " + err.Error())
				}
				defer file.Close()

				//Write file
				_, err = io.Copy(file, stream)
				if err != nil {
					return Errorf("Could not write file " + filename + ": " + err.Error())
				}

			// unknown filetype, bad news bears
			default:
				return Errorf("WAT: Unexpected TypeFlag " + string(hdr.Typeflag))
		}

		return nil
	}

	//Export files
	Println("Exporting files...")
	err := Export(r, forEachFile)
	if err != nil {
		return err
	}

	//Sort headers
	Println("Exporting metadata...")
	format.SortByName(headers)

	//Open metadata folder
	metaFilename := path.Join(fsPath, ".guitar")
	metaFile, err := os.Create(metaFilename)
	if err != nil {
		return Errorf("Could not create metadata file " + metaFilename + ": " + err.Error())
	}

	//Print the sorted JSON
	for _, hdr := range headers {
		//Convert header to JSON
		out, err := json.Marshal(hdr)
		if err != nil {
			return Errorf("Error JSON encoding file metadata from tar: " + err.Error())
		}

		//Write metadata
		_, err = metaFile.Write(out)
		if err != nil {
			return Errorf("Could not write file " + metaFilename + ": " + err.Error())
		}

		//Write newline to metadata!
		_, err = metaFile.WriteString("\n")
		if err != nil {
			return Errorf("Could not write file " + metaFilename + ": " + err.Error())
		}
	}

	return nil
}

//Given a reader to a tar stream, run a closure for every file with its header.
func Export(r io.Reader, fn func(*tar.Reader, *tar.Header) error) error {
	//Connect a tar reader
	stream := tar.NewReader(r)

	// Iterate through the files in the archive.
	for {
		//Read the header for the next file.
		hdr, err := stream.Next()

		//Check for errors
		if err == io.EOF {
			break //End of the archive
		} else if err != nil {
			return Errorf("Error extracting files from tar: " + err.Error())
		}

		//Run
		err = fn(stream, hdr)
		if (err != nil) {
			return Errorf("Error exporting file " + hdr.Name + ": " + err.Error())
		}
	}

	return nil
}
