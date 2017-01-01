package main

import (
	"fmt"
	"os"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/builders/customerimage-import/config"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/steps"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/template/interpolate"

	confighelper "github.com/mitchellh/packer/helper/config"
)

// BuilderID is the unique Id for the ddcloud builder
const BuilderID = "ddcloud.image"

// Builder is the customer image import builder plugin for Packer.
type Builder struct {
	settings             *config.Settings
	interpolationContext interpolate.Context
	client               *compute.Client
	runner               multistep.Runner
}

// Prepare the plugin to run.
func (builder *Builder) Prepare(settings ...interface{}) (warnings []string, err error) {
	if len(settings) == 0 {
		err = fmt.Errorf("No settings")

		return
	}

	// Builder settings.
	builder.settings = &config.Settings{}
	err = confighelper.Decode(builder.settings, &confighelper.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &builder.interpolationContext,
	}, settings...)
	if err != nil {
		return
	}

	err = builder.settings.Validate()
	if err != nil {
		return
	}

	builder.client = compute.NewClient(
		builder.settings.McpRegion,
		builder.settings.McpUser,
		builder.settings.McpPassword,
	)
	if os.Getenv("MCP_EXTENDED_LOGGING") != "" {
		builder.client.EnableExtendedLogging()
	}

	// Resolve the target datacenter.
	var targetDatacenter *compute.Datacenter
	targetDatacenter, err = builder.client.GetDatacenter(
		builder.settings.DatacenterID,
	)
	if err != nil {
		return
	}
	if targetDatacenter == nil {
		err = fmt.Errorf("Cannot find target datacenter '%s'",
			builder.settings.DatacenterID,
		)

		return
	}

	// Configure builder execution logic.
	builder.runner = &multistep.BasicRunner{
		Steps: []multistep.Step{
			&steps.ResolveDatacenter{
				AsTarget: true,
			},
			// TODO: Implement ImportCustomerImage step.
		},
	}

	return
}

// Run the plugin.
func (builder *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	settings := builder.settings
	packerConfig := &settings.PackerConfig
	client := builder.client

	stepState := helpers.ForStateBag(
		&multistep.BasicStateBag{},
	)
	stepState.SetUI(ui)
	stepState.SetHook(hook)
	stepState.SetPackerConfig(packerConfig)
	stepState.SetSettings(settings)
	stepState.SetClient(client)
	stepState.SetBuilderID(BuilderID)
	builder.runner.Run(stepState.Data)

	err := stepState.GetLastError()
	if err != nil {
		return nil, err
	}

	imageArtifact := stepState.GetTargetImageArtifact()
	if imageArtifact == nil {
		// TODO: This currently fails because the step sequence has not actually been implemented yet.
		return nil, fmt.Errorf("One or more steps failed to complete")
	}

	return imageArtifact, nil
}

// Cancel plugin execution.
func (builder *Builder) Cancel() {
	if builder.client != nil {
		builder.client.Cancel()
	}
}
