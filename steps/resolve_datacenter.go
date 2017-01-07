package steps

import (
	"fmt"

	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers"
	"github.com/mitchellh/multistep"
)

// ResolveDatacenter is the step that resolves the target datacenter from CloudControl.
type ResolveDatacenter struct {
	// The Id of the datacenter to resolve.
	DatacenterID string

	// Should the datacenter being resolved be treated as the source datacenter?
	AsSource bool

	// Should the datacenter being resolved be treated as the target datacenter?
	AsTarget bool
}

// Run is called to perform the step's action.
//
// The return value determines whether multi-step sequences should continue or halt.
func (step *ResolveDatacenter) Run(stateBag multistep.StateBag) multistep.StepAction {
	state := helpers.ForStateBag(stateBag)
	ui := state.GetUI()

	client := state.GetClient()

	ui.Message(fmt.Sprintf(
		"Resolving datacenter '%s'...",
		step.DatacenterID,
	))

	datacenter, err := client.GetDatacenter(step.DatacenterID)
	if err != nil {
		ui.Error(err.Error())

		return multistep.ActionHalt
	}
	if datacenter == nil {
		ui.Error(fmt.Sprintf(
			"Unable to find datacenter '%s'.",
			step.DatacenterID,
		))

		return multistep.ActionHalt
	}
	if step.AsTarget {
		state.SetTargetDatacenter(datacenter)
	}

	ui.Message(fmt.Sprintf(
		"Resolved datacenter '%s'.",
		step.DatacenterID,
	))

	return multistep.ActionContinue
}

// Cleanup is called in reverse order of the steps that have run
// and allow steps to clean up after themselves. Do not assume if this
// ran that the entire multi-step sequence completed successfully. This
// method can be ran in the face of errors and cancellations as well.
//
// The parameter is the same "state bag" as Run, and represents the
// state at the latest possible time prior to calling Cleanup.
func (step *ResolveDatacenter) Cleanup(state multistep.StateBag) {
}

var _ multistep.Step = &ResolveDatacenter{}
