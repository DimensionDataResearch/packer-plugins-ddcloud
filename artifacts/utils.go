package artifacts

import (
	"path"

	"github.com/mitchellh/packer/packer"
)

// GetFirstFileWithExtension retrieves the first file with the specified extension from the artifact's list of files.
func GetFirstFileWithExtension(extension string, artifact packer.Artifact) string {
	for _, file := range artifact.Files() {
		if path.Ext(file) == extension {
			return file
		}
	}

	return ""
}
