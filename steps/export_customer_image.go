package steps

import (
	"fmt"

	"time"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/artifacts"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/postprocessors/customerimage-export/config"
	"github.com/mitchellh/multistep"
)

// ExportCustomerImage is the step that exports a customer image as an OVF package.
type ExportCustomerImage struct{}

// Run is called to perform the step's action.
//
// The return value determines whether multi-step sequences should continue or halt.
func (step *ExportCustomerImage) Run(stateBag multistep.StateBag) multistep.StepAction {
	state := helpers.ForStateBag(stateBag)
	ui := state.GetUI()

	settings := state.GetSettings().(*config.Settings)
	client := state.GetClient()

	targetImage := state.GetTargetImage()
	targetImageID := targetImage.GetID()
	targetImageName := targetImage.GetName()

	targetImageArtifact := state.GetTargetImageArtifact()
	if targetImageArtifact == nil {
		ui.Error("Cannot find an artifact for the image being exported.")

		return multistep.ActionHalt
	}

	ui.Message(fmt.Sprintf(
		"Export customer image '%s' ('%s') to OVF package '%s'.",
		targetImageName,
		targetImageID,
		settings.OVFPackagePrefix,
	))

	exportID, err := client.ExportCustomerImage(targetImageID, settings.OVFPackagePrefix)
	if err != nil {
		ui.Error(err.Error())

		return multistep.ActionHalt
	}

	ui.Message(fmt.Sprintf(
		"Export of customer image '%s' in progress with export ID '%s'...",
		targetImageName,
		exportID,
	))

	_, err = client.WaitForChange(compute.ResourceTypeCustomerImage, targetImageID, "Export", 30*time.Minute)
	if err != nil {
		ui.Error(err.Error())

		return multistep.ActionHalt
	}

	datacenterMetadata, err := client.GetDatacenter(targetImage.DataCenterID)
	if err != nil {
		ui.Error(err.Error())

		return multistep.ActionHalt
	}
	if datacenterMetadata == nil {
		ui.Error(fmt.Sprintf(
			"Can't find target datacenter '%s'",
			targetImage.DataCenterID,
		))

		return multistep.ActionHalt
	}

	state.SetRemoteOVFPackageArtifact(&artifacts.RemoteOVFPackage{
		FTPSHostName:  datacenterMetadata.FTPSHost,
		PackagePrefix: settings.OVFPackagePrefix,
		BuilderID:     state.GetBuilderID(),
	})

	ui.Message(fmt.Sprintf(
		"Export of customer image '%s' complete.",
		targetImageName,
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
func (step *ExportCustomerImage) Cleanup(state multistep.StateBag) {
}
