# Customer image builder

The customer image builder deploys a new server in CloudControl, runs any configured provisioners against that server, then clones it to create a new customer image.

## Settings

* `mcp_region` (Required) is the CloudControl region code (e.g. AU, NA, EU, etc).
* `mcp_user` (Required) is the CloudControl user name.  
Can also be specified via the `MCP_USER` environment variable.
* `mcp_password` (Required) is the CloudControl password.  
Can also be specified via the `MCP_PASSWORD` environment variable.
* `datacenter` (Required) is the datacenter Id (must be MCP 2.0).
* `networkdomain` (Required) is the name of the network domain in which to create the server.
* `vlan` is the name of the VLAN to which the server will be attached.
* `source_image` (Required) is the name of the image used to create the server.
* `target_image` (Required) is the name of the customer image to create.
* `use_private_ipv4` (Optional) configures the builder to use private IPv4 addresses rather than public ones (via NAT rules).  
Set this to `true` if you're running packer from inside the MCP 2.0 network domain where the image will be created.
* `client_ip` (Optional) is your client machine's public (external) IP address.  
Required if `use_private_ipv4` is not set.
* `initial_admin_password` (Required unless image does not require) The administrator password to use when deploying the server from which the image will be created.

## Sample configurations

### Create a new customer image in Cloud Control

`build.json`:

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
			"initial_admin_password": "sn4u$ag3$!",
			"client_ip": "1.2.3.4",
			"communicator": "ssh"
		}
	],
	"provisioners": [
		{
			"type": "shell",
			"inline": [
				"echo 'my_password' | sudo -S su",
				"sudo apt-get update"
			]
		}
	]
}
```
