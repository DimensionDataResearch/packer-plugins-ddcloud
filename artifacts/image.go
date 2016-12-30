package artifacts

import (
	"fmt"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
	"github.com/mitchellh/packer/packer"
)

// Image represents a CloudControl image as a Packer Artifact.
type Image struct {
	Server        compute.Server
	NetworkDomain compute.NetworkDomain
	Image         compute.CustomerImage
	BuilderID     string
	deleteImage   func() error
}

// BuilderId returns the ID of the builder that was used to create the artifact.
func (artifact *Image) BuilderId() string {
	return artifact.BuilderID
}

// Files determines the set of files that comprise the artifact.
// If an artifact is not made up of files, then this will be empty.
func (artifact *Image) Files() []string {
	return []string{}
}

// Id gets the ID for the artifact.
// In this case, it's the image Id.
func (artifact *Image) Id() string {
	return artifact.Image.ID
}

// Returns human-readable output that describes the artifact created.
// This is used for UI output. It can be multiple lines.
func (artifact *Image) String() string {
	return fmt.Sprintf("Customer image '%s' ('%s') in datacenter '%s'.",
		artifact.Image.Name,
		artifact.Image.ID,
		artifact.Image.DataCenterID,
	)
}

// State allows the caller to ask for builder specific state information
// relating to the artifact instance.
func (artifact *Image) State(name string) interface{} {
	return nil // No specific state yet.
}

// Destroy deletes the artifact. Packer calls this for various reasons,
// such as if a post-processor has processed this artifact and it is
// no longer needed.
func (artifact *Image) Destroy() error {
	if artifact.deleteImage == nil {
		return nil // Already deleted.
	}

	err := artifact.deleteImage()
	artifact.deleteImage = nil

	return err
}

var _ packer.Artifact = &Image{}
