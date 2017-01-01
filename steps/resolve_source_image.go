package steps

import (
	"fmt"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/artifacts"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/builders/customerimage/config"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers"
	"github.com/mitchellh/multistep"
)

// ResolveSourceImage is the step that resolves the source image from CloudControl.
type ResolveSourceImage struct {
	// If true, then the source image must be a customer image.
	MustBeCustomerImage bool
}

// Run is called to perform the step's action.
//
// The return value determines whether multi-step sequences should continue or halt.
func (step *ResolveSourceImage) Run(stateBag multistep.StateBag) multistep.StepAction {
	state := helpers.ForStateBag(stateBag)
	ui := state.GetUI()

	settings := state.GetSettings().(*config.Settings)
	client := state.GetClient()
	networkDomain := state.GetNetworkDomain()

	var image compute.Image

	osImage, err := client.FindOSImage(settings.SourceImage, networkDomain.DatacenterID)
	if err != nil {
		ui.Error(err.Error())

		return multistep.ActionHalt
	}
	if !step.MustBeCustomerImage && osImage != nil {
		image = osImage
	} else {
		// Fall back to customer image.
		customerImage, err := client.FindCustomerImage(settings.SourceImage, networkDomain.DatacenterID)
		if err != nil {
			ui.Error(err.Error())

			return multistep.ActionHalt
		}
		if customerImage != nil {
			image = customerImage
		}
	}

	if image == nil {
		ui.Error(fmt.Sprintf(
			"Image '%s' not found in datacenter '%s'.",
			settings.SourceImage,
			networkDomain.DatacenterID,
		))

		return multistep.ActionHalt
	}

	state.SetSourceImage(image)
	state.SetSourceImageArtifact(&artifacts.Image{
		Image: image,
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
func (step *ResolveSourceImage) Cleanup(state multistep.StateBag) {
}

var _ multistep.Step = &ResolveSourceImage{}
