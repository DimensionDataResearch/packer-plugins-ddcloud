package steps

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"io/ioutil"

	"github.com/DimensionDataResearch/packer-plugins-ddcloud/artifacts"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers/ovftool"
	"github.com/mitchellh/multistep"
)

// ConvertVMXToOVF is the step that converts a .vmx artifact from the VMWare builder to a .ovf artifact (for uploading to CloudControl).
type ConvertVMXToOVF struct {
	// Delete the output directory (and its contents) when the Cleanup function is called?
	CleanupOVF bool

	// The output directory for the OVF package files.
	//
	// If not specified, a new temporary directory will be created and used.
	OutputDir string

	// The path to the VMWare "ovftool" program.
	OVFToolPath string
}

// Run is called to perform the step's action.
//
// The return value determines whether multi-step sequences should continue or halt.
func (step *ConvertVMXToOVF) Run(stateBag multistep.StateBag) multistep.StepAction {
	state := helpers.ForStateBag(stateBag)

	ui := state.GetUI()

	var err error
	if step.OVFToolPath == "" {
		step.OVFToolPath, err = exec.LookPath("ovftool")
		if err != nil {
			state.ShowError(err)

			return multistep.ActionHalt
		}
	}

	log.Printf("ConvertVMXToOVF - using ovftool from '%s'.", step.OVFToolPath)

	sourceArtifact := state.GetSourceArtifact()
	if sourceArtifact == nil {
		err = fmt.Errorf("Cannot find source artifact in state data.")

		state.ShowError(err)

		return multistep.ActionHalt
	}

	if sourceArtifact.BuilderId() != "mitchellh.vmware" {
		err = fmt.Errorf("Source artifact '%s' is of type '%s' (expected 'mitchellh.vmware').",
			sourceArtifact.Id(),
			sourceArtifact.BuilderId(),
		)

		state.ShowError(err)

		return multistep.ActionHalt
	}

	vmxFile := artifacts.GetFirstFileWithExtension(".vmx", sourceArtifact)
	if vmxFile == "" {
		err = fmt.Errorf("Cannot find .vmx file in source artifact '%s'.",
			sourceArtifact.Id(),
		)

		state.ShowError(err)

		return multistep.ActionHalt
	}

	if step.OutputDir == "" {
		step.OutputDir, err = ioutil.TempDir(
			"",                // Use default temp directory
			"packer_vmx_ovf_", // Directory prefix
		)
		if err != nil {
			state.ShowError(err)

			return multistep.ActionHalt
		}
	}

	ovfFile := path.Join(step.OutputDir,
		strings.Replace(
			path.Base(vmxFile),
			".vmx", // Old extension
			".ovf", // New extension
			1,      // Number of replacements
		),
	)

	ovfRunner := ovftool.NewRunnerWithOutputHandler(step.OutputDir, func(programOutput string) {
		ui.Message(fmt.Sprintf(
			"[ovftool] %s",
			programOutput,
		))
	})
	success, err := ovfRunner.Run(
		vmxFile, // From VMX
		ovfFile, // To OVF
	)
	if err != nil {
		state.ShowError(err)

		return multistep.ActionHalt
	}

	if !success {
		err = fmt.Errorf("ovftool exit code does not indicate success")

		state.ShowError(err)

		return multistep.ActionHalt
	}

	// TODO: Implement artifacts.LocalOVFPackage
	// TODO: state.SetSourceArtifact(&artifacts.LocalOVFPackage{x,y,z})

	return multistep.ActionContinue
}

// Cleanup is called in reverse order of the steps that have run
// and allow steps to clean up after themselves. Do not assume if this
// ran that the entire multi-step sequence completed successfully. This
// method can be ran in the face of errors and cancellations as well.
//
// The parameter is the same "state bag" as Run, and represents the
// state at the latest possible time prior to calling Cleanup.
func (step *ConvertVMXToOVF) Cleanup(stateBag multistep.StateBag) {
	if !step.CleanupOVF {
		return // No cleanup required.
	}

	state := helpers.ForStateBag(stateBag)

	if step.OutputDir != "" {
		err := os.RemoveAll(step.OutputDir)
		if err != nil && !os.IsNotExist(err) {
			state.SetLastError(err)
			state.GetUI().Error(
				err.Error(),
			)
		}

		step.OutputDir = ""
	}
}
