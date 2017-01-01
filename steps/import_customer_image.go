package steps

import (
	"fmt"
	"time"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/artifacts"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers"
	"github.com/mitchellh/multistep"
)

// ImportCustomerImage is the step that imports a customer image from an OVF package.
type ImportCustomerImage struct {
	// The name of the target image to create.
	TargetImageName string

	// The Id of the datacenter where the customer image will be created.
	DatacenterID string

	// The prefix for the OVF package files.
	OVFPackagePrefix string

	// Configure the customer image to prevent guess OS customisation?
	PreventGuestOSCustomization bool
}

// Run is called to perform the step's action.
//
// The return value determines whether multi-step sequences should continue or halt.
func (step *ImportCustomerImage) Run(stateBag multistep.StateBag) multistep.StepAction {
	state := helpers.ForStateBag(stateBag)
	ui := state.GetUI()

	client := state.GetClient()

	ui.Message(fmt.Sprintf(
		"Create customer image '%s' in datacenter '%s' from OVF package '%s'.",
		step.TargetImageName,
		step.DatacenterID,
		step.OVFPackagePrefix,
	))

	imageID, err := client.ImportCustomerImage(
		step.TargetImageName,
		step.TargetImageName+" (created by Packer).",
		step.PreventGuestOSCustomization,
		step.OVFPackagePrefix,
		step.DatacenterID,
	)
	if err != nil {
		ui.Error(err.Error())

		return multistep.ActionHalt
	}

	ui.Message(fmt.Sprintf(
		"Import of customer image '%s' ('%s') in progress...",
		step.TargetImageName,
		imageID,
	))

	resource, err := client.WaitForDeploy(compute.ResourceTypeCustomerImage, imageID, 30*time.Minute)
	if err != nil {
		ui.Error(err.Error())

		return multistep.ActionHalt
	}

	image := resource.(*compute.CustomerImage)

	ui.Message(fmt.Sprintf(
		"Import of customer image '%s' ('%s) complete.",
		image.Name,
		image.ID,
	))

	state.SetTargetImage(image)
	state.SetTargetImageArtifact(&artifacts.Image{
		Image:     image,
		BuilderID: state.GetBuilderID(),
	})

	return multistep.ActionContinue
}

// Cleanup is called in reverse order of the steps that have run
// and allow steps to clean up after themselves. Do not assume if this
// ran that the entire multi-step sequence completed successfully. This
// method can be ran in the face of errors and cancellations as well.
//
// The parameter is the same "state bag" as Run, and represents the
// state at the latest possible time prior to calling Cleanup.
func (step *ImportCustomerImage) Cleanup(state multistep.StateBag) {
}
