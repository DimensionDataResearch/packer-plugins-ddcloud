package steps

import (
	"fmt"

	"github.com/DimensionDataResearch/packer-plugins-ddcloud/builders/customerimage/config"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers"
	"github.com/mitchellh/multistep"
)

// CheckTargetImage is the step that ensures the target image does not already exist in CloudControl.
type CheckTargetImage struct{}

// Run is called to perform the step's action.
//
// The return value determines whether multi-step sequences should continue or halt.
func (step *CheckTargetImage) Run(stateBag multistep.StateBag) multistep.StepAction {
	state := helpers.ForStateBag(stateBag)
	ui := state.GetUI()

	settings := state.GetConfig().(*config.Settings)
	client := state.GetClient()

	targetImage, err := client.FindCustomerImage(settings.TargetImage, settings.DatacenterID)
	if err != nil {
		ui.Error(err.Error())

		return multistep.ActionHalt
	}

	if targetImage != nil {
		ui.Error(fmt.Sprintf(
			"Target image '%s' already exists in datacenter '%s'.",
			settings.TargetImage,
			settings.DatacenterID,
		))

		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

// Cleanup is called in reverse order of the steps that have run
// and allow steps to clean up after themselves. Do not assume if this
// ran that the entire multi-step sequence completed successfully. This
// method can be ran in the face of errors and cancellations as well.
//
// The parameter is the same "state bag" as Run, and represents the
// state at the latest possible time prior to calling Cleanup.
func (step *CheckTargetImage) Cleanup(state multistep.StateBag) {
}

var _ multistep.Step = &CheckTargetImage{}
