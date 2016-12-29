package main

import (
	"fmt"
	"os"

	"encoding/hex"
	"math/rand"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/builders/customerimage/artifacts"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/builders/customerimage/config"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/builders/customerimage/steps"
	"github.com/mitchellh/mapstructure"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

// Builder is the Builder plugin for Packer.
type Builder struct {
	settings *config.Settings
	client   *compute.Client
	runner   multistep.Runner
}

// Prepare the plugin to run.
func (builder *Builder) Prepare(settings ...interface{}) (warnings []string, err error) {
	if len(settings) == 0 {
		err = fmt.Errorf("No settings")

		return
	}

	uniquenessKeyBytes := make([]byte, 5)
	_, err = rand.Read(uniquenessKeyBytes)
	if err != nil {
		return
	}
	uniquenessKey := hex.EncodeToString(uniquenessKeyBytes)

	builder.settings = &config.Settings{
		UniquenessKey: uniquenessKey,
		ServerName:    fmt.Sprintf("packer-build-%s", uniquenessKey),
	}
	err = mapstructure.Decode(settings[0], builder.settings)
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

	// Configure builder execution logic.
	builder.runner = &multistep.BasicRunner{
		Steps: []multistep.Step{
			steps.ResolveNetworkDomain{},
			steps.ResolveSourceImage{},
			steps.DeployServer{},
			steps.CloneServer{},
			steps.DestroyServer{},
		},
	}

	return
}

// Run the plugin.
func (builder *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	settings := builder.settings
	client := builder.client

	stepState := &multistep.BasicStateBag{}
	stepState.Put("ui", ui)
	stepState.Put("hook", hook)
	stepState.Put("settings", settings)
	stepState.Put("client", client)
	builder.runner.Run(stepState)

	rawError, ok := stepState.GetOk("error")
	if ok {
		return nil, rawError.(error)
	}

	rawImageArtifact, ok := stepState.GetOk("target_image_artifact")
	if !ok {
		return nil, fmt.Errorf("One or more steps failed to complete")
	}

	imageArtifact := rawImageArtifact.(*artifacts.Image)

	return imageArtifact, nil
}

// Cancel plugin execution.
func (builder *Builder) Cancel() {
	if builder.client != nil {
		builder.client.Cancel()
	}
}
