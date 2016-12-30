package steps

import (
	"fmt"

	"github.com/DimensionDataResearch/packer-plugins-ddcloud/builders/customerimage/config"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers"
	"github.com/mitchellh/multistep"
)

// ResolveVLAN is the step that resolves the target VLAN by name from CloudControl.
type ResolveVLAN struct{}

// Run is called to perform the step's action.
//
// The return value determines whether multi-step sequences should continue or halt.
func (step *ResolveVLAN) Run(stateBag multistep.StateBag) multistep.StepAction {
	state := helpers.ForStateBag(stateBag)
	ui := state.GetUI()

	settings := state.GetConfig().(*config.Settings)
	client := state.GetClient()
	networkDomain := state.GetNetworkDomain()

	vlan, err := client.GetVLANByName(
		settings.VLANName,
		networkDomain.ID,
	)
	if err != nil {
		ui.Error(err.Error())

		return multistep.ActionHalt
	}
	if vlan == nil {
		ui.Error(fmt.Sprintf(
			"Unable to find VLAN '%s' in network domain '%s' ('%s') in datacenter '%s'.",
			settings.VLANName,
			networkDomain.Name,
			networkDomain.ID,
			settings.DatacenterID,
		))

		return multistep.ActionHalt
	}

	state.SetVLAN(vlan)

	return multistep.ActionContinue
}

// Cleanup is called in reverse order of the steps that have run
// and allow steps to clean up after themselves. Do not assume if this
// ran that the entire multi-step sequence completed successfully. This
// method can be ran in the face of errors and cancellations as well.
//
// The parameter is the same "state bag" as Run, and represents the
// state at the latest possible time prior to calling Cleanup.
func (step *ResolveVLAN) Cleanup(state multistep.StateBag) {
}

var _ multistep.Step = &ResolveVLAN{}
