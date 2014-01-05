package conf

type Settings struct {
	// If true, zero all modtime properties from the fs stream when generating the metadata file.
	Epoch bool
}
