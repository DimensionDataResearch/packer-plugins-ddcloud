package main

import (
	"fmt"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
	"github.com/mitchellh/packer/packer"
)

// ServerArtifact represents a CloudControl server as a Packer Artifact.
type ServerArtifact struct {
	Server        compute.Server
	NetworkDomain compute.NetworkDomain
	deleteServer  func() error
}

// BuilderId returns the ID of the builder that was used to create the artifact.
func (artifact *ServerArtifact) BuilderId() string {
	return BuilderID
}

// Files determines the set of files that comprise the artifact.
// If an artifact is not made up of files, then this will be empty.
func (artifact *ServerArtifact) Files() []string {
	return []string{}
}

// Id gets the ID for the artifact.
// In this case, it's the server Id.
func (artifact *ServerArtifact) Id() string {
	return artifact.Server.ID
}

// Returns human-readable output that describes the artifact created.
// This is used for UI output. It can be multiple lines.
func (artifact *ServerArtifact) String() string {
	return fmt.Sprintf("Server '%s' ('%s') in network domain '%s' ('%s').",
		artifact.Server.Name,
		artifact.Server.ID,
		artifact.NetworkDomain.Name,
		artifact.NetworkDomain.ID,
	)
}

// State allows the caller to ask for builder specific state information
// relating to the artifact instance.
func (artifact *ServerArtifact) State(name string) interface{} {
	return nil // No specific state yet.
}

// Destroy deletes the artifact. Packer calls this for various reasons,
// such as if a post-processor has processed this artifact and it is
// no longer needed.
func (artifact *ServerArtifact) Destroy() error {
	if artifact.deleteServer == nil {
		return nil // Already deleted.
	}

	err := artifact.deleteServer()
	artifact.deleteServer = nil

	return err
}

var _ packer.Artifact = &ServerArtifact{}
