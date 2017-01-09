# Customer image import post-processor

The customer image import post-processor converts a local VMWare (`.vmx`) virtual machine into OVF (`.ovf`) format, uploads it to CloudControl, and then imports it as a customer image.

## Settings

* `mcp_region` (Required) is the CloudControl region code (e.g. AU, NA, EU, etc).
* `mcp_user` (Required) is the CloudControl user name.  
Can also be specified via the `MCP_USER` environment variable.
* `mcp_password` (Required) is the CloudControl password.  
Can also be specified via the `MCP_PASSWORD` environment variable.
* `datacenter` (Required) is the Id of the datacenter where the image will be imported (must be MCP 2.0).
* `target_image` (Required) is the name of the customer image to create.
* `ovf_package_prefix` (Optional) is the prefix used to name the OVF package files.  
If not specified, `target_image` is used.
