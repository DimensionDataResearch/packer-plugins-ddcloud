package artifacts

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/packer/packer"
)

// LocalFiles is an Artifact representing a set of local files.
type LocalFiles struct {
	BaseDirectory string
	FilePaths     []string
	BuilderID     string
}

// BuilderId returns the ID of the builder that was used to create the artifact.
func (artifact *LocalFiles) BuilderId() string {
	return artifact.BuilderID
}

// Files determines the set of files that comprise the artifact.
// If an artifact is not made up of files, then this will be empty.
func (artifact *LocalFiles) Files() []string {
	return artifact.FilePaths
}

// Id gets the ID for the artifact.
// In this case, it's the directory name.
func (artifact *LocalFiles) Id() string {
	return artifact.BaseDirectory
}

// Returns human-readable output that describes the artifact created.
// This is used for UI output. It can be multiple lines.
func (artifact *LocalFiles) String() string {
	result := fmt.Sprintf("Files in local directory '%s':\n",
		artifact.BaseDirectory,
	)

	skipBaseDirectoryChars := len(artifact.BaseDirectory) + 1
	for _, filePath := range artifact.FilePaths {
		result += fmt.Sprintf("- '%s'\n",
			filePath[skipBaseDirectoryChars:],
		)
	}

	return result
}

// State allows the caller to ask for builder specific state information
// relating to the artifact instance.
func (artifact *LocalFiles) State(name string) interface{} {
	return nil // No state.
}

// Destroy deletes the artifact. Packer calls this for various reasons,
// such as if a post-processor has processed this artifact and it is
// no longer needed.
func (artifact *LocalFiles) Destroy() error {
	return nil // TODO: Implement os.Remove()
}

var _ packer.Artifact = &LocalFiles{}

// NewFromFilesInLocalDirectory creates a new LocalFiles artifact from the files in the specified directory.
func NewFromFilesInLocalDirectory(directory string, builderID string) (artifact *LocalFiles, err error) {
	var files []string
	visitFileSystemEntry := func(path string, info os.FileInfo, visitError error) error {
		if !info.IsDir() {
			files = append(files, path)
		}

		return visitError
	}

	err = filepath.Walk(directory, visitFileSystemEntry)
	if err != nil {
		return
	}

	artifact = &LocalFiles{
		BaseDirectory: directory,
		FilePaths:     files,
		BuilderID:     builderID,
	}

	return
}
