# CloudControl plugins for Packer

Plugins for [Hashicorp Packer](https://packer.io/) that target Dimension Data CloudControl.

**Note**: this is a work-in-progress; it's not production-ready yet.

Currently, the following plugins are available:

* `ddcloud-customerimage` (Builder) Deploys a server in CloudControl, runs configured provisioners (if any), then clones the server to create a Customer image.  
The deployed server is destroyed once cloning is complete.

We're also planning to create a plugin that uploads an OVF and imports it to create a customer image.

## Installing

There are no pre-built binaries yet, so you'll have to build it yourself for now.
Needs Packer <= 0.8.6, OSX or Linux, Go >= 1.7, and GNU Make.

Sorry, the dependencies are a bit messy at the moment.

1. Fetch correct dependency versions by running `./init.sh` (one-time only).
2. Run `make dev`.

## Example configuration

```json
{
	"builders": [
		{
			"type": "ddcloud-customerimage",
			"mcp_region": "AU",
			"datacenter": "AU9",
			"networkdomain": "MyNetworkDomain",
			"vlan": "MyVLAN",
            "source_image": "Ubuntu 14.04 2 CPU",
			"target_image": "packertest",
			"client_ip": "1.2.3.4",
			"communicator": "ssh"
		}
	],
	"provisioners": [
		{
			"type": "shell",
			"inline": [
				"ls -l /"
			]
		}
	]
}
```

* `datacenter` is the datacenter Id (must be MCP 2.0).
* `networkdomain` is the name of the network domain in which to create the server.
* `vlan` is the name of the VLAN to which the server will be attached.
* `source_image` is the name of the image used to create the server.
* `target_image` is the name of the customer image to create.
* `client_ip` is your client machine's public (external) IP address.

Specify CloudControl username and password using `MCP_USER` and `MCP_PASSWORD` environment variables.