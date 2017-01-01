package main

import (
	"fmt"
	"os"

	"crypto/rand"
	"encoding/hex"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/builders/customerimage/config"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/steps"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/helper/communicator"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/template/interpolate"

	confighelper "github.com/mitchellh/packer/helper/config"
	gossh "golang.org/x/crypto/ssh"
)

// BuilderID is the unique Id for the ddcloud builder
const BuilderID = "ddcloud.image"

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
			&steps.CheckTargetImage{},
			&steps.DeployServer{},
			&steps.CreateNATRule{},
			&steps.CreateFirewallRule{},
			&communicator.StepConnect{
				Config:      &builder.settings.CommunicatorConfig,
				Host:        getSSHHost,
				SSHPort:     getSSHPort,
				SSHConfig:   getSSHConfig,
				WinRMConfig: getWinRMConfig,
			},
			&common.StepProvision{},
			&steps.CloneServer{},
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

func createUniquenessKey() string {
	uniquenessKeyBytes := make([]byte, 5)
	rand.Read(uniquenessKeyBytes)

	return hex.EncodeToString(uniquenessKeyBytes)
}

func getSSHHost(state multistep.StateBag) (host string, err error) {
	settings := state.Get("settings").(*config.Settings)
	host = settings.CommunicatorConfig.SSHHost

	return
}

func getSSHPort(state multistep.StateBag) (port int, err error) {
	settings := state.Get("settings").(*config.Settings)
	port = settings.CommunicatorConfig.SSHPort

	return
}

func getSSHConfig(state multistep.StateBag) (clientConfig *gossh.ClientConfig, err error) {
	settings := state.Get("settings").(*config.Settings)
	clientConfig = &gossh.ClientConfig{
		User: settings.CommunicatorConfig.SSHUsername,
		Auth: []gossh.AuthMethod{
			gossh.Password(settings.CommunicatorConfig.SSHPassword),
		},
	}

	return
}

func getWinRMConfig(state multistep.StateBag) (winRMConfig *communicator.WinRMConfig, err error) {
	settings := state.Get("settings").(*config.Settings)
	winRMConfig = &communicator.WinRMConfig{
		Username: settings.CommunicatorConfig.WinRMUser,
		Password: settings.CommunicatorConfig.WinRMPassword,
	}

	return
}
