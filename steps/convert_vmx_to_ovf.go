package steps

import (
	"fmt"
	"os"
	"path"

	"io/ioutil"

	"github.com/DimensionDataResearch/packer-plugins-ddcloud/artifacts"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

// ConvertVMXToOVF is the step that converts a .vmx artifact from the VMWare builder to a .ovf artifact (for uploading to CloudControl).
type ConvertVMXToOVF struct {
	// Delete the output directory (and its contents) when the Cleanup function is called?
	CleanupOVF bool

	// The base name for .ovf package files.
	PackageName string

	// The output directory for the OVF package files.
	//
	// If not specified, a new temporary directory will be created and used.
	OutputDir string

	// The degree of disk compression (1-9, where 1 is minimum compression and 9 is maximum compression).
	DiskCompression int

	// The path to the VMWare "ovftool" executable.
	OVFExecutable string
}

// Run is called to perform the step's action.
//
// The return value determines whether multi-step sequences should continue or halt.
func (step *ConvertVMXToOVF) Run(stateBag multistep.StateBag) multistep.StepAction {
	var err error

	state := helpers.ForStateBag(stateBag)
	ui := state.GetUI()

	// Minimum compression, by default.
	if step.DiskCompression < 1 || step.DiskCompression > 9 {
		step.DiskCompression = 1
	}

	// Auto-detect tool location if not already specified.
	if step.OVFExecutable == "" {
		step.OVFExecutable = "ovftool"
	}

	vmxFile, err := step.getSourceVMXFile(state)
	if err != nil {
		state.ShowError(err)

		return multistep.ActionHalt
	}

	ovfFile, err := step.getTargetOVFFile()
	if err != nil {
		state.ShowError(err)

		return multistep.ActionHalt
	}

	ovfTool, err := step.createOVFTool(ui)
	if err != nil {
		state.ShowError(err)

		return multistep.ActionHalt
	}

	compression := fmt.Sprintf("--compress=%d",
		step.DiskCompression,
	)
	diskMode := "--diskMode=monolithicSparse"

	success, err := ovfTool.Run(
		compression, // Disk compression level
		diskMode,    // VM disk format (single file, sparse)
		vmxFile,     // From VMX
		ovfFile,     // To OVF
	)
	if err != nil {
		state.ShowError(err)

		return multistep.ActionHalt
	}

	if !success {
		state.ShowErrorMessage("ovftool exit code does not indicate success")

		return multistep.ActionHalt
	}

	ovfArtifact, err := artifacts.NewFromFilesInLocalDirectory(step.OutputDir, "ddcloud.ovf")
	if err != nil {
		state.ShowError(err)

		return multistep.ActionHalt
	}
	state.SetSourceArtifact(ovfArtifact)

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

// Verify that the output directory has been configured and does not yet exist.
func (step *ConvertVMXToOVF) ensureOutputDirectory() (err error) {
	if step.OutputDir == "" {
		step.OutputDir, err = ioutil.TempDir(
			"",                // Use default temp directory
			"packer_vmx_ovf_", // Directory prefix
		)
		if err != nil {
			return
		}
	} else {
		// Verify that target directory does not exist.
		_, err = os.Stat(step.OutputDir)
		if err != nil && !os.IsNotExist(err) {
			err = fmt.Errorf(
				"Output directory '%s' already exists",
				step.OutputDir,
			)

			return
		}

		err = os.MkdirAll(step.OutputDir, 0700 /* u=rwx */)
		if err != nil {
			return
		}
	}

	return
}

// Retrieve the .vmx file path from the source artifact in state data.
func (step *ConvertVMXToOVF) getSourceVMXFile(state helpers.State) (vmxFilePath string, err error) {
	sourceArtifact := state.GetSourceArtifact()
	if sourceArtifact == nil {
		err = fmt.Errorf(
			"Cannot find source artifact in state data.",
		)

		return
	}

	if sourceArtifact.BuilderId() != "mitchellh.vmware" {
		err = fmt.Errorf(
			"Source artifact '%s' is of type '%s' (expected 'mitchellh.vmware').",
			sourceArtifact.Id(),
			sourceArtifact.BuilderId(),
		)

		return
	}

	vmxFile := artifacts.GetFirstFileWithExtension(".vmx", sourceArtifact)
	if vmxFile == "" {
		err = fmt.Errorf(
			"Cannot find .vmx file in source artifact '%s'.",
			sourceArtifact.Id(),
		)

		return
	}

	return
}

// Get the target .ovf file path.
func (step *ConvertVMXToOVF) getTargetOVFFile() (targetOVFFile string, err error) {
	err = step.ensureOutputDirectory()
	if err != nil {
		return
	}

	targetOVFFile = path.Join(step.OutputDir,
		fmt.Sprintf("%s.ovf", step.PackageName),
	)

	return
}

func (step *ConvertVMXToOVF) createOVFTool(ui packer.Ui) (ovfTool *helpers.ToolHelper, err error) {
	return helpers.ForTool(step.OVFExecutable, step.OutputDir, func(programOutput string) {
		ui.Message(fmt.Sprintf(
			"[ovftool] %s",
			programOutput,
		))
	})
}
