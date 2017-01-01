package artifacts

import (
	"fmt"

	"github.com/mitchellh/packer/packer"
)

// RemoteOVFPackage represents an OVF package in CloudControl.
type RemoteOVFPackage struct {
	FTPSHostName       string
	PackagePrefix      string
	BuilderID          string
	deletePackageFiles func() error
}

// BuilderId returns the ID of the builder that was used to create the artifact.
func (artifact *RemoteOVFPackage) BuilderId() string {
	return artifact.BuilderID
}

// Files determines the set of files that comprise the artifact.
// If an artifact is not made up of files, then this will be empty.
func (artifact *RemoteOVFPackage) Files() []string {
	return []string{}
}

// Id gets the ID for the artifact.
// In this case, it's "FTPSHostName/PackagePrefix".
func (artifact *RemoteOVFPackage) Id() string {
	return fmt.Sprintf("%s/%s",
		artifact.FTPSHostName,
		artifact.PackagePrefix,
	)
}

// Returns human-readable output that describes the artifact created.
// This is used for UI output. It can be multiple lines.
func (artifact *RemoteOVFPackage) String() string {
	return fmt.Sprintf("OVF package with prefix '%s' on FTPS host '%s'.",
		artifact.PackagePrefix,
		artifact.FTPSHostName,
	)
}

// State allows the caller to ask for builder specific state information
// relating to the artifact instance.
func (artifact *RemoteOVFPackage) State(name string) interface{} {
	return nil // No specific state.
}

// Destroy deletes the artifact. Packer calls this for various reasons,
// such as if a post-processor has processed this artifact and it is
// no longer needed.
func (artifact *RemoteOVFPackage) Destroy() error {
	if artifact.deletePackageFiles == nil {
		return nil // Already deleted.
	}

	err := artifact.deletePackageFiles()
	artifact.deletePackageFiles = nil

	return err
}

var _ packer.Artifact = &RemoteOVFPackage{}
