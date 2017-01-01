package main

import (
	"fmt"
	"os"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/postprocessors/customerimage-export/config"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/steps"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/template/interpolate"

	"github.com/DimensionDataResearch/packer-plugins-ddcloud/artifacts"
	confighelper "github.com/mitchellh/packer/helper/config"
)

// PostProcessor is the customer image export post-processor plugin for Packer.
type PostProcessor struct {
	settings             *config.Settings
	interpolationContext interpolate.Context
	client               *compute.Client
	runner               multistep.Runner
}

// Configure is responsible for setting up configuration, storing the state for later,
// and returning and errors, such as validation errors.
func (postProcessor *PostProcessor) Configure(settings ...interface{}) (err error) {
	if len(settings) == 0 {
		err = fmt.Errorf("No settings")

		return
	}

	// Builder settings.
	postProcessor.settings = &config.Settings{}
	err = confighelper.Decode(postProcessor.settings, &confighelper.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &postProcessor.interpolationContext,
	}, settings...)
	if err != nil {
		return
	}

	err = postProcessor.settings.Validate()
	if err != nil {
		return
	}

	postProcessor.client = compute.NewClient(
		postProcessor.settings.McpRegion,
		postProcessor.settings.McpUser,
		postProcessor.settings.McpPassword,
	)
	if os.Getenv("MCP_EXTENDED_LOGGING") != "" {
		postProcessor.client.EnableExtendedLogging()
	}

	// Configure post-processor execution logic.
	postProcessor.runner = &multistep.BasicRunner{
		Steps: []multistep.Step{
			&steps.ResolveSourceImage{
				MustBeCustomerImage: true,
			},
			&steps.ExportCustomerImage{},
		},
	}

	return nil
}

// PostProcess takes a previously created Artifact and produces another Artifact.
//
// If an error occurs, it should return that error.
// If `keep` is to true, then the previous artifact is forcibly kept.
func (postProcessor *PostProcessor) PostProcess(ui packer.Ui, sourceArtifact packer.Artifact) (destinationArtifact packer.Artifact, keep bool, err error) {
	if sourceArtifact.BuilderId() != "ddcloud.image" {
		err = fmt.Errorf("The source artifact is not a CloudControl image.")

		return
	}

	settings := postProcessor.settings
	packerConfig := &settings.PackerConfig
	client := postProcessor.client

	var targetImage *compute.CustomerImage
	targetImageID := sourceArtifact.Id()
	targetImage, err = client.GetCustomerImage(targetImageID)
	if err != nil {
		return
	}
	if targetImage == nil {
		err = fmt.Errorf("Cannot find customer image '%s'",
			targetImageID,
		)

		return
	}

	stepState := helpers.ForStateBag(
		&multistep.BasicStateBag{},
	)
	stepState.SetUI(ui)
	stepState.SetPackerConfig(packerConfig)
	stepState.SetSettings(settings)
	stepState.SetClient(client)
	stepState.SetTargetImage(targetImage)
	stepState.SetTargetImageArtifact(&artifacts.Image{
		Image:     targetImage,
		BuilderID: sourceArtifact.BuilderId(),
	})
	postProcessor.runner.Run(stepState.Data)

	err = stepState.GetLastError()
	if err != nil {
		return
	}

	destinationArtifact = stepState.GetRemoteOVFPackageArtifact()

	return
}

var _ packer.PostProcessor = &PostProcessor{}
