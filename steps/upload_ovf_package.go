package steps

import (
	"crypto/tls"
	"fmt"
	"log"
	"path"

	"os"

	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"github.com/secsy/goftp"
)

// UploadOVFPackage is the step that uploads the files comprising an OVF package to CloudControl.
type UploadOVFPackage struct{}

// Run is called to perform the step's action.
//
// The return value determines whether multi-step sequences should continue or halt.
func (step *UploadOVFPackage) Run(stateBag multistep.StateBag) multistep.StepAction {
	state := helpers.ForStateBag(stateBag)
	ui := state.GetUI()
	targetDatacenter := state.GetTargetDatacenter()

	ui.Message(fmt.Sprintf(
		"Uploading OVF package files to datacenter '%s'...",
		targetDatacenter.ID,
	))

	sourceArtifact := state.GetSourceArtifact()
	if sourceArtifact.BuilderId() != "ddcloud.ovf" {
		state.ShowError(fmt.Errorf(
			"Source image '%s' is of type '%s' (expected 'ddcloud.ovf')",
			sourceArtifact.Id(),
			sourceArtifact.BuilderId(),
		))

		return multistep.ActionHalt
	}

	sourceFiles := sourceArtifact.Files()
	err := validateOVFPackageArtifactFiles(sourceFiles)
	if err != nil {
		state.ShowError(err)

		return multistep.ActionHalt
	}

	settings := state.GetSettings()
	ftpConfig := goftp.Config{
		User:     settings.GetMCPUser(),
		Password: settings.GetMCPPassword(),
		TLSConfig: &tls.Config{
			ServerName: targetDatacenter.FTPSHost,
		},
	}
	ftpClient, err := goftp.DialConfig(ftpConfig, targetDatacenter.FTPSHost)
	if err != nil {
		ui.Error(fmt.Sprintf(
			"Failed to connect (ftps://%s).",
			targetDatacenter.FTPSHost,
		))
		state.ShowError(err)

		return multistep.ActionHalt
	}
	defer ftpClient.Close()

	for _, sourceFileName := range sourceFiles {
		if !isOVFPackageFile(sourceFileName) {
			log.Printf("Skipping file '%s' (does not look like an OVF package file).", sourceFileName)

			continue
		}

		targetFileName := path.Base(sourceFileName)
		ui.Message(fmt.Sprintf(
			"Uploading '%s'...", targetFileName,
		))

		sourceFile, err := os.Open(sourceFileName)
		if err != nil {
			state.ShowError(err)

			return multistep.ActionHalt
		}
		defer sourceFile.Close()

		err = ftpClient.Store(targetFileName, sourceFile)
		if err != nil {
			state.ShowError(err)

			return multistep.ActionHalt
		}

		ui.Message(fmt.Sprintf(
			"Uploaded '%s'.", targetFileName,
		))
	}

	ui.Message(fmt.Sprintf(
		"Uploaded OVF package files to datacenter '%s'.",
		targetDatacenter.ID,
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
func (step *UploadOVFPackage) Cleanup(state multistep.StateBag) {
}

var _ multistep.Step = &UploadOVFPackage{}

// Is the specified file part of an OVF package (from Cloud Control's point of view)?
func isOVFPackageFile(fileName string) bool {
	switch path.Base(fileName) {
	case ".vmdk": // VM disk
	case ".ovf": // OVF
	case ".mf": // Manifest
		return true
	}

	return false
}

func validateOVFPackageArtifactFiles(packageFiles []string) (err error) {
	var haveVMDK, haveOVF, haveMF bool
	for _, packageFile := range packageFiles {
		switch path.Ext(packageFile) {
		case ".vmdk":
			haveVMDK = true

		case ".ovf":
			haveOVF = true

		case ".mf":
			haveMF = true
		}
	}

	if !haveMF {
		err = packer.MultiErrorAppend(err, fmt.Errorf(
			"Source artifact is missing .mf file",
		))
	}
	if !haveOVF {
		err = packer.MultiErrorAppend(err, fmt.Errorf(
			"Source artifact is missing .ovf file",
		))
	}
	if !haveVMDK {
		err = packer.MultiErrorAppend(err, fmt.Errorf(
			"Source artifact is missing .vmdk file",
		))
	}

	return
}
