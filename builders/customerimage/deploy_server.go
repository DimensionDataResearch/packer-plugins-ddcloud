package main

import (
	"fmt"
	"time"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
)

func deployServer(config Configuration, image compute.Image, networkDomain compute.NetworkDomain, client *compute.Client) (*compute.Server, error) {
	deploymentConfiguration := compute.ServerDeploymentConfiguration{
		Name:        config.serverName,
		Description: fmt.Sprintf("Temporary server created by Packer for image '%s'", config.TargetImage),
		Network: compute.VirtualMachineNetwork{
			NetworkDomainID: config.NetworkDomainID,
			PrimaryAdapter: compute.VirtualMachineNetworkAdapter{
				VLANID: &config.VLANID,
			},
		},
	}
	image.ApplyTo(&deploymentConfiguration)

	serverID, err := client.DeployServer(deploymentConfiguration)
	if err != nil {
		return nil, err
	}

	resource, err := client.WaitForDeploy(compute.ResourceTypeServer, serverID, 20*time.Minute)
	if err != nil {
		return nil, err
	}

	return resource.(*compute.Server), nil
}
