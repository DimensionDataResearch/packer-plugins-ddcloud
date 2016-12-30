package main

import (
	"github.com/mitchellh/packer/packer"
)

// PostProcessor is the customer image export post-processor plugin for Packer.
type PostProcessor struct{}

// Configure is responsible for setting up configuration, storing the state for later,
// and returning and errors, such as validation errors.
func (postProcessor *PostProcessor) Configure(...interface{}) error {
	return nil
}

// PostProcess takes a previously created Artifact and produces another Artifact.
//
// If an error occurs, it should return that error.
// If `keep` is to true, then the previous artifact is forcibly kept.
func (postProcessor *PostProcessor) PostProcess(ui packer.Ui, sourceArtifact packer.Artifact) (destinationArtifact packer.Artifact, keep bool, err error) {
	return
}

var _ packer.PostProcessor = &PostProcessor{}
