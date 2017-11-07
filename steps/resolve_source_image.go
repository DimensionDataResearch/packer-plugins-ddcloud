package steps

import (
	"fmt"
	"log"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/artifacts"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers"
	"github.com/mitchellh/multistep"
)

// ResolveSourceImage is the step that resolves the source image from CloudControl.
type ResolveSourceImage struct {
	// The name of the source image.
	ImageName string

	// The name of the target image.
	DatacenterID string

	// If true, then the source image must be a customer image.
	MustBeCustomerImage bool
}

// Run is called to perform the step's action.
//
// The return value determines whether multi-step sequences should continue or halt.
func (step *ResolveSourceImage) Run(stateBag multistep.StateBag) multistep.StepAction {
	state := helpers.ForStateBag(stateBag)
	ui := state.GetUI()

	client := state.GetClient()

	var (
		imageType string
		image     compute.Image
	)

	if step.MustBeCustomerImage {
		imageType = "Customer"
	} else {
		imageType = "any"
	}
	log.Printf(
		"Searching for %s image named '%s' in datacenter '%s'.", imageType, step.ImageName, step.DatacenterID,
	)

	osImage, err := client.FindOSImage(step.ImageName, step.DatacenterID)
	if err != nil {
		ui.Error(err.Error())

		return multistep.ActionHalt
	}
	if osImage != nil {
		log.Printf(
			"Found OS image '%s' ('%s') in datacenter '%s'.", step.ImageName, osImage.ID, step.DatacenterID,
		)
	} else {
		log.Printf(
			"No OS image named '%s' found in datacenter '%s'.", step.ImageName, step.DatacenterID,
		)
	}

	if osImage != nil && !step.MustBeCustomerImage {
		log.Printf(
			"Using OS image '%s' in datacenter '%s'.", step.ImageName, step.DatacenterID,
		)

		image = osImage
	} else {
		if osImage != nil && step.MustBeCustomerImage {
			log.Printf(
				"Ignoring OS image '%s' in datacenter '%s' (a Customer image is required for this step).", step.ImageName, step.DatacenterID,
			)
		}

		// Fall back to customer image.
		customerImage, err := client.FindCustomerImage(step.ImageName, step.DatacenterID)
		if err != nil {
			ui.Error(err.Error())

			return multistep.ActionHalt
		}
		if customerImage != nil {
			log.Printf(
				"Found Customer image '%s' ('%s') in datacenter '%s'.", step.ImageName, customerImage.ID, step.DatacenterID,
			)

			image = customerImage
		} else {
			log.Printf(
				"No Customer image named '%s' found in datacenter '%s'.", step.ImageName, step.DatacenterID,
			)
		}
	}

	if image == nil {
		ui.Error(fmt.Sprintf(
			"Unable to find any image named '%s' in datacenter '%s'.",
			step.ImageName,
			step.DatacenterID,
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
