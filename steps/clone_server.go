package steps

import (
	"fmt"
	"time"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/artifacts"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/builders/customerimage/config"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers"
	"github.com/mitchellh/multistep"
)

// CloneServer is the step that clones the target server in CloudControl.
type CloneServer struct{}

// Run is called to perform the step's action.
//
// The return value determines whether multi-step sequences should continue or halt.
func (step *CloneServer) Run(stateBag multistep.StateBag) multistep.StepAction {
	state := helpers.ForStateBag(stateBag)
	ui := state.GetUI()

	settings := state.GetConfig().(*config.Settings)
	client := state.GetClient()
	networkDomain := state.GetNetworkDomain()
	server := state.GetServer()

	ui.Message(fmt.Sprintf(
		"Shutting down server '%s' ('%s')...",
		server.Name,
		server.ID,
	))

	err := client.ShutdownServer(server.ID)
	if err != nil {
		ui.Error(err.Error())

		return multistep.ActionHalt
	}

	resource, err := client.WaitForChange(
		compute.ResourceTypeServer,
		server.ID,
		"Shut down",
		5*time.Minute,
	)
	if err != nil {
		ui.Error(err.Error())

		return multistep.ActionHalt
	}

	server = resource.(*compute.Server)
	state.SetServer(server)

	ui.Message(fmt.Sprintf(
		"Server '%s' ('%s') has been shut down.",
		server.Name,
		server.ID,
	))

	ui.Message(fmt.Sprintf(
		"Cloning server '%s' ('%s')...",
		server.Name,
		server.ID,
	))

	imageID, err := client.CloneServer(
		server.ID,
		settings.TargetImage,
		fmt.Sprintf("%s (created by Packer)", settings.TargetImage),
		false, // preventGuestOSCustomisation
	)
	if err != nil {
		ui.Error(err.Error())

		return multistep.ActionHalt
	}

	resource, err = client.WaitForServerClone(
		imageID,
		15*time.Minute,
	)
	if err != nil {
		ui.Error(err.Error())

		return multistep.ActionHalt
	}

	customerImage := resource.(*compute.CustomerImage)
	state.SetTargetImage(customerImage)

	ui.Message(fmt.Sprintf(
		"Cloned server '%s' ('%s') to customer image '%s' ('%s') in datacenter '%s'.",
		server.Name,
		server.ID,
		customerImage.Name,
		customerImage.ID,
		customerImage.DataCenterID,
	))

	imageArtifact := &artifacts.Image{
		Server:        *server,
		NetworkDomain: *networkDomain,
		Image:         *customerImage,
	}
	state.SetTargetImageArtifact(imageArtifact)

	return multistep.ActionContinue
}

// Cleanup is called in reverse order of the steps that have run
// and allow steps to clean up after themselves. Do not assume if this
// ran that the entire multi-step sequence completed successfully. This
// method can be ran in the face of errors and cancellations as well.
//
// The parameter is the same "state bag" as Run, and represents the
// state at the latest possible time prior to calling Cleanup.
func (step *CloneServer) Cleanup(state multistep.StateBag) {
}

var _ multistep.Step = &CloneServer{}
