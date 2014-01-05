package stream

//Given a tar, export the files and metadata regarding the file's properites.
//This is because git does not store everything we need, so we store that separately in a metadata file.

import (
	"archive/tar"
	"encoding/json"
	. "fmt"
	"io"
	"os"
	"path/filepath"
	"polydawn.net/guitar/conf"
	"polydawn.net/guitar/format"
)

//Given a reader, export the contents to the filesystem.
func ExportFromReaderToFilesystem(r *io.Reader, fsPath string, settings *conf.Settings) error {
	//Connect a tar reader
	stream := tar.NewReader(*r)

	return ExportToFilesystem(stream, fsPath, settings)
}

//Given a reader to a tar stream, export the contents to the filesystem.
func ExportToFilesystem(r *tar.Reader, fsPath string, settings *conf.Settings) error {
	//A set of headers. These are cached then sorted before writing as metadata.
	//This ensures the same filesystem will always the same metadata, because tar archives do not guarantee ordering.
	headers := make([]*format.Header, 0)

	//A closure that exports each file
	forEachFile := func(stream *tar.Reader, hdr *tar.Header) error {

		//Convert the tar header into a human-friendly format, then cache for later
		export, err := format.Export(hdr, settings)
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
				os.MkdirAll(filepath.Join(fsPath, hdr.Name), 0755)

			//regular file
			case tar.TypeReg, tar.TypeRegA:
				filename := filepath.Join(fsPath, hdr.Name)
				folder   := filepath.Join(fsPath, filepath.Dir(hdr.Name))

				//Create any folders
				err := os.MkdirAll(filepath.Join(folder), 0755)
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
	metaFilename := filepath.Join(fsPath, ".guitar")
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
func Export(r *tar.Reader, fn func(*tar.Reader, *tar.Header) error) error {
	// Iterate through the files in the archive.
	for {
		//Read the header for the next file.
		hdr, err := r.Next()

		//Check for errors
		if err == io.EOF {
			break //End of the archive
		} else if err != nil {
			return Errorf("Error extracting files from tar: " + err.Error())
		}

		//Run
		err = fn(r, hdr)
		if (err != nil) {
			return Errorf("Error exporting file " + hdr.Name + ": " + err.Error())
		}
	}

	return nil
}
