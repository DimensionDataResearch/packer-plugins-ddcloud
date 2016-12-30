package main

import (
	"fmt"
	"os"

	"crypto/rand"
	"encoding/hex"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/builders/customerimage/artifacts"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/builders/customerimage/config"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/builders/customerimage/steps"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/helper/communicator"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/template/interpolate"

	confighelper "github.com/mitchellh/packer/helper/config"
	gossh "golang.org/x/crypto/ssh"
)

// Builder is the Builder plugin for Packer.
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

	builder.settings.UniquenessKey = createUniquenessKey()
	builder.settings.ServerName = fmt.Sprintf("packer-build-%s", builder.settings.UniquenessKey)

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
			&steps.ResolveNetworkDomain{},
			&steps.ResolveVLAN{},
			&steps.ResolveSourceImage{},
			&steps.DeployServer{},
			&steps.CreateNATRule{},
			&steps.CreateFirewallRule{},
			&communicator.StepConnect{
				Config: &builder.settings.CommunicatorConfig,
				Host: func(state multistep.StateBag) (host string, err error) {
					settings := state.Get("settings").(*config.Settings)
					host = settings.CommunicatorConfig.SSHHost

					return
				},
				SSHPort: func(state multistep.StateBag) (port int, err error) {
					settings := state.Get("settings").(*config.Settings)
					port = settings.CommunicatorConfig.SSHPort

					return
				},
				SSHConfig: func(state multistep.StateBag) (clientConfig *gossh.ClientConfig, err error) {
					settings := state.Get("settings").(*config.Settings)
					clientConfig = &gossh.ClientConfig{
						User: settings.CommunicatorConfig.SSHUsername,
						Auth: []gossh.AuthMethod{
							gossh.Password(settings.CommunicatorConfig.SSHPassword),
						},
					}

					return
				},
			},
			&common.StepProvision{},
			&steps.CloneServer{},
			&steps.DestroyServer{},
		},
	}

	return
}

// Run the plugin.
func (builder *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	settings := builder.settings
	packerConfig := &settings.PackerConfig
	client := builder.client

	stepState := &multistep.BasicStateBag{}
	stepState.Put("ui", ui)
	stepState.Put("hook", hook)
	stepState.Put("config", packerConfig)
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

func createUniquenessKey() string {
	uniquenessKeyBytes := make([]byte, 5)
	rand.Read(uniquenessKeyBytes)

	return hex.EncodeToString(uniquenessKeyBytes)
}
