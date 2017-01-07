package steps

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/DimensionDataResearch/packer-plugins-ddcloud/artifacts"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"github.com/secsy/goftp"
)

// UploadOVFPackage is the step that uploads the files comprising an OVF package to CloudControl.
//
// Expects:
//   - Target data center in state from ResolveDatacenter step.
//   - OVF package files from source artifact in state from ConvertVMXToOVF step.
type UploadOVFPackage struct{}

// Run is called to perform the step's action.
//
// The return value determines whether multi-step sequences should continue or halt.
func (step *UploadOVFPackage) Run(stateBag multistep.StateBag) multistep.StepAction {
	state := helpers.ForStateBag(stateBag)
	ui := state.GetUI()

	targetDatacenter := state.GetTargetDatacenter()
	if targetDatacenter == nil {
		state.ShowErrorMessage("Cannot found target datacenter in state data.")

		return multistep.ActionHalt
	}

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

	log.Printf("UploadOVFPackage: source package is %s", sourceArtifact.String())

	sourceFiles := sourceArtifact.Files()
	err := step.validateOVFPackageFiles(sourceFiles)
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

	packageBaseName := ""
	for _, sourceFileName := range sourceFiles {
		if !isOVFPackageFile(sourceFileName) {
			log.Printf("Skipping file '%s' (does not look like an OVF package file).", sourceFileName)

			continue
		}

		targetFileName := path.Base(sourceFileName)
		ui.Message(fmt.Sprintf(
			"Uploading '%s'...", targetFileName,
		))

		if path.Ext(targetFileName) == ".ovf" {
			packageBaseName = strings.Replace(targetFileName, ".ofv", "", 1)
		}

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

	state.SetRemoteOVFPackageArtifact(&artifacts.RemoteOVFPackage{
		FTPSHostName:  targetDatacenter.FTPSHost,
		PackagePrefix: packageBaseName,
		BuilderID:     "ddcloud.ovf",
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

func (step *UploadOVFPackage) validateOVFPackageFiles(packageFiles []string) (err error) {
	var haveVMDK, haveOVF, haveMF bool
	for _, packageFile := range packageFiles {
		log.Printf("UploadOVFPackage: validating package file '%s'...", packageFile)

		haveVMDK = haveVMDK || strings.HasSuffix(packageFile, ".vmdk") || strings.HasSuffix(packageFile, ".vmdk.gz")
		haveOVF = haveOVF || strings.HasSuffix(packageFile, ".ovf")
		haveMF = haveMF || strings.HasSuffix(packageFile, ".mf")
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
