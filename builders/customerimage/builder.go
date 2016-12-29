package main

import (
	"fmt"
	"os"

	"encoding/hex"
	"math/rand"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
	"github.com/mitchellh/mapstructure"
	"github.com/mitchellh/packer/packer"
)

// BuilderID is the unique Id for the ddcloud builder
const BuilderID = "dimension-data-research.ddcloud"

// Builder is the Builder plugin for Packer.
type Builder struct {
	config *Configuration
	client *compute.Client
}

// Prepare the plugin to run.
func (builder *Builder) Prepare(configuration ...interface{}) (warnings []string, err error) {
	if len(configuration) == 0 {
		err = fmt.Errorf("No configuration")

		return
	}

	uniquenessKeyBytes := make([]byte, 5)
	_, err = rand.Read(uniquenessKeyBytes)
	if err != nil {
		return
	}
	uniquenessKey := hex.EncodeToString(uniquenessKeyBytes)

	builder.config = &Configuration{
		uniquenessKey: uniquenessKey,
		serverName:    fmt.Sprintf("packer-build-%s", uniquenessKey),
	}
	err = mapstructure.Decode(configuration[0], builder.config)
	if err != nil {
		return
	}

	err = builder.config.Validate()
	if err != nil {
		return
	}
	builder.client = compute.NewClient(
		builder.config.McpRegion,
		builder.config.McpUser,
		builder.config.McpPassword,
	)
	if os.Getenv("MCP_EXTENDED_LOGGING") != "" {
		builder.client.EnableExtendedLogging()
	}

	return
}

// Run the plugin.
func (builder *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	client := builder.client
	config := builder.config

	/*
	 * TODO: Use https://github.com/mitchellh/packer/blob/master/common/step_provision.go (and the machinery that goes with it)
	 */

	networkDomain, err := client.GetNetworkDomain(config.NetworkDomainID)
	if err != nil {
		return nil, err
	}
	if networkDomain == nil {
		return nil, fmt.Errorf("Network domain '%s' not found.", config.NetworkDomainID)
	}

	ui.Message(fmt.Sprintf(
		"Deploy server '%s' from image '%s' in network domain '%s' (datacenter '%s').",
		config.serverName,
		builder.config.SourceImage,
		networkDomain.Name,
		networkDomain.DatacenterID,
	))

	image, err := resolveServerImage(config.SourceImage, networkDomain.DatacenterID, client)
	if err != nil {
		return nil, err
	}
	if image == nil {
		return nil, fmt.Errorf("Cannot find source image named '%s' in datacenter '%s'",
			config.SourceImage,
			networkDomain.DatacenterID,
		)
	}

	ui.Message(fmt.Sprintf(
		"Deploying server '%s' in network domain '%s' ('%s')...",
		config.serverName,
		networkDomain.Name,
		networkDomain.ID,
	))

	server, err := deployServer(*config, image, *networkDomain, client)
	if err != nil {
		return nil, err
	}

	ui.Message(fmt.Sprintf(
		"Deployed server '%s' ('%s').",
		server.Name,
		server.ID,
	))

	ui.Message(fmt.Sprintf(
		"Creating image '%s' from server '%s' ('%s')...",
		config.TargetImage,
		server.Name,
		server.ID,
	))

	customerImage, err := cloneServer(*config, *server, *networkDomain, client)
	if err != nil {
		return nil, err
	}

	ui.Message(fmt.Sprintf(
		"Created customer image '%s' ('%s') from server '%s' ('%s').",
		customerImage.Name,
		customerImage.ID,
		server.Name,
		server.ID,
	))

	ui.Message(fmt.Sprintf(
		"Destroying server '%s' ('%s')...",
		server.Name,
		server.ID,
	))

	err = destroyServer(server.ID, client)
	if err != nil {
		return nil, err
	}

	imageArtifact := &ImageArtifact{
		Server:        *server,
		NetworkDomain: *networkDomain,
		Image:         *customerImage,

		// TODO: deleteImage.
	}

	return imageArtifact, nil
}

// Cancel plugin execution.
func (builder *Builder) Cancel() {
	if builder.client != nil {
		builder.client.Cancel()
	}
}
