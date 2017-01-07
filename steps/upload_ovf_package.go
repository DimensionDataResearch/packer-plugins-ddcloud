package steps

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/DimensionDataResearch/packer-plugins-ddcloud/artifacts"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

// UploadOVFPackage is the step that uploads the files comprising an OVF package to CloudControl.
//
// Expects:
//   - Target data center in state from ResolveDatacenter step.
//   - OVF package files from source artifact in state from ConvertVMXToOVF step.
type UploadOVFPackage struct {
	// The path to the "curl" executable.
	CurlExecutable string
}

// Run is called to perform the step's action.
//
// The return value determines whether multi-step sequences should continue or halt.
func (step *UploadOVFPackage) Run(stateBag multistep.StateBag) multistep.StepAction {
	state := helpers.ForStateBag(stateBag)
	ui := state.GetUI()

	// Auto-detect tool location if not already specified.
	if step.CurlExecutable == "" {
		step.CurlExecutable = "curl"
	}

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

	curlTool, err := step.createCurlTool(ui)
	if err != nil {
		state.ShowError(err)

		return multistep.ActionHalt
	}

	packageBaseName := ""
	for _, sourceFile := range sourceFiles {
		if !isOVFPackageFile(sourceFile) {
			log.Printf("Skipping file '%s' (does not look like an OVF package file).", sourceFile)

			continue
		}

		targetFileName := path.Base(sourceFile)
		ui.Message(fmt.Sprintf(
			"Uploading '%s'...", targetFileName,
		))

		if strings.HasSuffix(targetFileName, ".ovf") {
			packageBaseName = strings.Replace(targetFileName, ".ofv", "", 1)
		}

		success, err := curlTool.Run(
			"-s", // No progress bar
			"-S", // But still show errors
			"--user",
			fmt.Sprintf("%s:%s",
				settings.GetMCPUser(),
				settings.GetMCPPassword(),
			),
			"--upload-file",
			sourceFile,
			"--ssl", // FTPS
			fmt.Sprintf("ftp://%s/%s",
				targetDatacenter.FTPSHost,
				targetFileName,
			),
		)
		if err != nil {
			state.ShowError(err)

			return multistep.ActionHalt
		}
		if !success {
			state.ShowErrorMessage("Failed to upload file '%s' to '%s'",
				sourceFile,
				targetDatacenter.FTPSHost,
			)

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
	return strings.HasSuffix(fileName, ".vmdk") ||
		strings.HasSuffix(fileName, ".mf") ||
		strings.HasSuffix(fileName, ".ovf")
}

// Verify that the OFV package files include all required file types (OVF, MF, and VMDK).
func (step *UploadOVFPackage) validateOVFPackageFiles(packageFiles []string) (err error) {
	var haveVMDK, haveOVF, haveMF bool
	for _, packageFile := range packageFiles {
		log.Printf("UploadOVFPackage: validating package file '%s'...", packageFile)

		haveVMDK = haveVMDK || strings.HasSuffix(packageFile, ".vmdk")
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

func (step *UploadOVFPackage) createCurlTool(ui packer.Ui) (*helpers.Tool, error) {
	workDir, _ := os.Getwd()

	return helpers.ForTool(step.CurlExecutable, workDir, func(programOutput string) {
		ui.Message(fmt.Sprintf(
			"[curl] %s",
			programOutput,
		))
	})
}
