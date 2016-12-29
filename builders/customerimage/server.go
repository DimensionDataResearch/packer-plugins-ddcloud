package main

import (
	"fmt"
	"time"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
)

func deployServer(config Configuration, image compute.Image, networkDomain compute.NetworkDomain, client *compute.Client) (*compute.Server, error) {
	deploymentConfiguration := compute.ServerDeploymentConfiguration{
		Name:                  config.serverName,
		Description:           fmt.Sprintf("Temporary server created by Packer for image '%s'", config.TargetImage),
		AdministratorPassword: config.uniquenessKey,
		Network: compute.VirtualMachineNetwork{
			NetworkDomainID: config.NetworkDomainID,
			PrimaryAdapter: compute.VirtualMachineNetworkAdapter{
				VLANID: &config.VLANID,
			},
		},
		Start: false,
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

// Destroy a server.
func destroyServer(serverID string, client *compute.Client) error {
	err := client.DeleteServer(serverID)
	if err != nil {
		return err
	}

	return client.WaitForDelete(
		compute.ResourceTypeServer,
		serverID,
		5*time.Minute,
	)
}

// Clone a server to create a customer image.
func cloneServer(config Configuration, server compute.Server, networkDomain compute.NetworkDomain, client *compute.Client) (customerImage *compute.CustomerImage, err error) {
	var imageID string
	imageID, err = client.CloneServer(
		server.ID,
		config.TargetImage,
		fmt.Sprintf("%s (created by Packer)", config.TargetImage),
		false, // preventGuestOSCustomisation
	)
	if err != nil {
		return
	}

	var resource compute.Resource
	resource, err = client.WaitForServerClone(
		imageID,
		15*time.Minute,
	)
	if err != nil {
		return
	}

	customerImage = resource.(*compute.CustomerImage)

	return
}
