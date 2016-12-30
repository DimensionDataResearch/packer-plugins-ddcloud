package steps

import (
	"fmt"
	"time"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/builders/customerimage/config"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

// DeployServer is the step that deploys the target server in CloudControl.
type DeployServer struct{}

// Run is called to perform the step's action.
//
// The return value determines whether multi-step sequences should continue or halt.
func (step *DeployServer) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)

	settings := state.Get("settings").(*config.Settings)
	client := state.Get("client").(*compute.Client)
	networkDomain := state.Get("network_domain").(*compute.NetworkDomain)
	vlan := state.Get("vlan").(*compute.VLAN)
	image := state.Get("source_image").(compute.Image)

	ui.Message(fmt.Sprintf(
		"Deploying server '%s' in network domain '%s' ('%s')...",
		settings.ServerName,
		networkDomain.Name,
		networkDomain.ID,
	))

	deploymentConfiguration := compute.ServerDeploymentConfiguration{
		Name:                  settings.ServerName,
		Description:           fmt.Sprintf("Temporary server created by Packer for image '%s'", settings.TargetImage),
		AdministratorPassword: settings.UniquenessKey,
		Network: compute.VirtualMachineNetwork{
			NetworkDomainID: networkDomain.ID,
			PrimaryAdapter: compute.VirtualMachineNetworkAdapter{
				VLANID: &vlan.ID,
			},
		},
		Start: true, // TODO: Is it possible to auto-start the server only when one or more provisioners are configured?
	}
	image.ApplyTo(&deploymentConfiguration)

	serverID, err := client.DeployServer(deploymentConfiguration)
	if err != nil {
		ui.Error(err.Error())

		return multistep.ActionHalt
	}

	resource, err := client.WaitForDeploy(compute.ResourceTypeServer, serverID, 20*time.Minute)
	if err != nil {
		ui.Error(err.Error())

		return multistep.ActionHalt
	}

	server := resource.(*compute.Server)
	state.Put("server", server)

	serverIPv4 := *server.Network.PrimaryAdapter.PrivateIPv4Address
	ui.Message(fmt.Sprintf(
		"Server '%s' has IPv4 address '%s'.",
		server.Name,
		serverIPv4,
	))
	settings.CommunicatorConfig.SSHHost = serverIPv4
	settings.CommunicatorConfig.WinRMHost = serverIPv4

	ui.Message(fmt.Sprintf(
		"Deployed server '%s' ('%s') in network domain '%s' ('%s')...",
		server.Name,
		server.ID,
		networkDomain.Name,
		networkDomain.ID,
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
func (step *DeployServer) Cleanup(state multistep.StateBag) {}

var _ multistep.Step = &DeployServer{}
