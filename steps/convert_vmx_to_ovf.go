package steps

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"strings"

	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers"
	"github.com/mitchellh/multistep"
)

// ConvertVMXToOVF is the step that converts a .vmx artifact from the VMWare builder to a .ovf artifact (for uploading to CloudControl).
type ConvertVMXToOVF struct {
	// Delete the .OVF file (and it's associated artefacts) when the Cleanup function is called.
	CleanupOVF bool

	// The path to the VMWare "ovftool" program.
	ovfToolPath string

	// OVF file (and related files) generated by the step.
	generatedFiles []string
}

// Run is called to perform the step's action.
//
// The return value determines whether multi-step sequences should continue or halt.
func (step *ConvertVMXToOVF) Run(stateBag multistep.StateBag) multistep.StepAction {
	state := helpers.ForStateBag(stateBag)

	ui := state.GetUI()

	var err error
	step.ovfToolPath, err = exec.LookPath("ovftool")
	if err != nil {
		state.SetLastError(err)
		ui.Error(err.Error())

		return multistep.ActionHalt
	}

	log.Printf("ConvertVMXToOVF - using ovftool from '%s'.", step.ovfToolPath)

	step.generatedFiles = []string{}

	sourceArtifact := state.GetSourceArtifact()
	if sourceArtifact == nil {
		err = fmt.Errorf("Cannot find source artifact in state data.")

		state.SetLastError(err)
		ui.Error(err.Error())

		return multistep.ActionHalt
	}

	// TODO: Verify sourceArtifact.BuilderId()

	var vmxFile string
	for _, sourceFile := range sourceArtifact.Files() {
		if strings.HasSuffix(sourceFile, ".vmx") {
			vmxFile = sourceFile

			break
		}
	}
	if vmxFile == "" {
		err = fmt.Errorf("Cannot find .vmx file in source artifact '%s'.",
			sourceArtifact.Id(),
		)

		state.SetLastError(err)
		ui.Error(err.Error())

		return multistep.ActionHalt
	}

	// TODO: Use ovftool.NewRunnerWithOutputHandler().Run().
	// TODO: Remember to pipe output to UI with "[ovftool]" prefix.

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
	if !step.CleanupOVF || step.generatedFiles == nil {
		return // No cleanup required.
	}

	for _, generatedFile := range step.generatedFiles {
		os.Remove(generatedFile)
	}
	step.generatedFiles = nil
}
