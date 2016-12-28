package ddcloud

import "github.com/mitchellh/packer/packer"
import "github.com/DimensionDataResearch/go-dd-cloud-compute/compute"

// BuilderID is the unique Id for the ddcloud builder
const BuilderID = "dimension-data-research.ddcloud"

// Builder is the Builder plugin for Packer.
type Builder struct {
	client *compute.Client
}

// Prepare the plugin to run.
func (builder *Builder) Prepare(configuration ...interface{}) (warnings []string, err error) {
	return nil, nil
}

// Run the plugin.
func (builder *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	return nil, nil
}

// Cancel plugin execution.
func (builder *Builder) Cancel() {
	if builder.client != nil {
		builder.client.Cancel()
	}
}
